package user

import (
	"database/sql"

	"github.com/jimvid/dionysus/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	GetUser(username string) (model.User, error)
	InsertUser(model.User) error
}

type UserService struct {
	repo *UserRepository
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		repo: NewUserRepository(db),
	}
}

func (s *UserService) NewUserWithHashedPassword(registerUser model.UserRequest) (model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		Username:     registerUser.Username,
		Email:        registerUser.Email,
		PasswordHash: string(hashedPassword)}, nil
}

func (s *UserService) DoesUserExist(username string) (bool, error) {
	doesUserExist, err := s.repo.DoesUserExist(username)

	if err != nil {
		return false, err

	}

	return doesUserExist, nil
}

func (s *UserService) InsertUser(user model.User) error {
	err := s.repo.InsertUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUser(username string) (model.User, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) ValidatePassword(hashedPassowrd, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassowrd), []byte(password))
	return err == nil
}
