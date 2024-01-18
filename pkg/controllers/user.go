package controllers

import (
	"github.com/PandaX185/tatsumaki-chat/pkg/models"
	"github.com/PandaX185/tatsumaki-chat/pkg/repository"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	repository.Repository
}

func (uc *UserController) SetupController(router *gin.Engine) {
	router.POST("/register", uc.CreateUser)
	router.POST("/login", uc.Login)
	router.GET("/users/:username", uc.GetUser)
	router.GET("/users/", uc.GetAllUsers)
}

func NewUserController(r repository.Repository) *UserController {
	return &UserController{
		Repository: r,
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := uc.Repository.CreateUser(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user created successfully"})
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.Repository.GetAllUsers()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

func (uc *UserController) GetUser(c *gin.Context) {
	email := c.Param("username")
	user, err := uc.Repository.GetUser(email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": user})
}

func (uc *UserController) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := uc.Repository.Login(user.Username, user.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
