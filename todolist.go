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
	"strconv"
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


func Store(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	log.WithFields(log.Fields{"description": description}).Info("New Todo Added.")
	todo := &TodoItemModel{Description: description, Completed: false}
	db.Create(&todo)
	result := db.Last(&todo)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
}

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	response := GetItemById(id)

	if ! response {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": false, "error:" "Record not found"}`)
	} else {
		completed, _ := strconv.ParseBool(r.FormValue("completed"))

		log.WithFields(log.Fields{"id": id, "Completed": completed}).Info("Updating Todod")
		todo := &TodoItemModel{}
		db.First(&todo, id)
		todo.Completed = completed
		db.Save(&todo)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": true`)
	}
}

func DeletedTodoItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	response := GetItemById(id)
	if ! response {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": false, "error": "Record not Found`)
	} else {
		log.WithFields(log.Fields{"id": id}).Info("Deleting TodoItem")

		todo := &TodoItemModel{}
		db.First(&todo, id)
		db.Delete(&todo)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": true`)
	}
}

func GetItemById(Id int) bool {
	todo := &TodoItemModel{}
	result := db.First(&todo, Id)

	if result.Error != nil {
		log.Warn("Item not found")
		return false
	}

	return true
}

func GetCompletedTodoItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get completed Items")
	CompletedTodoItems := GetTodoItems(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CompletedTodoItems)
}

func GetTodoItems(completed bool) interface{} {
	var todos []TodoItemModel
	TodoItems := db.Where("completed = ?", completed).Find(&todos).Value
	return TodoItems
}

func GetIncompleteTodoItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get Incomplete TodoItems")
	IncompleteTodoItems := GetTodoItems(false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(IncompleteTodoItems)
}

func main() {
	defer db.Close()
	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})
	log.Info("Starting TodoList API server")
	router := mux.NewRouter()
	router.HandleFunc("/test", ApiHealth).Methods("GET")
	router.HandleFunc("/todo", Store).Methods("POST")
	http.ListenAndServe(":8000", router)

	if err != nil {
		log.Println("Error when opening DB", err)
	}

	db.Close()
}