package crawler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/eiprice/crawlers/mercadolivre/domain"
	"gitlab.com/eiprice/crawlers/mercadolivre/repositories"
	"gitlab.com/eiprice/crawlers/mercadolivre/utils"
)

type ProductCrawler interface {
	GetByPruducList(url string)
	PrepareUrl(url string) string
	GetData(subCategories []*domain.SubCategory, categories []*domain.Category)
	GetByPage(subCategories []*domain.SubCategory, categories []*domain.Category)
	GetByUrlPage(url string, size int) ([]*domain.Product, error)
	GetByCategory(category *domain.Category, size int)
	getShipment(delivery string) ([]domain.Shipment, error)
	getSeller(content *goquery.Document)
	getOtherSellers(content *goquery.Document) []domain.Seller
}

type ProductCrawlerInit struct {
	ProductRepository repositories.ProductRepository
	State             string
	City              string
	Cep               string
	List              string
}

func (craw *ProductCrawlerInit) GetByPruducList(url string) {
	var sku string
	var name string
	var brand string
	var seller []domain.Seller
	var level []domain.Level
	var attributes []domain.Attributes
	var image []domain.Image
	var products []*domain.Product

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("cookie", `cp=`+craw.Cep+`%7C1622117690815`)
	resp, err := client.Do(request)

	if err != nil {
		return
	}

	page_product, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return
	}

	name = getName(page_product)
	brand = getBrand(page_product)
	sku = getSku(page_product)
	image = getImage(page_product)
	seller = craw.getSeller(page_product)

	fmt.Println(`get item: `, name)

	level = getLevels(page_product)
	attributes = getAttributes(page_product)
	fmt.Println(url)
	obj, _ := domain.NewProduct(
		utils.StringNotNull(name),
		0,
		sku,
		sku,
		sku,
		level,
		brand,
		attributes,
		"",
		image,
		url,
		domain.ReviewInfo{},
		1,
		"Mobile",
		seller,
	)

	products = append(products, obj)
	_, err = craw.ProductRepository.Insert(obj)

	if err != nil {
		log.Fatalf("Error to insert")
	}
}

func (craw *ProductCrawlerInit) PrepareUrl(url string) string {
	var best_seller string = `/_Loja_all_BestSellers_YES_SHIPPING*ORIGIN_10215068#applied_filter_id=power_seller&applied_filter_name=Filtro+MercadoLíderes&applied_filter_order=15&applied_value_id=yes&applied_value_name=Melhores+vendedores&applied_value_order=1&is_custom=false#`
	var city_and_state string
	city_and_state = `-em-` + craw.State + "/#"

	if craw.City != "" {
		city_and_state = `-em-` + craw.City + `-` + craw.State + "/#"
	}

	if craw.Cep != "" {
		url = strings.Replace(url, "/#", best_seller, -1)
	} else {
		url = strings.Replace(url, "/#", city_and_state, -1)
	}

	return url
}
func (craw *ProductCrawlerInit) GetByPage(url string) {
	var SIZE_PAGE int = 48
	var size int = 0

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	request.Header.Set("cookie", `cp=`+craw.Cep+`%7C1622117690815`)
	resp, err := client.Do(request)

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return
	}

	total_text := doc.Find(".ui-search-search-result__quantity-results").Text()

	t := strings.Replace(total_text, " resultados", "", -1)
	t = strings.Replace(t, ".", "", -1)

	total, _ := strconv.Atoi(t)
	fmt.Println("Total Products to insert from : "+url+": ", t)
	for size < total {
		fmt.Println(`Page size inserted `+url+`: `, size)

		_, _ = craw.GetByUrlPage(url, size)

		size = size + SIZE_PAGE

	}
}

