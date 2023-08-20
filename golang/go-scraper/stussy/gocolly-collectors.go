package stussy

import (
	"fmt"
	"log"
	"regexp"
    "strings"

	"github.com/gocolly/colly"
	"github.com/json-iterator/go"

	"github.com/quinn-caverly/forward-utils/endpointstructs"
)

// the collector for the collection pages. ex) sweats, shorts, tees, ...
func collectionConfigure(prods *[]endpointstructs.UniqueProductExpanded, products_colors_xml *map[string][]string, visited_ids *map[string]bool, collection_identifier *string) *colly.Collector {
	collection := colly.NewCollector()

	collection.OnHTML("div[class]", func(e *colly.HTMLElement) {
		if e.Attr("class") != "product-card" {
			return
		}

		dom := e.DOM
		href, exists := dom.Find("a[href]").Attr("href")
		if exists == false {
			log.Println("href which corresponds to product's page did not exist. This should happen once per collection. The collection is: ", fmt.Sprintf(*collection_identifier))
			return
		}
		var Price, Name, Description string
		product := productConfigure(&Price, &Name, &Description)
		UrlOnBrandSite := "https://www.stussy.com" + href
		product.Visit(UrlOnBrandSite)

		exp := regexp.MustCompile(`\/products\/([\w]+)`)
		match := exp.FindString(href)
		Id := regexp.MustCompile(`\/products\/`).ReplaceAllString(match, "")

		_, in_visited_ids := (*visited_ids)[Id]
		if in_visited_ids == false {
			imgUrls := []string{}
			colorCollector := colorConfigure(&imgUrls)
			colorUrls := (*products_colors_xml)[Id]
			for i := range colorUrls {
				colorCollector.Visit(colorUrls[i])
			}
			(*visited_ids)[Id] = true

			URLColorContainers, err := parseImgUrlsToURLColorContainers(imgUrls)
			if err != nil {
				log.Println("Error when generating URLColorContainers, ", err)
			}

			*prods = append(*prods, endpointstructs.UniqueProductExpanded{
				Id:                 Id,
				Brand:              "stussy",
				Name:               Name,
				Price:              Price,
				Description:        Description,
				UrlOnBrandSite:     UrlOnBrandSite,
				ClothingType:       *collection_identifier,
				URLColorContainers: URLColorContainers,
			})
		}
	})

	return collection
}

func productConfigure(price, name, description *string) *colly.Collector {
	product := colly.NewCollector()

	product.OnHTML("div[class='shopify-section']", func(e *colly.HTMLElement) {
		script := e.DOM.Find("script[type='application/ld+json']")
		jsonBytes := []byte(script.Text())

		*price = jsoniter.Get(jsonBytes, "offers", 0, "price").ToString()
		*name = jsoniter.Get(jsonBytes, "name").ToString()
		*description = jsoniter.Get(jsonBytes, "description").ToString()
	})

	return product
}

func colorConfigure(img_urls *[]string) *colly.Collector {
	color_collector := colly.NewCollector()

	color_collector.OnHTML("img[loading='lazy']", func(e *colly.HTMLElement) {
		srcset := e.Attr("srcset")
		srcs := strings.Split(srcset, ",")

		for i := range srcs {
			//TODO clean up this code, I don't know why the src turns out to be index 4
			if strings.Contains(srcs[i], "1440w") {
				src := "https:" + strings.Split(srcs[i], " ")[4]

				for _, val := range *img_urls {
					if val == src {
						return
					}
				}

				*img_urls = append(*img_urls, src)
			}
		}
	})

	return color_collector
}
