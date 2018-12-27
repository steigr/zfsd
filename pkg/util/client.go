package util

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var Client *grpc.ClientConn

// NewGRPCConn is a helper wrapper around grpc.Dial.
func newGRPCConn(
	address string,
	serverCAFileName string,
	clientCertFileName string,
	clientKeyFileName string,
) (*grpc.ClientConn, error) {
	endpointAddressTuple := strings.Split(address,"://")
	if serverCAFileName == "" {
		if endpointAddressTuple[0] == "unix" {
			return grpc.Dial(endpointAddressTuple[1],
				grpc.WithInsecure(),
				grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
					return net.DialTimeout("unix", addr, timeout)
				}),
				grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
		} else {
			return grpc.Dial(endpointAddressTuple[1],
				grpc.WithInsecure(),
				grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
		}
	}

	caCert, err := ioutil.ReadFile(serverCAFileName)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cfg := &tls.Config{
		RootCAs: caCertPool,
	}

	if clientCertFileName != "" && clientKeyFileName != "" {
		peerCert, err := tls.LoadX509KeyPair(clientCertFileName, clientKeyFileName)
		if err != nil {
			return nil, err
		}
		cfg.Certificates = []tls.Certificate{peerCert}
	}
	return grpc.Dial(endpointAddressTuple[1],
		grpc.WithTransportCredentials(credentials.NewTLS(cfg)),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
}

func Connect(endpoint string) (err error) {
	if Client == nil {
		Client, err = newGRPCConn(endpoint,"","","")
		if err != nil {
			return err
		}
	}
	return nil
}

func Disconnect() error {
	if Client == nil {
		return nil
	}
	if err := Client.Close(); err != nil {
		return err
	}
	return nil
}