package api

import (
	"fmt"

	"github.com/alanchchen/ethermis/api/ethereum"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func DeployContract(cmd *cobra.Command, args []string) {
	var opts []grpc.DialOption
	creds := credentials.NewClientTLSFromCert(demoCertPool, "localhost:10000")
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := ethereum.NewEthereumClient(conn)

	msg, err := client.Deploy(context.Background(), &ethereum.CompiledContract{
		Abi:  args[0],
		Code: args[1],
	})

	if err != nil {
		cmd.Println(err)
		return
	}

	cmd.Println(msg)
}
