package userservice

import (
	"GameApp/entity"
	"GameApp/param"
	"GameApp/pkg/richerror"
	"GameApp/repository/mysql"
	"fmt"
)

type Repository interface {
	RegisterUser(u entity.User) (entity.User, error)
	GetUserByPhone(phoneNumber string) (entity.User, bool, error)
	GetUserByID(userID uint) (entity.User, error)
}

type AuthService interface {
	CreateRefreshToken(user entity.User) (string, error)
	CreateAccessToken(user entity.User) (string, error)
}

type Service struct {
	repo Repository
	auth AuthService
}

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {

	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    mysql.GetMD5Hash(req.Password),
	}
	// create new user in storage
	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return param.RegisterResponse{}, fmt.Errorf("failed to register user: %w", err)
	}

	//	 return created user
	return param.RegisterResponse{User: createdUser}, nil

}

func New(repo Repository, auth AuthService) Service {
	return Service{repo: repo, auth: auth}

}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User   param.UserInfo `json:"user"`
	Tokens Tokens         `json:"tokens"`
}

//RefreshToken string `json:"refresh_token"`

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	const op = "userservice.Login"
	//  check the existence of phone number from repository
	//	get the user by phone number
	// TODO : Its better to  separate method for check existence check  of get user by phone number
	user, exists, err := s.repo.GetUserByPhone(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, richerror.New(op).
			WithMessage("failed to get user by phone number").
			WithMeta(map[string]interface{}{"Req": req}).WithWrappedError(err)
	}
	if !exists {
		return LoginResponse{}, fmt.Errorf("user or password not valid")
	}
	if user.Password != mysql.GetMD5Hash(req.Password) {
		return LoginResponse{}, fmt.Errorf("user or password not valid ")

	}
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to create token: %w", err)
	}
	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to refresh token: %w", err)
	}

	return LoginResponse{
		User: param.UserInfo{
			ID:          user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
		Tokens: Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil

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
		return ProfileResponse{}, richerror.New("userService.Profile").WithWrappedError(err)
	}
	fmt.Println(ProfileResponse{Name: user.Name})
	return ProfileResponse{Name: user.Name}, nil

}
