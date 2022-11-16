package domain

import (
	"time"
)

type SubCategory struct {
	Base
	CategoryID    int    `json:"category_id"`
	DepartamentID int    `json:"departament_id"`
	Name          string `json:"name" gorm:"type:varchar(255)"`
	Url           string `json:"url" gorm:"type:varchar(255)"`
	Scan          string `json:"scan" gorm:"type:varchar(255)"`
}

func NewSubCategory(
	CategoryID int,
	DepartamentID int,
	Name string,
	Url string,
	Scan string,
) (*SubCategory, error) {

	sub_category := &SubCategory{
		CategoryID:    CategoryID,
		DepartamentID: DepartamentID,
		Name:          Name,
		Url:           Url,
		Scan:          Scan,
	}

	//category.ID = uuid.NewV4().String()
	sub_category.CreatedAt = time.Now()

	return sub_category, nil
}
