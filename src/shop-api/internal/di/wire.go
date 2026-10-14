//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
)

// InitializeApp — точка входа для Wire. Возвращает полностью собранное приложение.
func InitializeApp() (*App, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
