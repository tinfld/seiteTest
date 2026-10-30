package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func main() {

	loger, err := loadLog()
	if err != nil {
		fmt.Println("Ошибка создания логов", err)
		return
	}
	defer loger.close()
	loger.LogMain.Println("Логи созданы")

	env, err := EnvLoad(loger)
	if err != nil {
		loger.LogMain.Println("Ошибка загрузкм env", err)
		return
	}
	loger.LogMain.Println("Успешная загрузка env")

	sqlData := fmt.Sprintf("postgres://postgres:%v@%v:%v/%v", env.PasswordDb, env.DbHost, env.PortDb, env.Namedb)
	db, err := sql.Open("pgx", sqlData)
	if err != nil {
		loger.LogMain.Println("Ошибка загрузки БД", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		loger.LogMain.Println("База не доступна", err)
		return
	}
	loger.LogMain.Println("Успешное подключение к базе")
	s := Storage{db: db}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":3030",
		Handler: mux,
	}

	mux.Handle("GET /", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("GET /movies", s.GetMovies)
	mux.Handle("POST /createmovie", jwtMiddle(http.HandlerFunc(s.PastMovie)))
	mux.HandleFunc("GET /movie/{id}", s.GetMovie)
	mux.Handle("DELETE /movie/{id}", jwtMiddle(http.HandlerFunc(s.deleteMovie)))
	mux.Handle("PUT /movie/{id}", jwtMiddle(http.HandlerFunc(s.PutMovie)))

	mux.Handle("POST /userreg", http.HandlerFunc(s.AddUser))
	mux.Handle("DELETE /deluser/{id}", middleAdmin(http.HandlerFunc(s.deletUser)))
	mux.HandleFunc("POST /loginuser", s.loginUser)
	mux.Handle("GET /users", middleAdmin(http.HandlerFunc(s.GetUsers)))

	loger.LogMain.Println("Запуск Сервера")
	err = server.ListenAndServe()
	if err != nil {
		loger.LogMain.Println("неудача в создании сервера ", err)
		return
	}
}
