package cmd

import (
	"fmt"
	"net"
	"net/rpc"
	"strings"

	ipfs "github.com/qri-io/cafs/ipfs"
	"github.com/qri-io/qri/config"
	"github.com/qri-io/qri/core"
	"github.com/qri-io/qri/p2p"
	"github.com/qri-io/qri/repo"
	"github.com/qri-io/qri/repo/fs"
)

var (
	repository repo.Repo
	rpcClient  *rpc.Client
)

func getRepo(online bool) repo.Repo {
	if repository != nil {
		return repository
	}

	if !QRIRepoInitialized() {
		ErrExit(fmt.Errorf("no qri repo found, please run `qri setup`"))
	}

	fs := getIpfsFilestore(online)
	r, err := fsrepo.NewRepo(fs, core.Config.Profile, QriRepoPath)
	ExitIfErr(err)

	return r
}

func getIpfsFilestore(online bool) *ipfs.Filestore {
	fs, err := ipfs.NewFilestore(func(cfg *ipfs.StoreCfg) {
		cfg.FsRepoPath = IpfsFsPath
		cfg.Online = online
	})
	ExitIfErr(err)
	return fs
}

func datasetRequests(online bool) (*core.DatasetRequests, error) {
	// TODO - bad bad hardcode
	if conn, err := net.Dial("tcp", ":2504"); err == nil {
		return core.NewDatasetRequests(nil, rpc.NewClient(conn)), nil
	}

	if !online {
		// TODO - make this not terrible
		r, cli, err := repoOrClient(online)
		if err != nil {
			return nil, err
		}
		return core.NewDatasetRequests(r, cli), nil
	}

	n, err := qriNode(online)
	if err != nil {
		return nil, err
	}

	req := core.NewDatasetRequests(n.Repo, nil)
	req.Node = n
	return req, nil
}

func profileRequests(online bool) (*core.ProfileRequests, error) {
	r, cli, err := repoOrClient(online)
	if err != nil {
		return nil, err
	}
	return core.NewProfileRequests(r, cli), nil
}

func searchRequests(online bool) (*core.SearchRequests, error) {
	r, cli, err := repoOrClient(online)
	if err != nil {
		return nil, err
	}
	return core.NewSearchRequests(r, cli), nil
}

func historyRequests(online bool) (*core.HistoryRequests, error) {
	// TODO - bad bad hardcode
	if conn, err := net.Dial("tcp", ":2504"); err == nil {
		return core.NewHistoryRequests(nil, rpc.NewClient(conn)), nil
	}

	if !online {
		// TODO - make this not terrible
		r, cli, err := repoOrClient(online)
		if err != nil {
			return nil, err
		}
		return core.NewHistoryRequests(r, cli), nil
	}

	n, err := qriNode(online)
	if err != nil {
		return nil, err
	}

	req := core.NewHistoryRequests(n.Repo, nil)
	req.Node = n
	return req, nil
}

func peerRequests(online bool) (*core.PeerRequests, error) {
	// return nil, nil

	// TODO - bad bad hardcode
	if conn, err := net.Dial("tcp", ":2504"); err == nil {
		return core.NewPeerRequests(nil, rpc.NewClient(conn)), nil
	}

	node, err := qriNode(online)
	if err != nil {
		return nil, err
	}
	return core.NewPeerRequests(node, nil), nil
}

func repoOrClient(online bool) (repo.Repo, *rpc.Client, error) {
	if repository != nil {
		return repository, nil, nil
	} else if rpcClient != nil {
		return nil, rpcClient, nil
	}

	if fs, err := ipfs.NewFilestore(func(cfg *ipfs.StoreCfg) {
		cfg.FsRepoPath = IpfsFsPath
		cfg.Online = online
	}); err == nil {
		r, err := fsrepo.NewRepo(fs, core.Config.Profile, QriRepoPath)
		ExitIfErr(err)

		return r, nil, err

	} else if strings.Contains(err.Error(), "lock") {
		conn, err := net.Dial("tcp", fmt.Sprintf(":%d", core.Config.RPC.Port))
		if err != nil {
			return nil, nil, err
		}
		return nil, rpc.NewClient(conn), nil
	} else {
		return nil, nil, err
	}

	return nil, nil, fmt.Errorf("badbadnotgood")
}

func qriNode(online bool) (node *p2p.QriNode, err error) {
	var (
		r  repo.Repo
		fs *ipfs.Filestore
	)

	fs, err = ipfs.NewFilestore(func(cfg *ipfs.StoreCfg) {
		cfg.FsRepoPath = IpfsFsPath
		cfg.Online = online
	})

	if err != nil {
		return
	}

	r, err = fsrepo.NewRepo(fs, core.Config.Profile, QriRepoPath)
	if err != nil {
		return
	}

	node, err = p2p.NewQriNode(r, func(c *config.P2P) {
		c.Enabled = online
		c.QriBootstrapAddrs = core.Config.P2P.QriBootstrapAddrs
	})
	if err != nil {
		return
	}

	// if online {
	// 	log.Info("p2p addresses:")
	// 	for _, a := range node.EncapsulatedAddresses() {
	// 		log.Infof("  %s", a.String())
	// 	}
	// }

	return
}
