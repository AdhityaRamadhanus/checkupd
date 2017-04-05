package commands

import (
	"context"
	"log"

	checkupservice "github.com/AdhityaRamadhanus/checkupd/grpc/service"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func listEndpoint(cliContext *cli.Context) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(cliContext.String("host"), grpc.WithInsecure())
	if err != nil {
		log.Println("could not connect to grpc server", err)
	}
	defer conn.Close()
	c := checkupservice.NewCheckupClient(conn)

	r, err := c.ListEndpoint(context.Background(), &checkupservice.ListEndpointRequest{Check: false})

	if err != nil {
		log.Println("failed to get list of endpoints", err)
	}
	for _, endpoint := range r.Endpoints {
		log.Println(endpoint.Name, " ", endpoint.Url, " ", endpoint.Status)
	}
}
