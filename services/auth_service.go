package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tech-azim/be-learnova/middlewares"
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

	claims := middlewares.ClaimStruct{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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
