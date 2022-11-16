package domain

import (
	"time"
)

type Departament struct {
	Base
	Name string `json:"name" gorm:"type:varchar(255)"`
	Url  string `json:"url" gorm:"type:varchar(255)"`
	Scan string `json:"scan" gorm:"type:varchar(255)"`
}

func NewDepartament(
	Name string,
	Url string,
	Scan string,
) (*Departament, error) {

	departament := &Departament{
		Name: Name,
		Url:  Url,
		Scan: Scan,
	}

	//category.ID = uuid.NewV4().String()
	departament.CreatedAt = time.Now()

	return departament, nil
}
