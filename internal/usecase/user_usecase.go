package usecase

import (
	"context"

	"github.com/codepnw/go-cart-system/config"
	"github.com/codepnw/go-cart-system/internal/domain"
	"github.com/codepnw/go-cart-system/internal/dto"
	"github.com/codepnw/go-cart-system/internal/repository"
	"github.com/codepnw/go-cart-system/internal/utils/errs"
	"github.com/codepnw/go-cart-system/internal/utils/security"
)

type UserUsecase interface {
	Register(ctx context.Context, req *dto.CreateUser) error
	Login(ctx context.Context, req *dto.UserCredential) (*dto.LoginResponse, error)
}

type userUsecase struct {
	repo   repository.UserRepository
	config config.EnvConfig
}

func NewUserUsecase(repo repository.UserRepository, config config.EnvConfig) UserUsecase {
	return &userUsecase{
		repo:   repo,
		config: config,
	}
}

func (u *userUsecase) Register(ctx context.Context, req *dto.CreateUser) error {
	// check email
	exists, _ := u.repo.GetByEmail(ctx, req.Email)
	if exists != nil {
		return errs.ErrUserAlreadyExists
	}

	// hash password
	hashed, err := security.GenerateHashPassword(req.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Email:    req.Email,
		Password: hashed,
		Role:     "user",
	}

	return u.repo.Create(ctx, &user)
}

func (u *userUsecase) Login(ctx context.Context, req *dto.UserCredential) (*dto.LoginResponse, error) {
	// check email
	user, err := u.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	// check password
	if err := security.CheckHashPassword(req.Password, user.Password); err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	// config jwt token
	jwt := security.NewTokenConfig(u.config, user)

	accessToken, err := jwt.GenerateAccessToken()
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
