package main

import (
	"fmt"
	"io"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux" // for creating router.
	log "github.com/sirupsen/logrus"
	"encoding/json"
)

const (
	username = "root"
	password = ""
	dbName = "goTodo"
)


func ApiHealth(w http.ResponseWriter, r *http.Request) {
	log.Info("API health is ok")
	w.Header().Set("content-type", "application/json")
	io.WriteString(w, `{alive: true}`)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func dsn(dbname string) string {
	return fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True&loc=Local", username, dbName)
}

var db , err  = gorm.Open("mysql", dsn("goTodo"))

type TodoItemModel struct {
	Id int `gorm:"primary_key"`
	Description string
	Completed bool
}

func main() {
	defer db.Close()
	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})
	log.Info("Starting TodoList API server")
	router := mux.NewRouter()
	router.HandleFunc("/health", ApiHealth).Methods("Get")
	http.ListenAndServe(":8000", router)

	if err != nil {
		log.Println("Error when opening DB", err)
	}

	db.Close()
}