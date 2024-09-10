package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JuHaNi654/cms/internal/database"
	"github.com/JuHaNi654/cms/internal/models"
	"github.com/JuHaNi654/cms/internal/routes"
	"github.com/JuHaNi654/cms/internal/vite"
	"github.com/joho/godotenv"
)

//go:embed all:templates/static
var dist embed.FS

func loadViteHandler() (*vite.Handler, error) {
	if models.Environment.IsProduction() {
		fs, err := fs.Sub(dist, "templates/static")
		if err != nil {
			return nil, err
		}

		return vite.NewHandler(vite.Config{
			FS:     fs,
			IsProd: models.Environment.IsProduction(),
		})
	} else {
		return vite.NewHandler(vite.Config{
			FS:      os.DirFS("."),
			IsProd:  models.Environment.IsProduction(),
			ViteURL: "http://localhost:5173",
		})
	}
}

/*
  Connect to the sqlite database. If file doesn't exsists then we also run init.sql scripts when connected to the database
*/

func initDatabase() (*database.SqlClient, error) {
	sqliteFile := models.Environment.WithRoot("/database/sqlite.db")
	_, err := os.Stat(sqliteFile)
	if err == nil {
		return database.NewSQLClient(sqliteFile)
	} else if errors.Is(err, os.ErrNotExist) {
		db, err := database.NewSQLClient(sqliteFile)
		if err != nil {
			return nil, fmt.Errorf("error occured while connecting to the database: %s", err)
		}

		if err := db.Migrate(models.Environment.WithRoot("/schemas/init.sql")); err != nil {
			return nil, fmt.Errorf("error occured while runnin migrate: %s", err)
		}
		return db, nil
	}

	return nil, err
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("could not load .env file\n")
	}

	if err := models.LoadEnvironment(); err != nil {
		log.Fatalf("error occurred while loading configuration: %v", err)
	}

	// Create vite handler
	viteHandler, err := vite.NewHandler(vite.Config{
		FS:      os.DirFS("."),
		IsProd:  models.Environment.IsProduction(),
		ViteURL: "http://localhost:5173",
	})
	if err != nil {
		log.Fatalf("error occurred while loading vite: %v", err)
	}

	db, err := initDatabase()
	if err != nil {
		log.Fatalf("error occurred while handling database initialization: %v", err)
	}

	// Group Services
	services := &models.Services{
		Vite: viteHandler,
		DB:   db,
	}

	srv := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", models.Environment.Port),
		Handler: routes.Routes(services),
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}

		serverStopCtx()
	}()

	// Run the server
	log.Printf("server is running on: http://localhost:%s\n", models.Environment.Port)
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}
