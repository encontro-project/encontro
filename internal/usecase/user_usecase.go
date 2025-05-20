package usecase

import (
	"encontro/internal/domain/entity"
	"encontro/internal/domain/service"
	"encontro/internal/infrastructure/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo repository.UserRepository
	jwtSvc   *service.JWTService
}

func NewUserUseCase(repo repository.UserRepository, jwtSvc *service.JWTService) *UserUseCase {
	return &UserUseCase{userRepo: repo, jwtSvc: jwtSvc}
}

func (uc *UserUseCase) Register(username, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &entity.User{
		Username: username,
		Email:    email,
		Password: string(hash),
	}
	return uc.userRepo.Create(user)
}

func (uc *UserUseCase) Login(email, password string) (string, error) {
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	return uc.jwtSvc.GenerateToken(user.ID)
}