func (craw *ProductCrawlerInit) GetData(subCategories []*domain.SubCategory, categories []*domain.Category) {
	var headers []utils.Header
	headers = append(headers, utils.Header{
		Key:   "cookie",
		Value: `cp=` + craw.Cep + `%7C1622117690815`,
	})

	if len(subCategories) > 0 {
		for _, item := range subCategories {
			var SIZE_PAGE int = 48
			var size int = 0

			client := &http.Client{
				Timeout: 30 * time.Second,
			}

			url := craw.PrepareUrl(item.Url)

			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err)
			}

			request.Header.Set("cookie", `cp=`+craw.Cep+`%7C1622117690815`)
			resp, err := client.Do(request)

			doc, err := goquery.NewDocumentFromReader(resp.Body)

			if err != nil {
				continue
			}

			total_text := doc.Find(".ui-search-search-result__quantity-results").Text()

			t := strings.Replace(total_text, " resultados", "", -1)
			t = strings.Replace(t, ".", "", -1)

			total, _ := strconv.Atoi(t)
			fmt.Println("Total Products to insert from : "+item.Name+": ", t)
			for size < total {
				fmt.Println(`Page size inserted `+item.Name+`: `, size)

				_, _ = craw.GetBySubCategory(item, size)

				size = size + SIZE_PAGE

			}
		}
	} else {

		for _, item := range categories {
			var SIZE_PAGE int = 48
			var size int = 0

			client := &http.Client{
				Timeout: 30 * time.Second,
			}

			url := craw.PrepareUrl(item.Url)

			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err)
			}

			request.Header.Set("cookie", `cp=`+craw.Cep+`%7C1622117690815`)
			resp, err := client.Do(request)

			doc, err := goquery.NewDocumentFromReader(resp.Body)

			if err != nil {
				continue
			}

			total_text := doc.Find(".ui-search-search-result__quantity-results").Text()

			t := strings.Replace(total_text, " resultados", "", -1)
			t = strings.Replace(t, ".", "", -1)

			total, _ := strconv.Atoi(t)
			fmt.Println("Total Products to insert from : "+item.Name+": ", t)
			for size < total {
				fmt.Println(`Page size inserted `+item.Name+`: `, size)

				_, _ = craw.GetByCategory(item, size)

				size = size + SIZE_PAGE

			}
		}
	}
}

func (craw *ProductCrawlerInit) GetByCategory(category *domain.Category, size int) ([]*domain.Product, error) {
	var sku string
	var name string
	var brand string
	var seller []domain.Seller
	var level []domain.Level
	var attributes []domain.Attributes
	var image []domain.Image
	var products []*domain.Product
	url := category.Url

	if size > 0 {
		url = category.Url + `/_Desde_` + strconv.Itoa(size)
	}

	doc, err := goquery.NewDocument(url)

	if err != nil {
		return nil, err
	}

	doc.Find(".ui-search-result__image").Each(func(i int, s *goquery.Selection) {

		tag_url := s.Find("a")

		url, _ := tag_url.Attr("href")
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}
		request.Header.Set("cookie", `cp=`+craw.Cep+`%7C1622117690815`)
		resp, _ := client.Do(request)

		page_product, _ := goquery.NewDocumentFromReader(resp.Body)

		name = getName(page_product)
		brand = getBrand(page_product)
		sku = getSku(page_product)
		image = getImage(page_product)
		seller = craw.getSeller(page_product)

		fmt.Println(`get item: `, name)

		level = getLevels(page_product)
		attributes = getAttributes(page_product)
		fmt.Println(url)
		obj, _ := domain.NewProduct(
			utils.StringNotNull(name),
			0,
			sku,
			sku,
			sku,
			level,
			brand,
			attributes,
			"",
			image,
			url,
			domain.ReviewInfo{},
			1,
			"Mobile",
			seller,
		)

		products = append(products, obj)
		_, err = craw.ProductRepository.Insert(obj)

		if err != nil {
			log.Fatalf("Error to insert")
		}

	})

	return products, nil
}

