package connections

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB returns a gorm DB connection
func DB() *gorm.DB {
	_ = godotenv.Load("db.env")
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host= %s port=%s sslmode=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_SSL"))
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db
}
