package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Segment struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Userid int    `json:"userid"`
}

func main() {
	// соединение с БД
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	//создать таблицу
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS segments (id SERIAL PRIMARY KEY, name VARCHAR(50) NOT NULL, userid INT NOT NULL, CONSTRAINT add UNIQUE (name, userid))")

	if err != nil {
		panic(err)
	}

	// маршрутизаторы
	router := mux.NewRouter()
	router.HandleFunc("/segments/{userid}", getUserSegments(db)).Methods("GET")
	router.HandleFunc("/segments/delete/{name}/{userid}", deleteSegment(db)).Methods("DELETE")
	router.HandleFunc("/segments/deleteall/{name}", deleteOneSegment(db)).Methods("DELETE")
	router.HandleFunc("/segments/add", addSegment(db)).Methods("POST")

	// запуск сервера
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))

}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

//найти все сегменты по id пользователя

func getUserSegments(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userid := vars["userid"]
		rows, err := db.Query("SELECT * FROM segments WHERE userid = $1", userid)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		segments := []Segment{}
		for rows.Next() {
			var s Segment
			if err := rows.Scan(&s.ID, &s.Name, &s.Userid); err != nil {
				panic(err)
			}
			segments = append(segments, s)
		}
		if err := rows.Err(); err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(segments)
	}
}

// cоздать и добавить сегмент конкретному юзеру

func addSegment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var s Segment
		json.NewDecoder(r.Body).Decode(&s)

		err := db.QueryRow("INSERT INTO segments (name, userid) VALUES ($1, $2) RETURNING id", s.Name, s.Userid).Scan(&s.ID)
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(s)
	}
}

//удалить конкретный сегмент у конкретного юзера

func deleteSegment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		userid := vars["userid"]

		var s Segment
		err := db.QueryRow("SELECT * FROM segments WHERE name =$1 AND userid =$2", name, userid).Scan(&s.ID, &s.Name, &s.Userid)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM segments WHERE name =$1 AND userid =$2", name, userid)
			if err != nil {
				panic(err)
			}
			json.NewEncoder(w).Encode("Segment deleted")
		}

	}
}

//удалить сегмент полностью из БД

func deleteOneSegment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		var s Segment
		err := db.QueryRow("SELECT * FROM segments WHERE name =$1", name).Scan(&s.ID, &s.Name, &s.Userid)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM segments WHERE name =$1", name)
			if err != nil {
				panic(err)
			}
			json.NewEncoder(w).Encode("Segment deleted from all users")
		}

	}
}
