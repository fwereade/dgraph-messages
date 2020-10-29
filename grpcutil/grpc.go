package grpcutil

import (
	"github.com/juju/errors"
	"google.golang.org/grpc"
)

// Spec specifies how we want to connect to some GRPC server.
type Spec struct {

	// Addr is the address.
	Addr string

	// Token, if specified, is the per-request authentication token
	// (specifically, suitable for use with hosted dgraph).
	Token string

	// Huge bumps the maximum allowed message size to a very large value.
	Huge bool
}

// Options returns the GRPC dial options implied by the spec.
func (spec Spec) Options() ([]grpc.DialOption, error) {
	var opts []grpc.DialOption

	if spec.Token == "" {
		opts = append(opts, grpc.WithInsecure())
	} else {
		tc, err := transportCredentials()
		if err != nil {
			return nil, errors.Annotate(err, "cannot create transport credentials")
		}
		opts = append(opts, grpc.WithTransportCredentials(tc))

		rc := rpcCredentials{spec.Token}
		opts = append(opts, grpc.WithPerRPCCredentials(rc))
	}

	if spec.Huge {
		const HugeMessageSize = 1 << 32
		opts = append(opts, grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(HugeMessageSize),
			grpc.MaxCallSendMsgSize(HugeMessageSize),
		))
	}

	return opts, nil
}

// Connect returns the GRPC connection implied by spec.
func Connect(spec Spec) (*grpc.ClientConn, error) {
	opts, err := spec.Options()
	if err != nil {
		return nil, errors.Annotate(err, "cannot configure connection")
	}
	return grpc.Dial(spec.Addr, opts...)

}
