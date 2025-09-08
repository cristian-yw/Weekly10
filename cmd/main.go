package main

import (
	"log"

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

	// @securityDefinitions.apikey Bearer
	// @in header
	// @name Authorization
	// @type token
	// @description Enter your user JWT token like: Bearer <token>
	db, err := config.InitDB()
	if err != nil {
		log.Println("Error connecting to database: ", err.Error())
		return
	}
	defer db.Close()

	if err := config.TestDB(db); err != nil {
		log.Println("Error pinging database: ", err.Error())
		return
	}
	log.Println("Database connection successful")

	router := routers.InitRouter(db)

	router.Run("0.0.0.0:8080")
}
