package stussy

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"

	"github.com/quinn-caverly/forward-utils/endpointstructs"
	"github.com/quinn-caverly/forward-utils/mongoio"
	"github.com/quinn-caverly/forward-utils/rpcimpls"
)

var dateScrapedGlobal string

func TestScrape(dateScraped string) {
	dateScrapedGlobal = dateScraped
	prods := []endpointstructs.UniqueProductExpanded{}

	products_colors_xml := make(map[string][]string)
	traverseSitemap("https://www.stussy.com/sitemap.xml", &products_colors_xml)

	var collection_identifier string
	visited_ids := make(map[string]bool)
	collection := collectionConfigure(&prods, &products_colors_xml, &visited_ids, &collection_identifier)

	collections := [...]string{"tees"}
	for _, val := range collections {
		collection_identifier = val
		before := len(prods)
		scrapeCollection(collection, "https://www.stussy.com/collections/"+val, &prods)
		log.Println("Collection ", string(collection_identifier), " has been scraped. ", len(prods)-before, " products added.")
	}

    fmt.Println(prods[3].Name, prods[3].Price, prods[3].Brand, prods[3].Id, prods[3].Description, prods[3].ClothingType, prods[3].UrlOnBrandSite)

	log.Println("Complete. Considered ", fmt.Sprint(len(prods)), " products. ")
}

func Scrape(dateScraped string) {
	dateScrapedGlobal = dateScraped

	prods := []endpointstructs.UniqueProductExpanded{}

	products_colors_xml := make(map[string][]string)
	traverseSitemap("https://www.stussy.com/sitemap.xml", &products_colors_xml)

	var collection_identifier string
	visited_ids := make(map[string]bool)
	collection := collectionConfigure(&prods, &products_colors_xml, &visited_ids, &collection_identifier)

	collections := [...]string{"tees", "water-shorts", "sweats", "tops-shirts", "knits", "bottoms", "outerwear", "headwear", "accessories", "eyewear"}
	for _, val := range collections {
		collection_identifier = val
		before := len(prods)
		scrapeCollection(collection, "https://www.stussy.com/collections/"+val, &prods)
		log.Println("Collection ", string(collection_identifier), " has been scraped. ", len(prods)-before, " products added.")
	}

	writeToDbs(prods)

	log.Println("Complete. Considered ", fmt.Sprint(len(prods)), " products. ")
}

func writeToDbs(prods []endpointstructs.UniqueProductExpanded) {

	coll, mongoClient, err := mongoio.CreateConnToBrand("stussy")
	if err != nil {
		log.Fatal("Could not create connection to stussy, ", err)
	}

	client, serviceMethod, err := rpcimpls.ConnectToGoWritePod()
	if err != nil {
		log.Fatal("Could not establish connection to go write pod, ", err)
	}

	for i := range prods {

		colorContainersToWrite, err := mongoio.WriteUPE(&prods[i], coll)
		if err != nil {
			log.Fatal("Error when attempting to write UPE to mongo, ", err)
		}

		for j := range colorContainersToWrite {
			prodContainerForDB := endpointstructs.ProductContainerForWritingToDB{Upi: endpointstructs.UniqueProductIdentifier{Brand: prods[i].Brand, Id: prods[i].Id}, ColorAttr: colorContainersToWrite[j].ColorAttr, ImageURLs: colorContainersToWrite[j].ImageURLs}
			err = client.Call(serviceMethod, prodContainerForDB, nil)
			if err != nil {
				log.Fatal("Error when writing color to go write pod, ", err)
			}
		}
	}

	mongoClient.Disconnect(context.Background())
}

func scrapeCollection(c *colly.Collector, base_url string, prods *[]endpointstructs.UniqueProductExpanded) {
	prev_products_len := -1
	page := 1
	for prev_products_len < len(*prods) {
		prev_products_len = len(*prods)
		c.Visit(base_url + "?page=" + fmt.Sprint(page))
	}
}

// start at the sitemap.xml, then find the xml elem which contains products because this link changes
// then, keys are products and values are the slice which has the links to each different color of product
func traverseSitemap(site_map_url string, products_colors_xml *map[string][]string) {
	original_sitemap_collector := colly.NewCollector()
	products_sitemap_collector := colly.NewCollector()

	original_sitemap_collector.OnXML("//loc", func(e *colly.XMLElement) {
		loc := e.Text
		if strings.Contains(loc, "product") {
			products_sitemap_collector.Visit(e.Text)
		}
	})

	products_sitemap_collector.OnXML("//loc", func(e *colly.XMLElement) {
		loc := e.Text

		exp := regexp.MustCompile(`\/products\/([\w]+)`)
		match := exp.FindString(loc)
		id := regexp.MustCompile(`\/products\/`).ReplaceAllString(match, "")

		cur_colors, exists := (*products_colors_xml)[id]
		if exists == false {
			(*products_colors_xml)[id] = []string{loc}
		} else {
			for _, val := range (*products_colors_xml)[id] {
				if val == loc {
					return
				}
			}
			(*products_colors_xml)[id] = append(cur_colors, loc)
		}
	})

	original_sitemap_collector.Visit(site_map_url)
}
