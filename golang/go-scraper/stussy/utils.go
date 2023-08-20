package stussy

import (
	"log"
	"regexp"

	"github.com/quinn-caverly/forward-utils/endpointstructs"
)

func parseImgUrlsToURLColorContainers(imgUrls []string) ([]endpointstructs.URLColorContainer, error) {
	colorsMap := map[string]endpointstructs.URLColorContainer{}
	for i := range imgUrls {
		//TODO static Regexp grab, not sure how this could be less fragile
		re := regexp.MustCompile(`_([a-zA-Z0-9]+)_`)
		matches := re.FindAllStringSubmatch(imgUrls[i], -1)

		// in this case, it will usually be a special product, like keychain and unnecessary
		// will only be 0 if it does not follow the schema for products
		if len(matches) == 0 {
			log.Println("Product did not follow schema, probably unnecessary special product.")
			return nil, nil
		}
		color := matches[0][1]

		colorContainer, ok := colorsMap[color]
		if ok == true {
			colorContainer.ImageURLs = append(colorContainer.ImageURLs, imgUrls[i])
			colorsMap[color] = colorContainer
		} else {
			colorsMap[color] = endpointstructs.URLColorContainer{ColorAttr: endpointstructs.ColorAttr{
				ColorName:   color,
				DateScraped: dateScrapedGlobal,
			},
				ImageURLs: []string{imgUrls[i]}}
		}
	}

	colorContainers := []endpointstructs.URLColorContainer{}
	for _, value := range colorsMap {
		colorContainers = append(colorContainers, value)
	}

	return colorContainers, nil
}
