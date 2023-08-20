package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/quinn-caverly/forward-utils/endpointstructs"
	"github.com/quinn-caverly/forward-utils/rpcimpls"
)

func main() {
	ReadService := new(ReadService)
	rpc.Register(ReadService)

	listener, err := rpcimpls.CreatePodListener()
	if err != nil {
		log.Fatal("Error when creating listener via rpc impls, ", err)
	}

	log.Println("Go read pod running on port 8080: ")

	rpc.Accept(listener)
}

type ReadService struct{}

// reply must initially be an empty Color slice
func (s *ReadService) Read(args endpointstructs.UniqueColorIdentifier, reply *endpointstructs.ProductContainerForFrontend) error {
	dir_path := "/data/" + args.Upi.Brand + "/" + args.ColorAttr.DateScraped + "/" + args.Upi.Id + "/" + args.ColorAttr.ColorName

	imgs, err := take_images_from_color_dir(dir_path, args.ColorAttr.ColorName)
	if err != nil {
		log.Println(err)
		return err
	}

    bytes_double_slice, err := convert_imgs_to_byte_slices(imgs)
    if err != nil {
        log.Println(err)
        return err
    }

	*reply = endpointstructs.ProductContainerForFrontend{Upi: args.Upi, ColorAttr: args.ColorAttr, ImageBytes: bytes_double_slice}
	return nil
}

func convert_imgs_to_byte_slices(imgs []image.Image) ([][]byte, error) {
    bytes_slice := [][]byte{}

    for i := range imgs {
        buf := new(bytes.Buffer)
        err := jpeg.Encode(buf, imgs[i], nil)
        if err != nil {
           return nil, err
        }

        bytes_slice = append(bytes_slice, buf.Bytes())
    }

    return bytes_slice, nil
}

func take_images_from_color_dir(color_dir_path string, color_name string) ([]image.Image, error) {
	dir, err := os.Open(color_dir_path)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Error when attempting to open color directory. %w", err)
	}

	img_files, err := dir.Readdir(0)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Error when attempting to parse children of color directory. %w", err)
	}

	imgs := []image.Image{}
	for _, entry := range img_files {
		img_path := filepath.Join(color_dir_path, entry.Name())

		img, err := read_to_img(img_path)
		if err != nil {
			return nil, err
		}

		imgs = append(imgs, img)
	}
	return imgs, nil
}

func read_to_img(filepath string) (image.Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Error when attempting to open img path as a file.  %w", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Error when attempting to convert image file into image.Image  %w", err)
	}
	return img, nil
}
