package domain

import (
	"time"
)

type Category struct {
	Base
	DepartamentID   int    `json:"departament_id"`
	DepartamentName string `json:"departament_name"`
	Name            string `json:"name" gorm:"type:varchar(255)"`
	Url             string `json:"url" gorm:"type:varchar(255)"`
	Scan            string `json:"scan" gorm:"type:varchar(255)"`
}

func NewCategory(
	DepartamentID int,
	DepartamentName string,
	Name string,
	Url string,
	Scan string,
) (*Category, error) {

	category := &Category{
		DepartamentID:   DepartamentID,
		DepartamentName: DepartamentName,
		Name:            Name,
		Url:             Url,
		Scan:            Scan,
	}

	//category.ID = uuid.NewV4().String()
	category.CreatedAt = time.Now()

	return category, nil
}
