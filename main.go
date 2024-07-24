package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Todo struct {
	ID        int
	Name      string
	Timestamp time.Time
	status    string
}

var todos = []Todo{}
var nextID = 1

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/todos", createTodo).Methods("POST")
	router.HandleFunc("/todos", getTodos).Methods("GET")
	router.HandleFunc("/todos/{id:[0-9]+}", getTodo).Methods("GET")
	router.HandleFunc("/todos/{id:[0-9]+", updateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id:[0-9]+", deleteTodo).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}
func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo.ID = nextID
	todo.Timestamp = time.Now()
	nextID++
	todos = append(todos, todo)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}
func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for _, todo := range todos {
		if todo.ID == id {
			w.Header().Set("content-type", "application/json")
			json.NewEncoder(w).Encode(todo)
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var updateTodo Todo
	err = json.NewDecoder(r.Body).Decode(&updateTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Name = updateTodo.Name
			todos[i].status = updateTodo.status
			w.Header().Set("content-type", "application/json")
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:1], todos[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)

}
