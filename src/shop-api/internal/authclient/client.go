package authclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "auth-service/api/proto/auth/v1"
)

// Client — обёртка над gRPC-клиентом auth-service.
type Client struct {
	conn *grpc.ClientConn
	api  authpb.AuthServiceClient
}

// New создаёт подключение к auth-service по адресу addr (host:port).
// Соединение ленивое — реальный коннект произойдёт при первом вызове.
func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
		api:  authpb.NewAuthServiceClient(conn),
	}, nil
}

// Close закрывает соединение.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Register создаёт пользователя в auth-service и возвращает токен.
func (c *Client) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.TokenResponse, error) {
	return c.api.Register(ctx, req)
}

// Login аутентифицирует пользователя и возвращает токен.
func (c *Client) Login(ctx context.Context, email, password string) (*authpb.TokenResponse, error) {
	return c.api.Login(ctx, &authpb.LoginRequest{
		Email:    email,
		Password: password,
	})
}

// ResetPassword запускает процедуру восстановления пароля.
func (c *Client) ResetPassword(ctx context.Context, email string) error {
	_, err := c.api.ResetPassword(ctx, &authpb.ResetPasswordRequest{Email: email})
	return err
}

// VerifyToken проверяет валидность токена. Возвращает ответ,
// в котором Valid=true/false, а также userId и email при успехе.
func (c *Client) VerifyToken(ctx context.Context, token string) (*authpb.VerifyResponse, error) {
	return c.api.VerifyToken(ctx, &authpb.VerifyTokenRequest{Token: token})
}