func (craw *ProductCrawlerInit) GetByUrlPage(url string, size int) ([]*domain.Product, error) {
	var sku string
	var name string
	var brand string
	var seller []domain.Seller
	var level []domain.Level
	var attributes []domain.Attributes
	var image []domain.Image
	var products []*domain.Product

	if size > 0 {
		url = url + `/_Desde_` + strconv.Itoa(size)
	}

	doc, err := goquery.NewDocument(url)

	if err != nil {
		return nil, err
	}

	doc.Find(".ui-search-result__image").Each(func(i int, s *goquery.Selection) {

		tag_url := s.Find("a")

		url, _ := tag_url.Attr("href")
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}
		request.Header.Set("cookie", `cp=`+craw.Cep+`%7C1622117690815`)
		resp, _ := client.Do(request)

		page_product, _ := goquery.NewDocumentFromReader(resp.Body)

		name = getName(page_product)
		brand = getBrand(page_product)
		sku = getSku(page_product)
		image = getImage(page_product)
		seller = craw.getSeller(page_product)

		fmt.Println(`get item: `, name)

		level = getLevels(page_product)
		attributes = getAttributes(page_product)
		fmt.Println(url)
		obj, _ := domain.NewProduct(
			utils.StringNotNull(name),
			0,
			sku,
			sku,
			sku,
			level,
			brand,
			attributes,
			"",
			image,
			url,
			domain.ReviewInfo{},
			1,
			"Mobile",
			seller,
		)

		products = append(products, obj)
		_, err = craw.ProductRepository.Insert(obj)

		if err != nil {
			log.Fatalf("Error to insert")
		}

	})

	return products, nil
}

func (craw *ProductCrawlerInit) GetBySubCategory(subcategory *domain.SubCategory, size int) ([]*domain.Product, error) {
	var sku string
	var name string
	var brand string
	var seller []domain.Seller
	var level []domain.Level
	var attributes []domain.Attributes
	var image []domain.Image
	var products []*domain.Product
	url := subcategory.Url

	if size > 0 {
		url = subcategory.Url + `/_Desde_` + strconv.Itoa(size)
	}

	doc, err := goquery.NewDocument(url)

	if err != nil {
		return nil, err
	}

	doc.Find(".ui-search-result__image").Each(func(i int, s *goquery.Selection) {

		tag_url := s.Find("a")

		url, _ := tag_url.Attr("href")
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}
		request.Header.Set("cookie", `cp=`+craw.Cep+`%7C1622117690815`)
		resp, _ := client.Do(request)

		page_product, _ := goquery.NewDocumentFromReader(resp.Body)

		name = getName(page_product)
		brand = getBrand(page_product)
		sku = getSku(page_product)
		image = getImage(page_product)
		seller = craw.getSeller(page_product)

		fmt.Println(`get item: `, name)

		level = getLevels(page_product)
		attributes = getAttributes(page_product)
		fmt.Println(url)
		obj, _ := domain.NewProduct(
			utils.StringNotNull(name),
			0,
			sku,
			sku,
			sku,
			level,
			brand,
			attributes,
			"",
			image,
			url,
			domain.ReviewInfo{},
			1,
			"Mobile",
			seller,
		)

		products = append(products, obj)
		_, err = craw.ProductRepository.Insert(obj)

		if err != nil {
			log.Fatalf("Error to insert")
		}

	})

	return products, nil
}

func (craw *ProductCrawlerInit) getShipment(delivery string) ([]domain.Shipment, error) {
	// var delivery string = ""
	var delivery_time int = 0
	var price string
	var week_number int = 0
	var day string = ""
	var month_number int = 0
	var shipments []domain.Shipment
	var hasweek bool = false

	weekday := time.Now().Weekday()
	month := []string{"jan", "fev", "mar", "abr", "mai", "jun", "jul", "ago", "set", "out", "nov", "dez"}
	week := []string{"segunda-feira", "terça-feira", "quarta-feira", "quinta-feira", "sexta-feira", "sabado", "domingo"}
	//fmt.Println(delivery)

	re := regexp.MustCompile(`R\$(.+?) `)
	if len(re.FindStringSubmatch(delivery)) > 0 {
		price = re.FindStringSubmatch(delivery)[1]
		price = strings.Replace(price, ",", ".", -1)
		if strings.Contains(price, "Enviar") {
			price = strings.Replace(price, "Enviar", "", -1)
		}
	}

	re1 := regexp.MustCompile(`R\$(.*)`)
	if len(re.FindStringSubmatch(delivery)) == 0 && len(re1.FindStringSubmatch(delivery)) > 0 {
		price = re1.FindStringSubmatch(delivery)[1]
		price = strings.Replace(price, ",", ".", -1)
		if strings.Contains(price, "Enviar") {
			price = strings.Replace(price, "Enviar", "", -1)
		}
	}

	for i, mon := range month {
		if strings.Contains(delivery, mon) {
			//fmt.Println(mon)
			if strings.Contains(delivery, " de ") {
				re := regexp.MustCompile(`a (.+?) de ` + mon)
				if len(re.FindStringSubmatch(delivery)) > 0 {
					day = re.FindStringSubmatch(delivery)[1]
				}
			} else {
				re := regexp.MustCompile(` e (.+?) ` + mon)

				if len(re.FindStringSubmatch(delivery)) > 0 {
					day = re.FindStringSubmatch(delivery)[1]
				}

			}
			month_number = i + 1
		}
	}

	if day == "" {

		for k, we := range week {
			if strings.Contains(delivery, we) {
				week_number = k
				hasweek = true
			}
		}
		if hasweek {
			delivery_time = (7 - (int(weekday)) + int(week_number))
			if delivery_time < 0 {
				delivery_time = 0
			}
		}

	} else {
		day, _ := strconv.Atoi(day)
		t1 := time.Now()
		t2 := time.Date(2021, time.Month(month_number), day, 0, 0, 0, 0, time.UTC)
		days := t2.Sub(t1).Hours() / 24
		delivery_time = int(days)
	}

	cep, _ := strconv.Atoi(craw.Cep)

	if price == "" {
		price = "0.00"
	}

	price_float, _ := strconv.ParseFloat(price, 64)

	obj := domain.Shipment{
		cep,
		delivery_time,
		"MERCADOLIVRE",
		price_float,
		1,
		delivery,
	}
	shipments = append(shipments, obj)
	return shipments, nil

}

