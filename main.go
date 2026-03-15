package main

import (
	"log"
	"net/http"
	"os"
	"rest-api/internal/database"
	"rest-api/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://tasks_user:tasks_password@localhost:5433/tasks_db?sslmode=disable"
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8081"
	}
	log.Printf("Начинаем запуск сервера  %s", serverPort)
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	defer db.Close()
	log.Println("Подключение к базе данных успешно!")
	taskstore := database.NewTaskStore(db)
	handlers := handlers.NewHandlers(taskstore)
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", methodHandler(handlers.GetAllTasks, http.MethodGet))
	mux.HandleFunc("/tasks/create", methodHandler(handlers.CreateTask, http.MethodPost))
	mux.HandleFunc("/tasks/", taskIDHandler(handlers))

	loggedMux := loggingMiddleware(mux)
	serverAddr := ":" + serverPort
	log.Printf("Сервер запущен на %s", serverAddr)
	err = http.ListenAndServe(serverAddr, loggedMux)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

}

func methodHandler(handlerfunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}
		handlerfunc(w, r)
	}
}

func taskIDHandler(handler *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTaskByID(w, r)
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
