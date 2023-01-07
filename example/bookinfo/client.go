package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"

	kitexreflect "github.com/jxskiss/kitex-reflect"
)

var destService = "cwg.bookinfo.details"

func main() {
	serverAddr := flag.String("server-addr", "localhost:8888", "the server address (host:port)")
	flag.Parse()

	kcOptions := []client.Option{
		client.WithHostPorts(*serverAddr),
	}

	ctx := context.Background()
	prov, err := kitexreflect.NewDescriptorProvider(ctx, destService, kcOptions)
	if err != nil {
		log.Fatal(err)
	}
	g, err := generic.MapThriftGenericForJSON(prov)
	if err != nil {
		log.Fatal(err)
	}
	gCli, err := genericclient.NewClient(destService, g, kcOptions...)
	if err != nil {
		log.Fatal(err)
	}

	req := map[string]interface{}{
		"ID": "12345",
	}
	resp, err := gCli.GenericCall(ctx, "GetProduct", req)
	if err != nil {
		log.Fatal(err)
	}

	buf, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))
}
