package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/quinn-caverly/forward-utils/endpointstructs"
	"github.com/quinn-caverly/forward-utils/rpcimpls"
)

func main() {
	readPodClient, serviceMethod, err := rpcimpls.ConnectToGoReadPod()
	if err != nil {
		log.Fatal("Error when connecting to go read pod. %w", err)
	}

	pcff := endpointstructs.ProductContainerForFrontend{}
	err = readPodClient.Call(serviceMethod, endpointstructs.UniqueColorIdentifier{Upi: endpointstructs.UniqueProductIdentifier{Brand: "stussy", Id: "1994879"},
		ColorAttr: endpointstructs.ColorAttr{ColorName: "BLAC", DateScraped: "2023-08-15"}}, &pcff)

	if err != nil {
		log.Println("Error when receiving info from go read pod. ", err)
		unwrapped := errors.Unwrap(err)
		log.Fatal(unwrapped)
	}

	fmt.Println(len(pcff.ImageBytes))
	for _, val := range pcff.ImageBytes {
		fmt.Println(val)
	}
}
