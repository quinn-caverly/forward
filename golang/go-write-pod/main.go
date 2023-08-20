package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/quinn-caverly/forward-utils/endpointstructs"
	"github.com/quinn-caverly/forward-utils/rpcimpls"
)

func main() {
	WriteService := new(WriteService)
	rpc.Register(WriteService)

	listener, err := rpcimpls.CreatePodListener()
	if err != nil {
		log.Fatal("Error when creating listener, ", err)
	}

	rpc.Accept(listener)
}

type WriteService struct{}

func (s *WriteService) Write(args endpointstructs.ProductContainerForWritingToDB, reply *int) error {
	dirPath := "/data/" + args.Upi.Brand + "/" + args.ColorAttr.DateScraped

	for i, url := range args.ImageURLs {
		err := write_image_to_loc(url, dirPath+"/"+args.Upi.Id +"/"+args.ColorAttr.ColorName+"/"+fmt.Sprint(i)+".jpg")
		if err != nil {
			return err
		}
	}

	return nil
}

func write_image_to_loc(url, fileloc string) error {
	err := os.MkdirAll(filepath.Dir(fileloc), 0755)
	if err != nil {
		return fmt.Errorf("Failed to create parent directories. %w", err)
	}

	file, err := os.Create(fileloc)
	if err != nil {
		return fmt.Errorf("Failed to create file. %w", err)
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
        fmt.Println(err)
		return fmt.Errorf("Failed to get response from url. %w", err)
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("Failed to write image to file. %w", err)
	}
	return nil
}
