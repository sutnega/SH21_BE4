package di

import (
	"database/sql"
	"net/http"

	"github.com/google/wire"
	"github.com/gorilla/mux"

	"shop-api/internal/authclient"
	"shop-api/internal/config"
	"shop-api/internal/db"
	"shop-api/internal/handlers"
	"shop-api/internal/middleware"
	"shop-api/internal/repository"
	"shop-api/internal/router"
)

var ProviderSet = wire.NewSet(
	config.Load,
	provideDB,
	provideAuthClient,

	repository.NewClientRepository,
	repository.NewProductRepository,
	repository.NewSupplierRepository,
	repository.NewImageRepository,

	handlers.NewClientHandler,
	handlers.NewProductHandler,
	handlers.NewSupplierHandler,
	handlers.NewImageHandler,
	handlers.NewAuthHandler,

	provideTokenVerifier,
	provideAuthMiddleware,
	router.New,
	wire.Struct(new(router.Handlers), "*"),

	provideHTTPServer,
	NewApp,
)

func provideDB(cfg *config.Config) (*sql.DB, error) {
	return db.Connect(cfg.DSN)
}

func provideAuthClient(cfg *config.Config) (*authclient.Client, error) {
	return authclient.New(cfg.AuthGRPCAddr)
}

func provideTokenVerifier(c *authclient.Client) middleware.TokenVerifier {
	return c.AsVerifier()
}

func provideAuthMiddleware(v middleware.TokenVerifier) mux.MiddlewareFunc {
	return middleware.Auth(v)
}

func provideHTTPServer(cfg *config.Config, h http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: h,
	}
}
