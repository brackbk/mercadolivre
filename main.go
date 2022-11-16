package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gosimple/slug"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gitlab.com/eiprice/crawlers/mercadolivre/crawler"
	"gitlab.com/eiprice/crawlers/mercadolivre/utils"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func run(state string, city string, cep string, departament string, category string, subCategory string, scan string, fileurls string, collection string) error {
	var api crawler.Crawler
	var err error
	ceps := strings.Split(cep, ",")

	for _, cp := range ceps {
		api = crawler.Crawler{departament, category, subCategory, scan, state, city, cp, fileurls, collection}
		log.Printf("Start Crawler from", state)

		err = api.Start()
		if err != nil {
			continue
		}
	}

	return nil
}

func main() {
	var cep string
	var state string
	var city string
	var subCategory string
	var category string
	var departament string
	var scan string
	var drop string
	var file string
	var fileurls string
	var collection string
	var sc int64
	flag.StringVar(&cep, "cep", "", "Set Cep")
	flag.StringVar(&state, "state", "Sao Paulo", "Set State")
	flag.StringVar(&city, "city", "", "Set City")
	flag.StringVar(&departament, "departament", "", "Set Departament")
	flag.StringVar(&category, "category", "", "Set Category")
	flag.StringVar(&subCategory, "subCategory", "", "Set SubCategory")
	flag.StringVar(&scan, "scan", "1", "Set Scan")
	flag.StringVar(&drop, "drop", "", "Set Drop")
	flag.StringVar(&file, "file", "", "Set File")
	flag.StringVar(&fileurls, "fileurls", "", "Set File Urls")
	flag.StringVar(&collection, "collection", "", "Set Collection")
	flag.Parse()

	if drop == "all" {

		if collection != "" {
			drop := crawler.Crawler{departament, category, subCategory, scan, state, city, cep, fileurls, collection}
			drop.Drop()
			fmt.Println(`droped:`, collection)
		} else {
			fmt.Println(`Informe o collection`)
		}
		os.Exit(0)
	}

	if city != "" {
		city = slug.Make(city)
	}
	if state != "" {
		state = slug.Make(state)
	}

	fmt.Println(state)
	//Crawler get data from mercadolivre website

	if file != "" {
		Coords := utils.Read(file)

		for _, item := range Coords {
			sc = sc + 1
			err := run(item.State, item.City, cep, departament, category, subCategory, fmt.Sprint(sc), "", collection)

			if err != nil {
				continue
			}
		}
		// os.Exit(0)
	} else {
		err := run(state, city, cep, departament, category, subCategory, scan, fileurls, collection)
		fmt.Println(`Error:`, err)
	}

}
