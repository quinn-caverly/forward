package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/quinn-caverly/forward-utils/endpointstructs"
)

func main() {
    fmt.Println("starting")
    client, err := rpc.Dial("tcp", "go-product-builder-service:8080")
	if err != nil {
		log.Fatal("RPC client dial error:", err)
	}

    fmt.Println("connection established")

	var reply endpointstructs.ProductDisplayContainer
    err = client.Call("BuildProduct.Read", endpointstructs.UniqueProductIdentifier{Id: "1994879", Brand: "stussy"}, &reply)
	if err != nil {
		log.Fatal("RPC call error:", err)
	}

    log.Println(len(reply.ColorContainers[0].ImageBytes))

	log.Println("RPC response:", reply)
}
