package main

import (
	"log"
	"net/http"
	"os"
	"sdn-xml-api/config"
	mypg "sdn-xml-api/internal/database/postgres"
)

func main() {
	// инициализация конфигурации
	cfg := config.LoadConfig()

	// инициализация подключения к БД
	dbConnStr := cfg.DBConnectionString()

	db, err := mypg.NewDB(dbConnStr)
	if err != nil {
		log.Fatalf("failed to initialize database connection: %v", err)
	}
	defer db.Close()

	// инициализация маршрутизатора
	router := initRouter(db.DB)

	// запуск веб-сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting web server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("failed to start web server: %v", err)
	}
}
