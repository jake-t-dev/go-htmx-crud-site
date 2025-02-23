package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var tmpl *template.Template

type Task struct {
	Id   int
	Task string
	Done bool
}

func init() {
	var err error

	tmpl, err = template.ParseGlob("templates/*.html")

	if err != nil {
		log.Fatal(err)
	}
}

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:root@(127.0.0.1:3333)/testdb")

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	initDB()
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", HomeHandler)

	// get tasks
	router.HandleFunc("/tasks", fetchTasks).Methods("GET")

	//fetch add task form
	router.HandleFunc("/getnewtaskform", getTaskForm)

	//add task
	router.HandleFunc("/tasks", addTask).Methods("POST")

	//fetch update task form
	router.HandleFunc("/gettaskupdateform/{id}", getTaskUpdateForm)

	//update task
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT", "POST")

	//delete task
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	http.ListenAndServe(":3000", router)

}

// Route Handlers

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "home.html", nil)

	if err != nil {
		http.Error(w, "Error executing teamplate: "+err.Error(), http.StatusInternalServerError)
	}
}

func fetchTasks(w http.ResponseWriter, r *http.Request) {
	todos, _ := getTasks(db)
	tmpl.ExecuteTemplate(w, "todoList", todos)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	task := r.FormValue("task")

	query := "INSERT INTO tasks (task) VALUES (?)"

	stmt, err := db.Prepare(query)

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, executeErr := stmt.Exec(task)

	if executeErr != nil {
		log.Fatal(executeErr)
	}

	todos, _ := getTasks(db)

	tmpl.ExecuteTemplate(w, "todoList", todos)
}

func getTaskForm(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "addTaskForm", nil)
}

func getTaskUpdateForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskId, _ := strconv.Atoi(vars["id"])

	task, err := getTaskById(db, taskId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	tmpl.ExecuteTemplate(w, "updateTaskForm", task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskItem := r.FormValue("task")
	isDone := r.FormValue("done")

	var taskStatus bool

	switch strings.ToLower(isDone) {
	case "yes", "on":
		taskStatus = true
	case "no", "off":
		taskStatus = false
	default:
		taskStatus = false
	}

	taskId, _ := strconv.Atoi(vars["id"])

	task := Task{
		Id:   taskId,
		Task: taskItem,
		Done: taskStatus,
	}

	query := "UPDATE tasks SET task = ?, done = ? WHERE id = ?"

	result, err := db.Exec(query, task.Task, task.Done, task.Id)

	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		http.Error(w, "No task found with id "+strconv.Itoa(task.Id), http.StatusNotFound)
	}

	todos, _ := getTasks(db)

	tmpl.ExecuteTemplate(w, "todoList", todos)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskId, _ := strconv.Atoi(vars["id"])

	query := "DELETE FROM tasks WHERE id = ?"

	stmt, err := db.Prepare(query)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, deleteError := stmt.Exec(taskId)

	if deleteError != nil {
		log.Fatal(deleteError)
	}

	todos, _ := getTasks(db)

	tmpl.ExecuteTemplate(w, "todoList", todos)

}

// Utility Functions

func getTasks(dbPointer *sql.DB) ([]Task, error) {
	query := "SELECT id, task, done FROM tasks"

	rows, err := dbPointer.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var todo Task

		rowErr := rows.Scan(&todo.Id, &todo.Task, &todo.Done)

		if rowErr != nil {
			return nil, rowErr
		}

		tasks = append(tasks, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func getTaskById(dbPointer *sql.DB, id int) (*Task, error) {
	query := "SELECT id, task, done FROM tasks WHERE id = ?"

	var task Task

	row := dbPointer.QueryRow(query, id)

	err := row.Scan(&task.Id, &task.Task, &task.Done)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No task found with id %d", id)
		}
		return nil, err
	}

	return &task, nil
}
