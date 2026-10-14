package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authpb "auth-service/api/proto/auth/v1"
	"auth-service/internal/hasher"
	"auth-service/internal/mailer"
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/internal/token"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer

	repo   repository.UserRepository
	hasher hasher.Hasher
	token  token.Manager
	mailer mailer.Mailer
}

func NewAuthService(
	repo repository.UserRepository,
	h hasher.Hasher,
	t token.Manager,
	m mailer.Mailer,
) *AuthService {
	return &AuthService{repo: repo, hasher: h, token: t, mailer: m}
}

func (s *AuthService) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.TokenResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	hash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "hash password failed")
	}

	u := &models.User{
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		Name:         req.Name,
		Surname:      req.Surname,
		Phone:        req.Phone,
		PasswordHash: hash,
	}

	created, err := s.repo.Create(ctx, u)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}
		return nil, status.Error(codes.Internal, "create user failed")
	}

	tokenStr, err := s.token.Issue(created.ID.String(), created.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "issue token failed")
	}

	return &authpb.TokenResponse{Token: tokenStr, UserId: created.ID.String()}, nil
}

func (s *AuthService) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.TokenResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	u, err := s.repo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "get user failed")
	}

	if !s.hasher.Compare(u.PasswordHash, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	tokenStr, err := s.token.Issue(u.ID.String(), u.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "issue token failed")
	}

	return &authpb.TokenResponse{Token: tokenStr, UserId: u.ID.String()}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, req *authpb.ChangePasswordRequest) (*authpb.Empty, error) {
	if req.UserId == "" || req.OldPassword == "" || req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, old_password, new_password required")
	}

	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "get user failed")
	}

	if !s.hasher.Compare(u.PasswordHash, req.OldPassword) {
		return nil, status.Error(codes.Unauthenticated, "invalid old password")
	}

	hash, err := s.hasher.Hash(req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, "hash password failed")
	}

	if err := s.repo.UpdatePassword(ctx, id, hash); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "update password failed")
	}

	return &authpb.Empty{}, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req *authpb.ResetPasswordRequest) (*authpb.Empty, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email required")
	}

	u, err := s.repo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "get user failed")
	}

	newPwd, err := generatePassword(12)
	if err != nil {
		return nil, status.Error(codes.Internal, "generate password failed")
	}

	hash, err := s.hasher.Hash(newPwd)
	if err != nil {
		return nil, status.Error(codes.Internal, "hash password failed")
	}

	if err := s.repo.UpdatePassword(ctx, u.ID, hash); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "update password failed")
	}

	// Заглушка: по заданию достаточно вывести в консоль
	s.mailer.SendPassword(u.Email, newPwd)

	return &authpb.Empty{}, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyResponse, error) {
	if req.Token == "" {
		return &authpb.VerifyResponse{Valid: false}, nil
	}

	userID, email, err := s.token.Parse(req.Token)
	if err != nil {
		return &authpb.VerifyResponse{Valid: false}, nil
	}

	return &authpb.VerifyResponse{
		Valid:  true,
		UserId: userID,
		Email:  email,
	}, nil
}

func validateRegister(req *authpb.RegisterRequest) error {
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return status.Error(codes.InvalidArgument, "invalid email")
	}
	if req.Name == "" || req.Surname == "" {
		return status.Error(codes.InvalidArgument, "name and surname required")
	}
	if req.Phone == "" {
		return status.Error(codes.InvalidArgument, "phone required")
	}
	if len(req.Password) < 6 {
		return status.Error(codes.InvalidArgument, "password must be at least 6 characters")
	}
	return nil
}

func generatePassword(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := base64.RawURLEncoding.EncodeToString(b)
	if len(s) > n {
		s = s[:n]
	}
	return s, nil
}
