package graph

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jaspeen/sample-users-app/config"
	"github.com/jaspeen/sample-users-app/db"
	"github.com/jaspeen/sample-users-app/services"
	"github.com/rs/zerolog/log"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			log.Debug().Msgf("Auth: %v", authHeader)
			// Allow unauthenticated users in
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenArr := strings.Split(authHeader, "Bearer ")
			if len(tokenArr) < 2 {
				services.GraphQLError(w, "Invalid authorization header", http.StatusForbidden)
				return
			}
			token := tokenArr[1]
			log.Debug().Msgf("Token: %v", token)

			jwtToken, err := services.ParseAndValidateToken(config.C.TokenSecret, token)
			if err != nil {
				log.Error().Err(err).Msgf("Token is not valid")
				if err.Error() == "Token is expired" {
					services.GraphQLError(w, "Token is expired", http.StatusForbidden)
				} else {
					services.GraphQLError(w, "Invalid token", http.StatusForbidden)
				}
				return
			}

			log.Debug().Msgf("Parsed token: %v", jwtToken)

			var user db.User

			if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
				services.ClaimsToUser(claims, &user)
			} else {
				services.GraphQLError(w, "Invalid token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func UserFromContext(ctx context.Context) *db.User {
	raw, _ := ctx.Value(userCtxKey).(*db.User)
	return raw
}
