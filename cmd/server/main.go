package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"pr-reviewer/internal/handlers"
	"pr-reviewer/internal/repository"
	"pr-reviewer/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Получаем параметры подключения из окружения
	dbHost := getEnv("POSTGRES_HOST", "localhost")
	dbPort := getEnv("POSTGRES_PORT", "5432")
	dbUser := getEnv("POSTGRES_USER", "postgres")
	dbPassword := getEnv("POSTGRES_PASSWORD", "postgres")
	dbName := getEnv("POSTGRES_DB", "pr_reviewer")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	// Создаём репозитории
	teamRepo := repository.NewTeamRepo(db)
	userRepo := repository.NewUserRepo(db)
	prRepo := repository.NewPRRepo(db)

	// Создаём сервисы
	teamService := service.NewTeamService(teamRepo, userRepo)
	userService := service.NewUserService(userRepo, teamRepo)
	prService := service.NewPRService(prRepo, userRepo, teamRepo)

	// Создаём обработчики
	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService, prService)
	prHandler := handlers.NewPRHandler(prService)

	// Создаём маршрутизатор
	r := mux.NewRouter()
	teamHandler.RegisterTeamRoutes(r)
	userHandler.RegisterUserRoutes(r)
	prHandler.RegisterPRRoutes(r)

	// Эндпоинт здоровья
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Запуск сервера
	serverAddr := ":8080"
	srv := &http.Server{
		Handler:      r,
		Addr:         serverAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is running at %s", serverAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// getEnv возвращает значение переменной окружения или дефолт
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
