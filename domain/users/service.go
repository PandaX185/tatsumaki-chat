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
	res, err := s.repository.Save(user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *UserService) GetByUserName(username string) (*User, error) {
	res, err := s.repository.GetByUserName(username)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *UserService) Login(username, password string) (*User, error) {
	res, err := s.repository.Login(username, password)
	if err != nil {
		return nil, err
	}

	return res, nil
}
