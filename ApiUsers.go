package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ClaimJwt struct {
	UserEmail string
	Role      string
	jwt.RegisteredClaims
}

func (s Storage) AddUser(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Ошибка чтения тела", err)
		http.Error(w, "Ошибка чтения тела", http.StatusInternalServerError)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка ожидания сервера", http.StatusInternalServerError)
		log.Println("Ошибка записи пароля", err)
		return
	}
	user.Password = string(hash)

	var test string
	row := s.db.QueryRow(`select email from users where email = $1`, user.Email)
	err = row.Scan(&test)
	if err != nil {
		log.Println(err)
		_, err = s.db.Exec(`insert into users (name, email, password) values ($1, $2, $3)`, user.Name, user.Email, user.Password)
		if err != nil {
			http.Error(w, "Ошибка ожидания сервера ", http.StatusInternalServerError)
			log.Println("неудачный запрос", err)
			return
		}

		fmt.Fprintln(w, "Регистрация успешна")
	}

	fmt.Fprintln(w, "Такой email уже есть")

}

func (s Storage) deletUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	nomer, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Ошибка запроса", http.StatusBadGateway)
		log.Println(err)
		return
	}

	_, err = s.db.Exec(`delete from users where id = $1`, nomer)
	if err != nil {
		http.Error(w, "Ошибка ожидания сервера", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	fmt.Fprintln(w, "Удаление успешно", http.StatusOK)
}

func (s Storage) loginUser(w http.ResponseWriter, r *http.Request) {
	var user User
	var loginUser User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка тела", http.StatusBadRequest)
		log.Println(err)
		return
	}

	row := s.db.QueryRow(`select id, name, email, password from users where email = $1`, user.Email)
	err = row.Scan(&loginUser.Id, &loginUser.Name, &loginUser.Email, &loginUser.Password)
	if err != nil {
		http.Error(w, "неверный пароль или email", http.StatusNotFound)
		log.Println(err)
		return
	}

	if user.Email == loginUser.Email {
		err := bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(user.Password))
		if err != nil {
			http.Error(w, "неверный пароль или email", http.StatusNotFound)
			return
		}

		claims := &ClaimJwt{
			UserEmail: user.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(
					time.Now().Add(15 * time.Minute),
				),
			},
		}
		if user.Email == "admin@piprch.ru" {
			claims.Role = "admin"
		} else {
			claims.Role = "user"
		}

		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			claims,
		)

		tokenString, err := token.SignedString([]byte(Env.Jwt_secret))
		if err != nil {
			http.Error(w, "Ошибка ожидания сервера", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		fmt.Fprintln(w, "Вход успешен, ваш ключ", string(tokenString))
	}
}

func (s Storage) GetUsers(w http.ResponseWriter, r *http.Request) {
	var userall []User
	var user User
	rows, err := s.db.Query(`select id, name, email from users`)
	if err != nil {
		http.Error(w, "Ошибка ожидания сервера", http.StatusInternalServerError)
		log.Println("не доставли пользователей ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			log.Println("Ошибка чтения ", err)
			return
		}
		userall = append(userall, user)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := json.Marshal(userall)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, string(data))
}
