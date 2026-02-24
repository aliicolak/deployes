package user

import (
	"errors"
	"time"

	"github.com/google/uuid"

	domain "deployes/internal/domain/user"
	"deployes/pkg/utils"
)

type Service struct {
	repo      domain.Repository
	jwtSecret string
}

func NewService(repo domain.Repository, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret}
}

func (s *Service) Register(req RegisterRequest) (*UserResponse, error) {

	// 1) Email var mı kontrol et
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	// 2) Password hashle
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 3) User entity oluştur
	newUser := &domain.User{
		ID:        uuid.NewString(),
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 4) DB’ye kaydet
	err = s.repo.Create(newUser)
	if err != nil {
		return nil, err
	}

	// 5) Response dön
	return &UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}, nil
}

func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {

	// 1) Kullanıcı email ile bulunur
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 2) Password kontrol edilir
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// 3) Token pair üretilir (access + refresh)
	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// 4) Response dönülür
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	}, nil
}

// RefreshToken validates the refresh token and generates a new token pair
func (s *Service) RefreshToken(req RefreshRequest) (*LoginResponse, error) {
	// 1) Refresh token'ı doğrula ve userId al
	userId, err := utils.ValidateRefreshToken(req.RefreshToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2) Kullanıcıyı bul
	user, err := s.repo.FindByID(userId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 3) Yeni token pair üret
	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// 4) Response dönülür
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	}, nil
}
