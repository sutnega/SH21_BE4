# Auth Service

gRPC-сервис авторизации для магазина бытовой техники.

## Требования
- Go 1.27+
- PostgreSQL 14+
- protoc + protoc-gen-go + protoc-gen-go-grpc

## Установка protoc плагинов
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Создание БД
```bash
createdb auth_db
```

## Генерация proto
```bash
make proto
```

## Миграции
```bash
export DB_DSN="postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"
make migrate
```

## Запуск
```bash
make run
```

## Проверка через grpcurl
```bash
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext -d '{"email":"a@b.c","name":"Ivan","surname":"Ivanov","phone":"+79001234567","password":"secret"}' \
  localhost:50051 auth.v1.AuthService/Register
```