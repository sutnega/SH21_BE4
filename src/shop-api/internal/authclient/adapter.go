package authclient

import (
	"context"

	"shop-api/internal/middleware"
)

// AsVerifier возвращает адаптер, который middleware.Auth понимает.
func (c *Client) AsVerifier() middleware.TokenVerifier {
	return verifierAdapter{client: c}
}

type verifierAdapter struct {
	client *Client
}

func (a verifierAdapter) VerifyToken(ctx context.Context, token string) (*middleware.VerifyResult, error) {
	resp, err := a.client.VerifyToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return &middleware.VerifyResult{
		Valid:  resp.GetValid(),
		UserID: resp.GetUserId(),
		Email:  resp.GetEmail(),
	}, nil
}

var _ middleware.TokenVerifier = verifierAdapter{}
