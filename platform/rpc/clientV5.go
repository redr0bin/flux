package rpc

import (
	"errors"
	"io"

	"github.com/weaveworks/flux/platform"
)

// RPCClient is the rpc-backed implementation of a platform, for
// talking to remote daemons.
type RPCClientV5 struct {
	*RPCClientV4
}

var _ platform.PlatformV5 = &RPCClientV5{}

// NewClient creates a new rpc-backed implementation of the platform.
func NewClientV5(conn io.ReadWriteCloser) *RPCClientV5 {
	return &RPCClientV5{NewClientV4(conn)}
}

// Export is used to get service configuration in platform-specific format
func (p *RPCClientV5) Export() ([]byte, error) {
	var config []byte
	err := p.client.Call("RPCServer.Export", struct{}{}, &config)
	return config, CategoriseRPCError(err)
}

func (p *RPCClientV5) Sync(spec platform.SyncDef) error {
	var result SyncResult
	if err := p.client.Call("RPCServer.Sync", spec, &result); err != nil {
		return CategoriseRPCError(err)
	}
	if len(result) > 0 {
		errs := platform.SyncError{}
		for id, msg := range result {
			errs[id] = errors.New(msg)
		}
		return platform.ClusterError(errs)
	}
	return nil
}