func getLevels(content *goquery.Document) []domain.Level {
	var levels []domain.Level

	content.Find(".andes-breadcrumb__item").Each(func(i int, leve *goquery.Selection) {
		lv := leve.Find(".andes-breadcrumb__link").Text()
		catList := domain.Level{
			lv,
			0,
		}

		levels = append(levels, catList)
	})
	return levels
}

func getPriceTo(content *goquery.Document) float64 {

	var price_to string = "0"
	var response float64

	price_to = content.Find(".ui-pdp-price.mt-16.ui-pdp-price--size-large").Find(".ui-pdp-price__second-line").Find(".price-tag-fraction").First().Text()
	if s, err := strconv.ParseFloat(price_to, 64); err == nil {
		response = s
	}
	return response
}

func getPriceToSeller(content *goquery.Selection) float64 {

	var price_to string = "0"
	var response float64

	price_to = content.Find(".ui-pdp-price").Find(".ui-pdp-price__second-line").Find(".price-tag-fraction").First().Text()
	if s, err := strconv.ParseFloat(price_to, 64); err == nil {
		response = s
	}
	return response
}

func getPriceFromSeller(content *goquery.Selection) float64 {
	var response float64
	var price_from string = "0"

	price_from = content.Find("s").Find(".price-tag-fraction").Text()

	if s, err := strconv.ParseFloat(price_from, 64); err == nil {
		response = s
	}
	return response
}

func getPriceFrom(content *goquery.Document) float64 {
	var response float64
	var price_from string = "0"

	price_from = content.Find(".ui-pdp-container__row.ui-pdp-container__row--price").Find(".price-tag.ui-pdp-price__part.ui-pdp-price__original-value.price-tag__disabled").Find(".price-tag-fraction").Text()

	if price_from == "" {
		price_from = content.Find(".ui-pdp-price.mt-16.ui-pdp-price--size-large").Find(".price-tag.ui-pdp-price__part.ui-pdp-price__original-value.price-tag__disabled").Find(".price-tag-fraction").Text()
	}

	if s, err := strconv.ParseFloat(price_from, 64); err == nil {
		response = s
	}
	return response
}

func getDescountPercent(content *goquery.Document) float64 {
	var response float64
	var descount string = "0"
	descount_percent := content.Find(".ui-pdp-price__second-line__label.ui-pdp-color--GREEN.ui-pdp-size--MEDIUM").Text()
	re := regexp.MustCompile(`(.+?)%`)
	match := re.FindStringSubmatch(descount_percent)

	if len(match) > 0 {
		descount = match[1]
	}

	response, _ = strconv.ParseFloat(descount, 64)

	return response
}

