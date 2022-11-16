# Crawler in Go to get Data from State City, departaments (Mercado Livre)

 
 You can pass parameters from command line to define what you want to get from mercadolivre.

# Steps Install

    1. git clone https://gitlab.com/eiprice/crawlers/mercadolivre
    2. cd /mercadolivre
    3. https://golang.org/doc/install ( Install Go By your OS )
    4. open directory "mercadolivre" and Write:
    ` go install `

    5. create .env file
    `cp env.example .env`

    6. create folder files
    `mkdir files`
    `sudo chmod -R 777 files`

    7. Up PostgreSql
    ` docker-compose up -d `

    8. Create mercadolivre schema on postgresql
    `CREATE SCHEMA mercadolivre;`

    
# Commands


## Get All ( State Default - Sao Paulo )

` go run main.go `

## Get from Other State

` go run main.go -state "Sao Paulo" `

## Get from Other State and City

` go run main.go -state "Sao Paulo" -city "Assis" `

## Get from a specific Departament

` go run main.go -departament "Acessórios para Veículos"`

## Get from a specific Departament and Category

` go run main.go -departament "Acessórios para Veículos" -category "Aces. de Carros e Caminhonetes"`

## Get from a specific Departament, Category and Sub Category

` go run main.go -departament "Acessórios para Veículos" -category "Aces. de Carros e Caminhonetes" -subCategory "Exterior"`

## Set a name or number for this scan

` go run main.go -departament "Acessórios para Veículos" -category "Aces. de Carros e Caminhonetes" -subCategory "Exterior" -state "Sao Paulo" -city "Assis" -scan "1"`

## get multiples States and City from file

` go run main.go -file="leitura.csv" `

## Clean Database 

` make clean ` 

## Upload S3

`bucket : eiprice.delivery /mercadolivre`




