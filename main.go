package main

import (
	"log"
	_ "music_library_api/docs"
	"music_library_api/internal/database"
	"music_library_api/internal/endpoints"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Music API
// @version 1.0
// @description Test task : This is a simple music API to manage songs
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

func main() {

	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // default port - 8080
	}

	database.ConnectToPostgres()
	http.HandleFunc("/api/getSongs", endpoints.GetSongs)
	http.HandleFunc("/api/AddSong", endpoints.AddSongs)
	http.HandleFunc("/api/deleteSong", endpoints.DeleteSong)
	http.HandleFunc("/api/EditSong", endpoints.EditSong)
	http.HandleFunc("/api/GetSongText", endpoints.GetSongText)

	http.Handle("/swagger/", httpSwagger.WrapHandler) //http://localhost:8080/swagger/index.html

	log.Printf("Server started on :%s", port)
	// Start
	http.ListenAndServe(":"+port, nil)

}
