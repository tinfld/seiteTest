package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func jwtStart(w http.ResponseWriter, r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "Нет токена авторизации", http.StatusNotFound)
		return "", nil
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	return tokenString, nil

}

func jwtMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString, err := jwtStart(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		claim := &ClaimJwt{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claim,
			func(token *jwt.Token) (any, error) {
				return []byte(Env.Jwt_secret), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
		)
		if err != nil || !token.Valid {
			http.Error(w, "Неверный токен авторизации", http.StatusNetworkAuthenticationRequired)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func middleAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")

		if auth == "" {
			http.Error(w, "Нет токена авторизации", http.StatusNotFound)
			return
		}

		tockenString := strings.TrimPrefix(auth, "Bearer ")

		claim := &ClaimJwt{}

		token, err := jwt.ParseWithClaims(
			tockenString,
			claim,
			func(token *jwt.Token) (any, error) {
				return []byte(Env.Jwt_secret), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
		)
		if err != nil || !token.Valid || claim.Role != "admin" {
			http.Error(w, "неверный токен авторизации ", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)

	})
}
