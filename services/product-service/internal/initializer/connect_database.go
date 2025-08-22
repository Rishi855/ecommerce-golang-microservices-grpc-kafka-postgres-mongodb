package initializer

import (
	"fmt"
	"log"
	"os"
)

func ConnectDatabase() *gorm.DB {
	var dns string
	if dbConfig := os.Getenv("DATABSAE_CONFIG"); dbConfig != "" {
		dns = dbConfig
	} else {
		dns = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSLMODE"),
		)
	}
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil{
		log.Fatalln(err)
	}

	db.AutoMigrate(&model.Product{}, &model.StockDecreaseLog{})
	return db
}
