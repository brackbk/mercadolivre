package repositories

import (
	"log"
	"time"

	"gitlab.com/eiprice/crawlers/mercadolivre/domain"
	"gitlab.com/eiprice/crawlers/mercadolivre/utils"
)

type ProductRepository interface {
	Insert(product *domain.Product) (*domain.Product, error)
}

type ProductRepositoryDb struct {
	C          *utils.Context
	Collection string
}

func (repo ProductRepositoryDb) Drop() {
	c := repo.C.Collection(repo.Collection)
	c.DropCollection()
	// if err := c.DropCollection().Error(); err != "" {
	// 	log.Fatal(err)
	// }
}
func (repo ProductRepositoryDb) Insert(product *domain.Product) (*domain.Product, error) {

	product.Created_at = time.Now()
	c := repo.C.Collection(repo.Collection)
	err := c.Insert(&product)

	if err != nil {
		log.Fatalf("Error inserting product")
	}

	return product, nil
}
