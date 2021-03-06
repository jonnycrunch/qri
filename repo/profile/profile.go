package profile

import (
	"time"

	"github.com/ipfs/go-datastore"
	// ma "gx/ipfs/QmXY77cVe7rVRQXZZQRioukUM7aRW3BTcAgJe12MCtb3Ji/go-multiaddr"
	peer "gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
)

// Profile defines peer profile details
type Profile struct {
	ID ID `json:"id"`
	// Created timestamp
	Created time.Time `json:"created,omitempty"`
	// Updated timestamp
	Updated time.Time `json:"updated,omitempty"`
	// Peername a handle for the user. min 1 character, max 80. composed of [_,-,a-z,A-Z,1-9]
	Peername string `json:"peername"`
	// specifies weather this is a user or an organization
	Type Type `json:"type"`
	// user's email address
	Email string `json:"email"`
	// user name field. could be first[space]last, but not strictly enforced
	Name string `json:"name"`
	// user-filled description of self
	Description string `json:"description"`
	// url this user wants the world to click
	HomeURL string `json:"homeUrl"`
	// color this user likes to use as their theme color
	Color string `json:"color"`
	// Thumb path for user's thumbnail
	Thumb datastore.Key `json:"thumb"`
	// Profile photo
	Photo datastore.Key `json:"photo"`
	// Poster photo for users's profile page
	Poster datastore.Key `json:"poster"`
	// Twitter is a  peer's twitter handle
	Twitter string `json:"twitter"`
	// Addresses lists any network addresses associated with this profile
	// in the form of peer.ID.Pretty() : []multiaddr strings
	// both peer.IDs and multiaddresses are converted to strings for
	// clean en/decoding
	Addresses map[string][]string `json:"addresses"`
}

// PeerIDs sifts through listed multaddrs looking for an IPFS peer ID
func (p *Profile) PeerIDs() (ids []peer.ID) {
	for idstr := range p.Addresses {
		if id, err := peer.IDB58Decode(idstr); err == nil {
			ids = append(ids, id)
		}
	}
	return
}
