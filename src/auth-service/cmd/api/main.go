package main

import (
	"log"
	"net"

	"google.golang.org/grpc/reflection"

	authpb "auth-service/api/proto/auth/v1"
	"auth-service/internal/config"
	"auth-service/internal/db"
	"auth-service/internal/grpcserver"
	"auth-service/internal/hasher"
	"auth-service/internal/mailer"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/internal/token"
)

func main() {
	cfg := config.Load()

	conn, err := db.Connect(cfg.DSN)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer conn.Close()

	userRepo := repository.NewUserRepository(conn)
	h := hasher.New()
	tk := token.New(cfg.JWTSecret)
	ml := mailer.New()

	authService := service.NewAuthService(userRepo, h, tk, ml)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpcserver.New()
	authpb.RegisterAuthServiceServer(srv.GRPC(), authService)
	reflection.Register(srv.GRPC())

	log.Printf("auth-service listening on :%s", cfg.GRPCPort)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
