package main

import (
	"fmt"
	"log"
	"net/rpc"
    "context"

	"github.com/quinn-caverly/forward-utils/endpointstructs"
	"github.com/quinn-caverly/forward-utils/mongoio"
	"github.com/quinn-caverly/forward-utils/rpcimpls"
)

func main() {
	BuildProduct := new(BuildProduct)
	rpc.Register(BuildProduct)

	listener, err := rpcimpls.CreatePodListener()
	if err != nil {
		log.Fatal("Listener error:", err)
	}

	rpc.Accept(listener)
}

type BuildProduct struct{}

// reply should be an initially empty Product
func (s *BuildProduct) Read(args endpointstructs.UniqueProductIdentifier, reply *endpointstructs.ProductDisplayContainer) error {

	coll, client, err := mongoio.CreateConnToBrand(args.Brand)
	if err != nil {
		return fmt.Errorf("Error when trying to connect to brand collection, %w", err)
	}
    defer client.Disconnect(context.TODO())

	upe, err := mongoio.ReadUPE(args.Id, coll)
	if err != nil {
		return fmt.Errorf("Error when reading from collection via id, %w", err)
	}

    client.Disconnect(context.TODO())

	readPodClient, serviceMethod, err := rpcimpls.ConnectToGoReadPod()
	if err != nil {
		return fmt.Errorf("Error when connecting to go read pod. %w", err)
	}

	colorContainers := []endpointstructs.ColorContainer{}
	for _, urlColorCont := range upe.URLColorContainers {
        pcff := endpointstructs.ProductContainerForFrontend{}
        err = readPodClient.Call(serviceMethod, endpointstructs.UniqueColorIdentifier{Upi: args, ColorAttr: urlColorCont.ColorAttr}, &pcff)
        if err != nil {
            return fmt.Errorf("Error when receiving images from read pod client. %w", err)
        }

        log.Println(pcff)

        colorContainers = append(colorContainers, endpointstructs.ColorContainer{ColorAttr: urlColorCont.ColorAttr, ImageBytes: pcff.ImageBytes})
	}

    log.Println(colorContainers)

    // UniqueProductExpanded is converted to UniqueProduct here so that we do not need to send the original url info to frontend, struct transferred is smaller
    uniqueProduct := endpointstructs.UniqueProduct {
        Id: upe.Id,
        Brand: upe.Brand,
        Name: upe.Name,
        UrlOnBrandSite: upe.UrlOnBrandSite,
        Price: upe.Price,
        Description: upe.Description,
        ClothingType: upe.ClothingType,
    }

    reply.Up = uniqueProduct
    reply.ColorContainers = colorContainers
	return nil
}
