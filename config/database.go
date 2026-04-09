package config

import (
	"fmt"
	"log"
	"os"

	"hookfy/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	godotenv.Load()

	dbPath := getEnv("DB_PATH", "hookfy.db")

	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		log.Fatal("Erro ao conectar com o banco:", err)
	}

	database.AutoMigrate(&models.Webhook{})

	fmt.Println("Banco de dados conectado com sucesso.")
	DB = database
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
