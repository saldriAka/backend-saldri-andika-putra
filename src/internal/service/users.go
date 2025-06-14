package service

import (
	"context"
	"errors"
	"saldri/backend-saldri-andika-putra/domain"
	"saldri/backend-saldri-andika-putra/dto"
	"saldri/backend-saldri-andika-putra/internal/config"
	"saldri/backend-saldri-andika-putra/internal/util"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type usersService struct {
	conf      *config.Config
	usersRepo domain.UsersRepository
}

func NewUsersService(usersRepo domain.UsersRepository, cnf *config.Config) domain.UsersService {
	return &usersService{
		usersRepo: usersRepo,
		conf:      cnf,
	}
}

func (s usersService) Register(ctx context.Context, req dto.RegisterUserRequest) error {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := domain.Users{
		ID:        uuid.NewString(),
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      req.Role,
		CreatedAt: time.Now(),
	}

	return s.usersRepo.Save(ctx, &user)
}

func (a usersService) Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error) {
	user, err := a.usersRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return dto.AuthResponse{}, errors.New("invalid credentials")
	}

	tokenStr, err := util.GenerateJWT(user.ID, user.Role, a.conf.Jwt.Key, a.conf.Jwt.Exp)
	if err != nil {
		return dto.AuthResponse{}, errors.New("authentication failed")
	}

	return dto.AuthResponse{
		Token: tokenStr,
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s usersService) GetProfile(ctx context.Context, id string) (dto.UserData, error) {
	user, err := s.usersRepo.FindById(ctx, id)
	if err != nil {
		return dto.UserData{}, err
	}

	return dto.UserData{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s usersService) List(ctx context.Context, page, limit int) ([]dto.UserData, int64, error) {
	offset := (page - 1) * limit
	users, total, err := s.usersRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var result []dto.UserData
	for _, u := range users {
		result = append(result, dto.UserData{
			ID:    u.ID,
			Email: u.Email,
			Role:  u.Role,
		})
	}

	return result, total, nil
}
