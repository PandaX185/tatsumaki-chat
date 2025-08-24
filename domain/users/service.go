package users

import (
	"crypto/sha256"
	"encoding/hex"
)

type UserService struct {
	repository UserRepository
}

func NewService(r UserRepository) *UserService {
	return &UserService{
		repository: r,
	}
}

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *UserService) Save(user User) (*User, error) {
	user.Password = hashPassword(user.Password)
	return s.repository.Save(user)
}

func (s *UserService) GetByExactUserName(username string) (*User, error) {
	return s.repository.GetByExactUserName(username)
}

func (s *UserService) SearchByUserName(username string) (UserSlice, error) {
	return s.repository.SearchByUserName(username)
}

func (s *UserService) Login(username, password string) (*User, error) {
	password = hashPassword(password)
	return s.repository.Login(username, password)
}
