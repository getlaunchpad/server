package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *chi.Mux
}

// Initializes server for user
// this includes composing routes, middleware and db
func (s *Server) Initialize() {
	/* Setup postgres */
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_NAME"))
	s.DB, err = gorm.Open(os.Getenv("DB_DRIVER"), DBURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Turn this off during production
	s.DB.LogMode(true)

	defer s.DB.Close()

	// server.DB.AutoMigrate(&models.User{})

	// Creates new chi mux and setup middlware
	s.Router = chi.NewRouter()
	setupMiddleware(s.Router)

	// Creates application routes
	s.initializeRoutes()

	// Define http server
	h := &http.Server{
		Handler:      s.Router,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Concurrently start http server process
	go func() {
		log.Println("Starting Server")
		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown(h)
}

// Gracefully shut down server
func waitForShutdown(s *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we recieve signal
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Panic("Error shutting down")
	}

	log.Println("Shutting down")
	os.Exit(0)
}

// Sets up various middlewares
func setupMiddleware(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Compress(6, "application/json"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))
}
