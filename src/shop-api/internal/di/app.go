package di

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"shop-api/internal/authclient"
	"shop-api/internal/config"
)

// App — собранное приложение. Все зависимости доступны через Wire.
type App struct {
	cfg        *config.Config
	srv        *http.Server
	db         *sql.DB
	authClient *authclient.Client
}

func NewApp(
	cfg *config.Config,
	srv *http.Server,
	db *sql.DB,
	authClient *authclient.Client,
) *App {
	return &App{
		cfg:        cfg,
		srv:        srv,
		db:         db,
		authClient: authClient,
	}
}

// Serve запускает HTTP-сервер (блокирующий вызов).
func (a *App) Serve() error {
	log.Printf("shop-api listening on :%s (swagger: http://localhost:%s/swagger/index.html)",
		a.cfg.Port, a.cfg.Port)
	err := a.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Close аккуратно закрывает все ресурсы.
func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = a.srv.Shutdown(ctx)
	if a.authClient != nil {
		_ = a.authClient.Close()
	}
	if a.db != nil {
		_ = a.db.Close()
	}
	return nil
}
