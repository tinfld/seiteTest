package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Movie struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Time  int    `json:"time"`
	Gb    int    `json:"gb"`
	Image string `json:"image"`
}

func (s Storage) GetMovies(w http.ResponseWriter, r *http.Request) {
	var movies []Movie

	rows, err := s.db.Query(`select * from "movies"`)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		log.Println("Ошибка запроса чтения ")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var movie Movie
		err = rows.Scan(&movie.Id, &movie.Name, &movie.Time, &movie.Gb, &movie.Image)
		if err != nil {
			http.Error(w, "Ошибка сервера ", http.StatusInternalServerError)
			log.Println("Ошибка преобразования данных", err)
			return
		}
		movies = append(movies, movie)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, "Ошибка сервера ", http.StatusInternalServerError)
		log.Println("Ошибка конвертации", err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintln(w, string(data))
}

func (s Storage) PastMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Ошибка запроса ", http.StatusBadRequest)
		log.Println("Ошибка тела", err)
		return
	}

	_, err = s.db.Exec(`insert into movies (name, time, gb, images) values ($1, $2, $3, $4)`, movie.Name, movie.Time, movie.Gb, movie.Image)
	if err != nil {
		http.Error(w, "Ошибка сервера ", http.StatusInternalServerError)
		log.Println("ошибка записи бд ", err)
		return
	}
	fmt.Fprintln(w, "Запись успешна")
}

func (s Storage) GetMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie

	id := r.PathValue("id")
	nomer, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Ошибка id ", http.StatusBadRequest)
		log.Println("Ошибка чтения", err)
		return
	}

	row := s.db.QueryRow(`select id, name, time, gb from movies where id = $1`, nomer)
	err = row.Scan(&movie.Id, &movie.Name, &movie.Time, &movie.Gb)
	if err != nil {
		http.Error(w, "Ошибка сервера ", http.StatusInternalServerError)
		log.Println("Ошибка базы ", err)
		return
	}

	data, err := json.Marshal(movie)
	if err != nil {
		http.Error(w, "Ошибка сервера ", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintln(w, string(data))
}

func (s Storage) deleteMovie(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	nomer, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Ошибка сервера ", http.StatusInternalServerError)
		log.Println("Ошибка чтения", err)
		return
	}

	_, err = s.db.Exec(`delete from movies where id = $1`, nomer)
	if err != nil {
		log.Println("Ошибка удаления ", err)
		return
	}

	fmt.Fprintln(w, "Удаление успешно")
}

func (s Storage) PutMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie

	id := r.PathValue("id")
	nomer, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Ошибка чтения", err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		log.Println("Ошибка чтнеия", err)
		return
	}

	_, err = s.db.Exec(`update movies set name = $1, time = $2, gb = $3 where id = $4`, movie.Name, movie.Time, movie.Gb, nomer)
	if err != nil {
		log.Println("Ошибка изменения", err)
		return
	}

	fmt.Fprintln(w, "Изменение успешно")
}