func getInstallmentValue(content *goquery.Document) float64 {
	var response float64 = getPriceTo(content)
	var installment_value string
	content_price_installment_value := content.Find(".ui-pdp-price__subtitles").Find(".ui-pdp-color--GREEN").Text()
	re := regexp.MustCompile(`R\$(.+?) s`)
	match := re.FindStringSubmatch(content_price_installment_value)

	if len(match) > 0 {
		installment_value = match[1]
		if strings.Contains(installment_value, ",") {
			installment_value = strings.ReplaceAll(installment_value, ",", ".")
		}
	}
	if s, err := strconv.ParseFloat(installment_value, 64); err == nil {
		response = s
	}

	return response
}

func getInstallmentValueSeller(content *goquery.Selection, contenTo *goquery.Document) float64 {
	var response float64 = getPriceTo(contenTo)
	var installment_value string
	content_price_installment_value := content.Find(".ui-pdp-media.ui-pdp-payment.ui-pdp-payment--md.ui-pdp-s-payment").Text()
	re := regexp.MustCompile(`R\$(.+?) s`)
	match := re.FindStringSubmatch(content_price_installment_value)

	if len(match) > 0 {
		installment_value = match[1]
		if strings.Contains(installment_value, ",") {
			installment_value = strings.ReplaceAll(installment_value, ",", ".")
		}
	}
	if s, err := strconv.ParseFloat(installment_value, 64); err == nil {
		response = s
	}

	return response
}

func getInstallment(content *goquery.Document) int {
	var response int
	var installment string = "1"
	content_installment := content.Find(".ui-pdp-price__subtitles").Find(".ui-pdp-color--GREEN").Text()
	re := regexp.MustCompile(`em (.+?)x`)
	match := re.FindStringSubmatch(content_installment)

	if len(match) > 0 {
		installment = match[1]
	}

	response, _ = strconv.Atoi(installment)

	return response
}

func getInstallmentSeller(content *goquery.Selection) int {
	var response int
	var installment string = "1"
	content_installment := content.Find(".ui-pdp-media.ui-pdp-payment.ui-pdp-payment--md.ui-pdp-s-payment").Text()
	re := regexp.MustCompile(`(.+?)x`)
	match := re.FindStringSubmatch(content_installment)

	if len(match) > 0 {
		installment = match[1]
	}

	response, _ = strconv.Atoi(installment)

	return response
}

func getInterestRate(content *goquery.Document) float64 {
	var interestRate string = "0.00"
	content_interestRate := content.Find(".ui-pdp-reviews__rating__summary__average").Text()

	if content_interestRate != "" {
		interestRate = content_interestRate
	}

	interestRate_float, _ := strconv.ParseFloat(interestRate, 64)

	return interestRate_float
}

func getImage(content *goquery.Document) []domain.Image {
	var image string = ""
	var images []domain.Image
	content_image := content.Find(".ui-pdp-gallery__figure__image")

	src_content_image, ok := content_image.Attr("src")
	if ok {
		image = src_content_image
	}
	imageList := domain.Image{
		image,
		"large",
	}

	images = append(images, imageList)

	return images
}

func getSku(content *goquery.Document) string {

	var sku string

	html, _ := content.Html()
	re := regexp.MustCompile(`sku":"(.+?)",`)
	match := re.FindStringSubmatch(html)

	if len(match) > 0 {
		sku = match[1]
	}
	return sku
}

func (craw *ProductCrawlerInit) getSeller(content *goquery.Document) []domain.Seller {
	var sellers []domain.Seller
	var installmentData []domain.Installment
	var seller string = ""
	content_seller := content.Find(".ui-seller-info").Find(".ui-pdp-seller__header__title").First().Text()

	if content_seller != "" {
		seller = content_seller
	}
	installment := getInstallment(content)
	installmentValue := getInstallmentValue(content)
	interest_rate := 0.00
	price_to := getPriceTo(content)
	price_from := getPriceFrom(content)
	descount_percent := getDescountPercent(content)

	// fmt.Println(price_to)
	// fmt.Println(price_from)
	price_term := 0.00
	obj := domain.Installment{
		"CARTAO_VISA",
		installment,
		utils.FloatNotNull(installmentValue),
		utils.FloatNotNull(price_to),
		utils.IntNotNull(interest_rate),
	}
	installmentData = append(installmentData, obj)
	delivery := content.Find(".ui-pdp-media__body").First().Text()
	shipment, _ := craw.getShipment(delivery)
	sellerr := domain.Seller{
		seller,
		domain.PricesInfo{
			domain.RetailPriceInfo{
				installmentData,
				utils.FloatNotNull(price_from),
				utils.FloatNotNull(price_to),
				utils.FloatNotNull(price_term),
				utils.FloatNotNull(0.00),
				utils.FloatNotNull(descount_percent),
			},
			domain.WholesalePriceInfo{},
		},
		"0",
		[]domain.Review{},
		shipment,
		true,
	}
	sellers = append(sellers, sellerr)
	otherSellers := craw.getOtherSellers(content)

	if len(otherSellers) > 0 {
		sellers = otherSellers
	}

	return sellers
}

