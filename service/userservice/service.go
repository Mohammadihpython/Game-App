package userservice

import (
	"GameApp/entity"
	"GameApp/pkg/phonenumber"
	"GameApp/repository/mysql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	RegisterUser(u entity.User) (entity.User, error)
	GetUserByPhone(phoneNumber string) (entity.User, bool, error)
	GetUserByID(userID uint) (entity.User, error)
}

var (
	secret = "%324sdf"
)

type Service struct {
	signKey string
	repo    Repository
}
type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}
type RegisterResponse struct {
	User entity.User `json:"user"`
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO : We should verify phone number by verification code
	// Validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("invalid phone")
	}
	// check uniqueness of phone number
	// we used shorthand if here
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			// %w wrap the error and show us the last errors corrupted for this error
			//ارور های قبلی را که مربوط به این خطا هست زا نیز نشان می دهد
			return RegisterResponse{}, fmt.Errorf("failed to check if phone is unique: %w", err)
		}
		if !isUnique {
			return RegisterResponse{}, errors.New("phone number is already used")
		}
	}

	//validate name
	// TODO : Add support for Persion word or not ASID word
	if len(req.Name) < 3 {
		return RegisterResponse{}, errors.New("name is too short")
	}

	// TODO validate password with regex
	if len(req.Password) > 8 {
		return RegisterResponse{}, errors.New("password length must greater than 8")
	}

	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    mysql.GetMD5Hash(req.Password),
	}
	// create new user in storage
	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("failed to register user: %w", err)
	}

	//	 return created user
	return RegisterResponse{User: createdUser}, nil

}

func New(repo Repository, signKey string) Service {
	return Service{repo: repo, signKey: signKey}

}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	//RefreshToken string `json:"refresh_token"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	//  check the existence of phone number from repository
	//	get the user by phone number
	// TODO : Its better to  separate method for check existence check  of get user by phone number
	user, exists, err := s.repo.GetUserByPhone(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}
	if !exists {
		return LoginResponse{}, fmt.Errorf("user or password not valid")
	}
	if user.Password != mysql.GetMD5Hash(req.Password) {
		return LoginResponse{}, fmt.Errorf("user or password not valid ")

	}
	token, err := createToken(user.ID, s.signKey)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to create token: %w", err)
	}

	return LoginResponse{AccessToken: token}, nil

	//	compare user.password with re.password

}

type ProfileRequest struct {
	UserID uint `json:"user_id"`
}
type ProfileResponse struct {
	Name string `json:"name"`
}

func (s Service) Profile(req ProfileRequest) (ProfileResponse, error) {
	// get user by ID
	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		//I have not expected the repository call return
		//"record not found" error,
		// because I assume the input is sanitized
		// TODO we can use Rich Error
		return ProfileResponse{}, fmt.Errorf("unexpected error : %w", err)
	}
	fmt.Println(ProfileResponse{Name: user.Name})
	return ProfileResponse{Name: user.Name}, nil

}

//type Claims struct {
//	RegisteredClaims jwt.RegisteredClaims
//	UserID           uint
//}
//
//func (c Claims) Valid() error {
//	return nil
//}

func createToken(userID uint, signKey string) (string, error) {
	// create a signer for rsa 256
	// TODO : replace with rsa 256 RS256
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}
