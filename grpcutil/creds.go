package grpcutil

import (
	"context"
	"crypto/x509"

	"github.com/juju/errors"
	"google.golang.org/grpc/credentials"
)

// rpcCredentials is an implementation of credentials.PerRPCCredentials
// that holds the Slash GraphQL API token.
type rpcCredentials struct {
	apiToken string
}

// GetRequestMetadata sets the value for "authorization" key.
func (rc rpcCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": rc.apiToken,
	}, nil
}

// RequireTransportSecurity should be true as we want to have the token
// encrypted over the wire.
func (rc rpcCredentials) RequireTransportSecurity() bool {
	return true
}

// transportCredentials returns credentials that check we're talking to a
// certified server.
func transportCredentials() (credentials.TransportCredentials, error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Annotate(err, "cannot get system certs")
	}
	return credentials.NewClientTLSFromCert(pool, ""), nil
}
