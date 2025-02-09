package database

import (
	"fmt"
	"log"
	"music_library_api/internal/models"
	testData "music_library_api/test"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// db postgres
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "postgres"
)

// Load .env file
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("loading .env file to database successfully")

}

var db *gorm.DB

func ConnectToPostgres() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	log.Println("✅ Successfully connected to the database!")

	Migrate()
	LoadTestData()
}

func GetDB() *gorm.DB {
	return db
}

func Migrate() {
	if !db.Migrator().HasTable(&models.Song{}) {
		err := db.AutoMigrate(&models.Song{})
		if err != nil {
			log.Fatal(" Migration failed:", err)
		}
		log.Println("Database migrated successfully!")
	} else {
		log.Println("Info: Table already exists, skipping migration.")
	}
}

func LoadTestData() {
	//check
	var count int64
	err := db.Model(&models.Song{}).Count(&count).Error
	if err != nil {
		log.Fatal("Failed to check existing data:", err)
	}

	if count > 0 {
		log.Println("Info: Test data already exists, skipping loading test data...")
		return
	}

	for _, song := range testData.TestData {
		err := db.Create(&song).Error
		if err != nil {
			log.Fatal("Failed to insert test data:", err)
		}
	}
	log.Println("Test data inserted successfully!")

}
