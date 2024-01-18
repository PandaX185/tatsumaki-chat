package controllers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/PandaX185/tatsumaki-chat/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatController struct {
	repository.Repository
}

func (c *ChatController) SetupController(router *gin.Engine) {
	router.GET("ws/chat", c.OpenChat)
	router.POST("send-message", c.SendMessage)
}

func NewChatController(r repository.Repository) *ChatController {
	return &ChatController{
		Repository: r,
	}
}

func (c *ChatController) SendMessage(ctx *gin.Context) {
	var body map[string]string
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	token := ctx.GetHeader("Authorization")
	username, err := extractUsernameFromToken(token)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	topic := username + "-" + body["user"]
	kConn, err := kafka.DialLeader(context.Background(), "tcp", os.Getenv("KAFKA_URL"), topic, 0)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := kConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, err := kConn.WriteMessages(kafka.Message{
		Topic: topic,
		Value: []byte(body["message"]),
	}); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Message sent successfully",
	})
}

func (c *ChatController) OpenChat(ctx *gin.Context) {
	userToChat := ctx.Query("user")
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Error in upgrading connection: "+err.Error()))
		return
	}

	token := ctx.Query("Authorization")
	username, err := extractUsernameFromToken(token)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Error in upgrading connection: "+err.Error()))
		return
	}

	topic := username + "-" + userToChat
	kConn, err := kafka.DialLeader(context.Background(), "tcp", os.Getenv("KAFKA_URL"), topic, 0)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Error in upgrading connection: "+err.Error()))
		return
	}

	kConn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	for {
		msg, readErr := kConn.ReadMessage(1e6)
		if errors.Is(readErr, io.EOF) {
			break
		}

		if readErr != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Error in reading message: "+readErr.Error()))
			break
		}

		if len(msg.Value) != 0 {
			conn.WriteMessage(websocket.TextMessage, []byte(userToChat+": "+string(msg.Value)))
		}
	}

	defer func() {
		conn.Close()
		kConn.Close()
	}()
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
