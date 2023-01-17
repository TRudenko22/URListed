package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Sorter struct {
	Links []string `yaml:"links"`
	Imgs  []string `yaml:"images"`
}

func (s *Sorter) Collect(website string, depth int) {

	collector := colly.NewCollector(colly.MaxDepth(depth))

	// VISITING <website>
	collector.OnRequest(func(r *colly.Request) {
		log.Println("VISITING", r.URL.String())
	})

	// Links
	collector.OnHTML("a[href]", func(h *colly.HTMLElement) {
		ScrapeData("href", &s.Links, true, h)
	})

	// Images
	collector.OnHTML("img", func(h *colly.HTMLElement) {
		ScrapeData("src", &s.Imgs, false, h)
	})

	err := collector.Visit(website)
	if err != nil {
		log.Fatal("Error with link")
	}
}

func ScrapeData(attribute string, list *[]string, recurse bool, element *colly.HTMLElement) {

	var prefix string
	link := element.Attr(attribute)

	if len(link) == 0 {
		log.Println("BROKEN LINK")
	} else {
		prefix = string(link[0])
	}

	if len(link) > 7 && (prefix != "/" && prefix != "#") {
		*list = append(*list, link)
		if recurse {
			element.Request.Visit(link)
		}
	}
}

func WriteConfig(file_name string, sorter Sorter) {
	yml_data, err := yaml.Marshal(&sorter)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(file_name, yml_data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	app := &cli.App{
		Name:  "pull",
		Usage: "Pulls down all links from url provided",
		Action: func(c *cli.Context) (err error) {
			website_bits := Sorter{}
			website_bits.Collect(c.Args().Get(1), 1)

			WriteConfig("websites", website_bits)

			fmt.Println("Amount of Links:", len(website_bits.Links))
			fmt.Println("Amount of Images:", len(website_bits.Imgs))
			return
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
