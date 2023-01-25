package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pmorelli92/how-i-test-microservices-2023/database"
	"github.com/pmorelli92/how-i-test-microservices-2023/features/create_post"
)

func Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err = database.Migrate(config.migrateDatabaseDSN()); err != nil {
		log.Fatal(err)
	}

	db, err := pgxpool.New(context.Background(), config.ConnectDatabaseDSN())
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	server := &http.Server{Addr: config.HTTPAddress, Handler: router}

	create_post.Setup(router, db)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-stop

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
