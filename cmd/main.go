package main

import (
	"log"
	"os"

	"github.com/cristian-yw/Weekly10/internal/config"
	"github.com/cristian-yw/Weekly10/internal/routers"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Printf("Error loading .env file: %v", err.Error())
	// 	return
	// }
	// log.Println(os.Getenv("DB_USER"))

	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization
	// @type token
	// @description Enter your user JWT token like: Bearer <token>
	log.Println("Check ENV:", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	db, err := config.InitDB()
	if err != nil {
		log.Println("Error connecting to database: ", err.Error())
		return
	}
	defer db.Close()
	rdb := config.InitClient()
	defer rdb.Close()
	if err := config.TestDB(db); err != nil {
		log.Println("Error pinging database: ", err.Error())
		return
	}
	log.Println("Database connection successful")

	router := routers.InitRouter(db, rdb)

	router.Run("0.0.0.0:8080")
}
