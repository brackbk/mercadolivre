package crawler

import (
	"fmt"
	"log"

	"gitlab.com/eiprice/crawlers/mercadolivre/repositories"
	"gitlab.com/eiprice/crawlers/mercadolivre/utils"
)

type Crawler struct {
	Departament string
	Category    string
	SubCategory string
	Scan        string
	State       string
	City        string
	Cep         string
	Fileurls    string
	Collection  string
}

func (craw *Crawler) Start() error {

	context := utils.NewContext()
	defer context.Close()

	departamentCraw := DepartamentCrawlerInit{craw.Departament, craw.Category, craw.Scan}
	categoryCraw := CategoryCrawlerInit{craw.Category, craw.SubCategory, craw.Scan}
	productRepo := repositories.ProductRepositoryDb{context, craw.Collection}
	productCraw := ProductCrawlerInit{productRepo, craw.State, craw.City, craw.Cep, craw.Scan}
	if craw.Fileurls != "" && craw.Departament != "supermercado" {

		urls := utils.ReadUrls(craw.Fileurls)

		for _, u := range urls {
			fmt.Println(u.Url)
			productCraw.GetByPruducList(u.Url)
		}
	} else if craw.Fileurls != "" && craw.Departament == "supermercado" {

		urls := utils.ReadUrls(craw.Fileurls)

		for _, u := range urls {
			fmt.Println(u.Url)
			productCraw.GetByPage(u.Url)
		}

	} else {
		departaments, _ := departamentCraw.GetData()
		for _, item := range departaments {

			fmt.Println("Get Categories from: ", item.Name)

			subCategories, categories, err := categoryCraw.GetData(item)
			productCraw.GetData(subCategories, categories)

			if err != nil {
				log.Println("Error Get Categories Data ", err)
				continue
			}
		}
	}

	return nil

}
func (craw *Crawler) Drop() {

	context := utils.NewContext()
	defer context.Close()
	productRepo := repositories.ProductRepositoryDb{context, craw.Collection}
	productRepo.Drop()

}
