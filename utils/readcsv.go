package utils

import (
	//"bufio"

	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/gosimple/slug"
)

type CitysAndStates struct {
	City  string
	State string
}

type Urls struct {
	Url string
}

func Read(file string) []CitysAndStates {

	var dados []CitysAndStates

	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		dados = append(dados, CitysAndStates{
			City:  slug.Make(line[0]),
			State: slug.Make(line[1]),
		})
	}
	return dados
}

func ReadUrls(file string) []Urls {

	var dados []Urls

	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		dados = append(dados, Urls{
			Url: line[0],
		})
	}
	return dados
}
