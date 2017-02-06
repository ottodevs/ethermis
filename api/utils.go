package api

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/philips/grpc-gateway-example/insecure"
)

var (
	demoKeyPair  *tls.Certificate
	demoCertPool *x509.CertPool
)

func init() {
	var err error
	pair, err := tls.X509KeyPair([]byte(insecure.Cert), []byte(insecure.Key))
	if err != nil {
		panic(err)
	}
	demoKeyPair = &pair
	demoCertPool = x509.NewCertPool()
	ok := demoCertPool.AppendCertsFromPEM([]byte(insecure.Cert))
	if !ok {
		panic("bad certs")
	}
}
