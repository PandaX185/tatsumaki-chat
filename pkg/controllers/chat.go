package controllers

import (
	"errors"
	"net/http"
	"os"

	"github.com/PandaX185/tatsumaki-chat/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatController struct {
	repository.Repository
	KConn *kafka.Conn
}

func (c *ChatController) SetupController(router *gin.Engine) {
	router.GET("ws/chat", c.OpenChat)
}

func NewChatController(r repository.Repository, kConn *kafka.Conn) *ChatController {
	return &ChatController{
		Repository: r,
		KConn:      kConn,
	}
}

func (c *ChatController) OpenChat(ctx *gin.Context) {
	userToChat := ctx.Query("user")
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	token := ctx.GetHeader("Authorization")
	username, err := extractUsernameFromToken(token)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.KConn.CreateTopics(kafka.TopicConfig{
		Topic:             userToChat + ":" + username,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})

	batch := c.KConn.ReadBatch(10e3, 1e6)
	b := make([]byte, 10e3)
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		ctx.JSON(http.StatusOK, string(b[:n]))
	}

	defer batch.Close()
}

func extractUsernameFromToken(tokenString string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if username, exists := claims["username"].(string); exists {
			return username, nil
		}
	}

	return "", errors.New("failed to extract username from token")
}
