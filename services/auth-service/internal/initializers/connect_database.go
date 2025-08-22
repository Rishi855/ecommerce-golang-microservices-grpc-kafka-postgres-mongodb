package initializers

import (
    "fmt"
    "log"
    "os"

    "auth-service/internal/model"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
    var dsn string
    
    // Check if we have a single DATABASE_CONFIG (Docker environment)
    if dbConfig := os.Getenv("DATABASE_CONFIG"); dbConfig != "" {
        dsn = dbConfig
    } else {
        // Fallback to individual environment variables (local development)
        dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
            os.Getenv("DB_HOST"),
            os.Getenv("DB_PORT"),
            os.Getenv("DB_USER"),
            os.Getenv("DB_PASSWORD"),
            os.Getenv("DB_NAME"),
            os.Getenv("DB_SSLMODE"),
        )
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalln(err)
    }

    db.AutoMigrate(&model.User{}, &model.Admin{})
    return db
}