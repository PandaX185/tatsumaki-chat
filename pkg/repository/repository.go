package repository

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/PandaX185/tatsumaki-chat/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

type Repository interface {
	GetUser(username string) (*models.User, error)
	CreateUser(user *models.User) error
	Login(username string, password string) (string, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetUser(username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow("SELECT * FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (r *repository) CreateUser(user *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", &user.Username, &user.Password)

	if err != nil {
		return err
	}

	return nil

}

func (r *repository) Login(username string, password string) (string, error) {
	var token string
	user, err := r.GetUser(username)
	if err != nil {
		return "", errors.New("invalid email")
	}

	if user.Password != password {
		return "", errors.New("incorrect password")
	}

	if err != nil {
		return "", err
	}

	token, err = createJWTToken(username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func createJWTToken(username string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiry time (24 hours from now)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
