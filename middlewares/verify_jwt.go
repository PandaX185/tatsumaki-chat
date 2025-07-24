package middlewares

import (
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
}

func VerifyJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range whitelist {
			if strings.HasPrefix(strings.Join([]string{r.Method, r.URL.Path}, ""),path) {
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
