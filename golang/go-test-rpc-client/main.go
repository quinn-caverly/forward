package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type Color struct {
	Color    string
	Pic_urls []string
}

type Args struct {
	Collection   string
	Date_scraped string
	Id           string
	Colors       []Color
}

func main() {
    fmt.Println("starting")
    client, err := rpc.Dial("tcp", "go-write-service:80")
	if err != nil {
		log.Fatal("RPC client dial error:", err)
	}

    fmt.Println("connection established")

	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal("Was not able to load EST timezone: ", err)
	}
	date_scraped := time.Now().In(est).Format("2006-01-02")

	collection := "stussy"

	color := Color{Color: "gray", Pic_urls: []string{
		"https://cdn.shopify.com/s/files/1/0087/6193/3920/products/118469_GH22_1_2b3ff92b-d809-46d3-9884-84197a8b3890.jpg?v=1683139297",
		"https://cdn.shopify.com/s/files/1/0087/6193/3920/products/118469_GH22_2_f6095ed6-f21b-4b78-b8bf-a556328efa83.jpg?v=1683139297"}}

        args := Args{Collection: collection, Date_scraped: date_scraped, Id: "118469", Colors: []Color{color}}

	var reply int
	err = client.Call("WriteService.Write", args, &reply)
	if err != nil {
		log.Fatal("RPC call error:", err)
	}

	log.Println("RPC response:", reply)
}
