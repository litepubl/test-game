package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/litepubl/test-game/pkg/controller"
	"github.com/litepubl/test-game/pkg/items"
	"github.com/litepubl/test-game/pkg/items/repo"
	"github.com/litepubl/test-game/pkg/logger"
	"github.com/litepubl/test-game/pkg/postgres"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type app struct {
	config          *Config
	logConfigurator *logger.Configurator
	pg              *postgres.Postgres
	redis           *redis.Client
	nc              *nats.Conn
	router          *gin.Engine
}

func NewApp() (*app, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}

	app := &app{
		config: config,
	}

	app.logConfigurator, err = logger.NewConfigurator(&config.Log)
	if err != nil {
		return nil, err
	}

	log.Info().Msg(fmt.Sprintf("%#v", config))

	app.pg, err = postgres.New(&config.PG)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connect to postgress")
		return nil, err
	}

	app.redis = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Error connect to redis")
		return nil, err
	}

	app.nc, err = nats.Connect("nats://" + config.Nats.Host + ":" + config.Nats.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connect to nats")
		return nil, err
	}

	app.router = gin.New()
	app.router.Use(gin.Recovery())
	app.router.Use(gin.Logger())

	ic := controller.NewItemsController(
		items.New(
			repo.New(app.pg),
			app.redis,
			app.nc,
		),
	)

	controller.InitRoutes(app.router, ic)

	return app, nil
}

func (app *app) Run() error {
	srv := &http.Server{
		Addr:    ":" + app.config.HTTP.Port,
		Handler: app.router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	failStartServer := make(chan any, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Error listen http server")
			failStartServer <- nil
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
	case <-failStartServer:
	}

	log.Info().Msg("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Info().Err(err).Msg("Server forced to shutdown")
		return err
	}

	log.Info().Msg("Server exiting")
	return nil
}

func (app *app) GetConfig() *Config {
	return app.config
}

func (app *app) Close() {
	app.pg.Close()
	app.pg = nil
	app.redis.Close()
	app.redis = nil
	app.nc.Close()
	app.nc = nil
	app.logConfigurator.Close()

}
