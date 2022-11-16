package crawler

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/eiprice/crawlers/mercadolivre/domain"
)

type DepartamentCrawler interface {
	GetData() ([]*domain.Departament, error)
}
type DepartamentCrawlerInit struct {
	Departament string
	Category    string
	List        string
}

func (craw *DepartamentCrawlerInit) GetData() ([]*domain.Departament, error) {

	var departament *domain.Departament
	var departaments []*domain.Departament
	var title string
	// var headers []utils.Header

	doc, err := goquery.NewDocument("https://www.mercadolivre.com.br/categorias#menu=categories")
	if err != nil {
		return nil, err
	}

	doc.Find("h2.categories__title").Each(func(i int, s *goquery.Selection) {

		li := s.Find("a")
		url, _ := li.Attr("href")
		title = li.Text()

		departament, err = domain.NewDepartament(
			title,
			url,
			craw.List,
		)

		if craw.Departament != "" {
			if strings.Contains(title, craw.Departament) {
				departaments = append(departaments, departament)
			}
		} else {
			departaments = append(departaments, departament)
		}

	})
	fmt.Println("Total Departaments inserted from Store "+title+": ", len(departaments))

	return departaments, nil
}

func (craw *DepartamentCrawlerInit) GetDataByUrl(url string) ([]*domain.Departament, error) {

	var departament *domain.Departament
	var departaments []*domain.Departament
	var title string
	// var headers []utils.Header

	departament, err := domain.NewDepartament(
		craw.Category,
		url,
		craw.List,
	)
	fmt.Println(departament)
	if err != nil {
		return nil, err
	}

	departaments = append(departaments, departament)

	fmt.Println("Total Departaments inserted from Store "+title+": ", len(departaments))

	return departaments, nil
}
