package crawler

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/eiprice/crawlers/mercadolivre/domain"
)

type CategoryCrawler interface {
	GetData(departament *domain.Departament) ([]*domain.SubCategory, []*domain.Category, error)
	GetSubCategory(category *domain.Category, subCategorySection *goquery.Selection) ([]*domain.SubCategory, error)
}
type CategoryCrawlerInit struct {
	Category    string
	SubCategory string
	List        string
}

func (craw *CategoryCrawlerInit) GetData(departament *domain.Departament) ([]*domain.SubCategory, []*domain.Category, error) {

	var category *domain.Category
	var categories []*domain.Category
	var subCategory []*domain.SubCategory
	var subCategories []*domain.SubCategory
	var title string
	// var headers []utils.Header

	doc, err := goquery.NewDocument(departament.Url)
	if err != nil {
		return nil, nil, err
	}

	doc.Find(".desktop__view-child").Each(func(i int, s *goquery.Selection) {

		li := s.Find("a")
		url, _ := li.Attr("href")
		title_find := li.Find(".category-list__permanlink-custom")
		title = title_find.Text()

		category, err = domain.NewCategory(
			departament.ID,
			departament.Name,
			title,
			url,
			craw.List,
		)

		if craw.Category != "" {
			if strings.Contains(title, craw.Category) {
				categories = append(categories, category)
				subCategory, err = craw.GetSubCategory(category, s)
				subCategories = append(subCategories, subCategory...)

			}

		} else {
			categories = append(categories, category)
			subCategory, err = craw.GetSubCategory(category, s)
			subCategories = append(subCategories, subCategory...)

		}

	})

	fmt.Println("Total Categories inserted from Departament "+departament.Name+": ", len(categories))

	return subCategories, categories, nil
}

func (craw *CategoryCrawlerInit) GetSubCategory(category *domain.Category, subCategorySection *goquery.Selection) ([]*domain.SubCategory, error) {

	var subCategory *domain.SubCategory
	var subCategories []*domain.SubCategory
	var title string

	subCategorySection.Find(".category-list__item").Each(func(i int, s *goquery.Selection) {

		li := s.Find("a")
		url, _ := li.Attr("href")
		title_find := li.Find(".category-list__permanlink-title")
		title = title_find.Text()

		subCategory, _ = domain.NewSubCategory(
			category.ID,
			category.DepartamentID,
			title,
			url,
			craw.List,
		)

		if craw.SubCategory != "" {
			if strings.Contains(title, craw.SubCategory) {
				subCategories = append(subCategories, subCategory)
			}
		} else {
			subCategories = append(subCategories, subCategory)
		}

	})
	fmt.Println("Total Sub Categories inserted from Category "+category.Name+": ", len(subCategories))

	return subCategories, nil
}
