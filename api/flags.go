package api

import flag "github.com/spf13/pflag"

var (
	APIServiceFlags = flag.NewFlagSet("api", flag.ExitOnError)

	host     string
	port     int
	grpcPort int
)

func init() {
	APIServiceFlags.StringVar(&host,
		"host",
		"localhost",
		"API service listening address",
	)

	APIServiceFlags.IntVar(&port,
		"port",
		8000,
		"API service listening port",
	)

	APIServiceFlags.IntVar(&grpcPort,
		"grpcport",
		9000,
		"gRPC service listening port",
	)
}
