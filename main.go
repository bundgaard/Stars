package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

type Star struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func handleRequest(engine *xorm.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			handlePOST(engine, w, r)
			return
		}
		handleGET(engine, w, r)
	}
}

func handleGET(engine *xorm.Engine, w http.ResponseWriter, r *http.Request) {
	var starz []Star
	if err := engine.Cols("id", "name").Find(&starz); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&starz); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handlePOST(engine *xorm.Engine, w http.ResponseWriter, r *http.Request) {

	var star Star
	if err := json.NewDecoder(r.Body).Decode(&star); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fmt.Println(star.ID)
	if star.ID == "" {
		star.ID = uuid.NewString()
	}
	fmt.Println(star)
	if _, err := engine.Insert(star); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func main() {

	fmt.Println("Starz")
	orm, err := xorm.NewEngine("sqlite3", "./starz.db")
	orm.SetMapper(names.GonicMapper{})
	if err != nil {
		log.Fatal(err)
	}
	err = orm.Sync2(new(Star))
	if err != nil {
		log.Fatal(err)
	}

	root := http.NewServeMux()
	root.HandleFunc("/", handleRequest(orm))
	log.Fatal(http.ListenAndServe(":8080", root))

}