func (craw *ProductCrawlerInit) getOtherSellers(content *goquery.Document) []domain.Seller {
	var sellers []domain.Seller
	var installmentData []domain.Installment
	var seller string = ""
	var page int = 0

	urlotherSellers, _ := content.Find(".ui-pdp-other-sellers__link").Attr("href")

	if urlotherSellers == "" {
		return nil
	}

	for page < 5 {

		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		fmt.Println(urlotherSellers)

		urlget := urlotherSellers + "&page=" + strconv.Itoa(page)
		request, err := http.NewRequest("GET", urlget, nil)
		if err != nil {
			fmt.Println(err)
		}
		request.Header.Set("cookie", `cp=`+craw.Cep+`%7C1622117690815`)
		resp, _ := client.Do(request)

		sellerslist, _ := goquery.NewDocumentFromReader(resp.Body)

		sellerslist.Find(".ui-pdp-buybox.ui-pdp-table__row.ui-pdp-s-table__row").Each(func(i int, s *goquery.Selection) {

			content_seller := s.Find(".ui-pdp-color--BLUE.ui-pdp-family--REGULAR").Text()

			if content_seller != "" {
				seller = content_seller
			}

			installment := getInstallmentSeller(s)
			installmentValue := getInstallmentValueSeller(s, content)
			interest_rate := 0.00
			price_to := getPriceToSeller(s)
			price_from := getPriceFromSeller(s)

			price_term := 0.00
			installmentData = []domain.Installment{}
			delivery := s.Find(".andes-tooltip__trigger").Text()
			shipment, _ := craw.getShipment(delivery)
			obj := domain.Installment{
				"CARTAO_VISA",
				installment,
				utils.FloatNotNull(installmentValue),
				utils.FloatNotNull(price_to),
				utils.IntNotNull(interest_rate),
			}
			descount_percent := getDescountPercent(content)
			installmentData = append(installmentData, obj)

			sellerr := domain.Seller{
				seller,
				domain.PricesInfo{
					domain.RetailPriceInfo{
						installmentData,
						utils.FloatNotNull(price_from),
						utils.FloatNotNull(price_to),
						utils.FloatNotNull(price_term),
						utils.FloatNotNull(0.00),
						utils.FloatNotNull(descount_percent),
					},
					domain.WholesalePriceInfo{},
				},
				"0",
				[]domain.Review{},
				shipment,
				true,
			}
			sellers = append(sellers, sellerr)
		})
		page = page + 1
	}
	return sellers
}

func getName(content *goquery.Document) string {
	var name string = ""
	content_name := content.Find(".ui-pdp-title").Text()

	if content_name != "" {
		name = content_name
	}

	return name
}

func getBrand(content *goquery.Document) string {

	var brand string = ""

	content.Find("tr").Each(func(i int, table *goquery.Selection) {

		tr := table.Find("th").Text()
		td := table.Find("td").Text()

		if strings.Contains(tr, "Marca") {
			brand = td
		}
	})

	return brand
}

func getAttributes(content *goquery.Document) []domain.Attributes {
	var attributes []domain.Attributes

	content.Find("tr").Each(func(i int, table *goquery.Selection) {

		th := table.Find("th").Text()
		td := table.Find("td").Text()

		attrList := domain.Attributes{
			utils.StringNotNull(th),
			utils.StringNotNull(td),
		}

		attributes = append(attributes, attrList)
	})

	return attributes
}
