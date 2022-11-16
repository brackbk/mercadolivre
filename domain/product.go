package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Level struct {
	Name       string `json:"name"`
	Menu_level int    `json:"menu_level"`
}

type Image struct {
	Url  string `json:"url"`
	Size string `json:"size"`
}

type Review struct {
	Url  string `json:"url"`
	Size string `json:"size"`
}

type ReviewInfo struct {
	Total   int      `json:"total"`
	Average float64  `json:"average"`
	Reviews []Review `json:"reviews"`
}

type Installment struct {
	Method        string  `json:"method"`
	Quantity      int     `json:"quantity"`
	Value         float64 `json:"value"`
	Term_value    float64 `json:"term_value"`
	Interest_rate int     `json:"interest_rate"`
}

type RetailPriceInfo struct {
	Installments     []Installment `json:"installments"`
	From_value       float64       `json:"from_value"`
	To_value         float64       `json:"to_value"`
	Billet_value     float64       `json:"billet_value"`
	Term_value       float64       `json:"term_value"`
	Discount_percent float64       `json:"discount_percent"`
}

type WholesalePriceInfo struct {
	Minimum_quantity int `json:"minimum_quantity"`
}

type Shipment struct {
	Postal_code   int     `json:"postal_code"`
	Delivery_time int     `json:"delivery_time"`
	Delivery_type string  `json:"delivery_type"`
	Freight_price float64 `json:"freight_price"`
	Stock         int     `json:"stock"`
	Obs           string  `json:"obs"`
}

type PricesInfo struct {
	Retail    RetailPriceInfo    `json:"retail"`
	Wholesale WholesalePriceInfo `json:"wholesale"`
}

type Seller struct {
	Name      string     `json:"name"`
	Prices    PricesInfo `json:"prices"`
	Seller_id string     `json:"seller_id"`
	Reviews   []Review   `json:"reviews"`
	Shipment  []Shipment `json:"shipment"`
	In_stock  bool       `json:"in_stock"`
}

type Attributes struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type Product struct {
	ID                bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name              string        `json:"name"`
	Ean               int           `json:"ean"`
	Sku               string        `json:"sku"`
	Product_id        string        `json:"product_id"`
	Product_reference string        `json:"product_reference"`
	Levels            []Level       `json:"levels"`
	Brand             string        `json:"brand"`
	Attributes        []Attributes  `json:"attributes"`
	Model             string        `json:"model"`
	Images            []Image       `json:"images"`
	Url               string        `json:"url"`
	Reviews           ReviewInfo    `json:"reviews"`
	Scan              int           `json:"scan"`
	Origem            string        `json:"origem"`
	Sellers           []Seller      `json:"sellers"`
	Created_at        time.Time     `json:"created_at"`
}

func NewProduct(
	Name string,
	Ean int,
	Sku string,
	Product_id string,
	Product_reference string,
	Levels []Level,
	Brand string,
	Attributes []Attributes,
	Model string,
	Images []Image,
	Url string,
	Reviews ReviewInfo,
	Scan int,
	Origem string,
	Sellers []Seller,
) (*Product, error) {

	Product := &Product{
		Name:              Name,
		Ean:               Ean,
		Sku:               Sku,
		Product_id:        Product_id,
		Product_reference: Product_reference,
		Levels:            Levels,
		Brand:             Brand,
		Attributes:        Attributes,
		Model:             Model,
		Images:            Images,
		Url:               Url,
		Reviews:           Reviews,
		Scan:              Scan,
		Origem:            Origem,
		Sellers:           Sellers,
	}

	Product.Created_at = time.Now()

	return Product, nil
}
