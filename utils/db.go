package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"gitlab.com/eiprice/crawlers/mercadolivre/domain"
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	err := godotenv.Load(basepath + "/../.env")

	if err != nil {
		log.Fatalf("Error loading .env files")
	}
}

func ConnectDB(drop string) *gorm.DB {
	var dsn string
	var db *gorm.DB
	var err error

	dsn = os.Getenv("dsn")
	db, err = gorm.Open(os.Getenv("dbType"), dsn)

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		panic(err)
	}

	if os.Getenv("debug") == "true" {
		db.LogMode(true)
	}

	if drop == "all" {
		log.Printf("Cleaning Database...")
		db.Exec("TRUNCATE mercadolivre.products;")
		db.Exec("TRUNCATE mercadolivre.departaments;")
		db.Exec("TRUNCATE mercadolivre.categories;")
		db.Exec("TRUNCATE mercadolivre.sub_categories;")
		log.Printf("Database is clean")
	} else {
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return os.Getenv("schema") + "." + defaultTableName
		}

		if os.Getenv("AutoMigrateDb") == "true" {
			db.AutoMigrate(&domain.Product{}, &domain.Departament{}, &domain.Category{}, &domain.SubCategory{})
			os.Setenv("AutoMigrateDb", "false")
		}
	}

	return db
}
