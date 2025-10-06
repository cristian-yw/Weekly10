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

// package main

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"path/filepath"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// const (
// 	tmdbBaseURL = "https://image.tmdb.org/t/p/original" // bisa diganti w500 kalau mau lebih kecil
// 	posterDir   = "./uploads/poster"
// 	backdropDir = "./uploads/backdrop"
// )

// type Movie struct {
// 	ID           int
// 	PosterPath   *string
// 	BackdropPath *string
// }

// func main() {
// 	ctx := context.Background()

// 	db, err := pgxpool.New(ctx, "postgres://cegans:12345@localhost:5423/Learn_db?sslmode=disable")

// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	rows, err := db.Query(ctx, `
// 		SELECT id, poster_path, backdrop_path
// 		FROM movies
// 		WHERE poster_path IS NOT NULL OR backdrop_path IS NOT NULL
// 	`)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer rows.Close()

// 	// Buat folder kalau belum ada
// 	os.MkdirAll(posterDir, os.ModePerm)
// 	os.MkdirAll(backdropDir, os.ModePerm)

// 	for rows.Next() {
// 		var m Movie
// 		if err := rows.Scan(&m.ID, &m.PosterPath, &m.BackdropPath); err != nil {
// 			panic(err)
// 		}

// 		if m.PosterPath != nil {
// 			originalName := filepath.Base(*m.PosterPath) // ambil nama asli dari TMDB path
// 			outPath := filepath.Join(posterDir, originalName)
// 			err := downloadFile(*m.PosterPath, outPath)
// 			if err != nil {
// 				fmt.Println("Poster gagal:", err)
// 			} else {
// 				fmt.Println("Poster tersimpan:", originalName)
// 			}
// 		}
// 		if m.BackdropPath != nil {
// 			originalName := filepath.Base(*m.BackdropPath)
// 			outPath := filepath.Join(backdropDir, originalName)
// 			err := downloadFile(*m.BackdropPath, outPath)
// 			if err != nil {
// 				fmt.Println("Backdrop gagal:", err)
// 			} else {
// 				fmt.Println("Backdrop tersimpan:", originalName)
// 			}
// 		}
// 	}
// }

// func downloadFile(path, outPath string) error {
// 	url := tmdbBaseURL + path
// 	fmt.Println("Downloading:", url)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("bad status: %s", resp.Status)
// 	}

// 	out, err := os.Create(outPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	_, err = io.Copy(out, resp.Body)
// 	return err
// }
