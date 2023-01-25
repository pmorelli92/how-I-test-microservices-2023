package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	env "github.com/Netflix/go-env"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pmorelli92/how-i-test-microservices-2023/database"
	"github.com/pmorelli92/how-i-test-microservices-2023/features/create_post"
)

type config struct {
	DBHost      string `env:"DB_HOST"`
	DBUser      string `env:"DB_USER"`
	DBPass      string `env:"DB_PASS"`
	DBPort      string `env:"DB_PORT"`
	DBName      string `env:"DB_NAME"`
	HTTPAddress string `env:"HTTP_ADDRESS"`
}

func (c config) databaseConn(driver string) string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		driver,
		c.DBUser,
		c.DBPass,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var config config
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Fatal(err)
	}

	if err = database.Migrate(config.databaseConn("pgx")); err != nil {
		log.Fatal(err)
	}

	db, err := pgxpool.New(context.Background(), config.databaseConn("postgresql"))
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
