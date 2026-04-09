package config
import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Check for environment variables (for Docker compatibility)
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "1234")
	dbname := getEnv("DB_NAME", "postgres")
	sslMode := getEnv("DB_SSLMODE", "require")

	// Construct the connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslMode,
	)
		




	// Connect to the database with optimized settings
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
	})

	if err != nil {
		log.Fatal("❌ Erro ao conectar com o banco:", err)
	}

	// Get underlying SQL database to configure connections
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("❌ Erro ao configurar pool de conexões:", err)
	}

	// 🚀 Pool otimizado para alta performance
	sqlDB.SetMaxIdleConns(25)                 // Mais conexões idle para reuso
	sqlDB.SetMaxOpenConns(200)                // Limite de conexões simultâneas
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Renovar conexões a cada 5min
	sqlDB.SetConnMaxIdleTime(2 * time.Minute) // Fechar conexões ociosas após 2min

	fmt.Println("✅ Banco de dados conectado com sucesso.")
	DB = database
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
