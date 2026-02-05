package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/utils"
)

type AuthService interface {
	Login(email string, password string) (string, models.User, error)
	Register(user models.User) (models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
}


func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo,
	}
}

// Login implements [AuthService].
func (a *authService) Login(email string, password string) (string, models.User, error) {
	user, err := a.userRepo.FindByEmail(email)

	if err != nil {
		return "", models.User{}, errors.New("Account not found")
	}

	if _, err := utils.Descrypt(password, user.Password); err != nil {
		return "", models.User{}, errors.New("Wrong password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", models.User{}, errors.New("failed to generate token")
	}

	user.Password = ""
	return tokenString, user, nil
}

// Register implements [AuthService].
func (a *authService) Register(user models.User) (models.User, error) {
	panic("unimplemented")
}

