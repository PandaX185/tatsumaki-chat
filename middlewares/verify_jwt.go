package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/PandaX185/tatsumaki-chat/domain/errors/codes"
	"github.com/golang-jwt/jwt/v5"
)

var whitelist = []string{
	"POST/api/users",
	"POST/api/users/login",
	"GET/api/realtime",
}

func VerifyJwtFromQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			w.WriteHeader(codes.UNAUTHORIZED)
			json.NewEncoder(w).Encode(map[string]any{
				"code":    codes.UNAUTHORIZED,
				"message": "Token query parameter is missing",
			})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrECDSAVerification
			}
			return []byte(os.Getenv("AUTH_KEY")), nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(codes.UNAUTHORIZED)
			json.NewEncoder(w).Encode(map[string]any{
				"code":    codes.UNAUTHORIZED,
				"message": "Invalid token",
			})
			return
		}

		claims := extractClaims(token)
		ctx := context.WithValue(context.Background(), "userId", claims["userId"])
		ctx = context.WithValue(ctx, "username", claims["username"])
		ctx = context.WithValue(ctx, "fullname", claims["fullname"])
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func VerifyJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range whitelist {
			if strings.HasPrefix(strings.Join([]string{r.Method, r.URL.Path}, ""), path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(codes.UNAUTHORIZED)
			json.NewEncoder(w).Encode(map[string]any{
				"code":    codes.UNAUTHORIZED,
				"message": "Authorization header is missing",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrECDSAVerification
			}
			return []byte(os.Getenv("AUTH_KEY")), nil
		})

		claims := extractClaims(token)
		ctx := context.WithValue(context.Background(), "userId", claims["userId"])
		ctx = context.WithValue(ctx, "username", claims["username"])
		ctx = context.WithValue(ctx, "fullname", claims["fullname"])
		r = r.WithContext(ctx)

		if err != nil || !token.Valid {
			w.WriteHeader(codes.UNAUTHORIZED)
			json.NewEncoder(w).Encode(map[string]any{
				"code":    codes.UNAUTHORIZED,
				"message": "Authorization header is missing",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func extractClaims(token *jwt.Token) jwt.MapClaims {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims
	}
	return nil
}
