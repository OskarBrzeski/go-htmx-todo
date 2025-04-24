package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var tmpl *template.Template

type Task struct {
	Id       int
	TaskName string
	Done     bool
}

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func initDB() {
	var err error

	db, err = sql.Open("mysql", "root:root@(127.0.0.1:3333)/testdb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
	defer db.Close()

	gRouter := mux.NewRouter()

	gRouter.HandleFunc("/", HomeHandler)
	gRouter.HandleFunc("/tasks", fetchTasks).Methods("GET")
	gRouter.HandleFunc("/new", newForm)
	gRouter.HandleFunc("/add", addTask).Methods("POST")
	gRouter.HandleFunc("/updateform/{id}", getTaskUpdateForm)
	gRouter.HandleFunc("/update/{id}", updateTask).Methods("PATCH")
	gRouter.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")

	fmt.Println("Listening on Port 3000")
	http.ListenAndServe(":3000", gRouter)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(),
			http.StatusInternalServerError)
	}
}

func fetchTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := getTasks(db)
	if err != nil {
		http.Error(w, "Error fetching tasks: "+err.Error(), http.StatusInternalServerError)
	}

	tmpl.ExecuteTemplate(w, "todoList", tasks)
}

func newForm(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "taskAddForm", nil)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	task := r.FormValue("task")

	query := "INSERT INTO tasks (task) VALUES (?);"

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, execErr := stmt.Exec(task)
	if execErr != nil {
		log.Fatal(err)
	}

	todos, _ := getTasks(db)

	tmpl.ExecuteTemplate(w, "todoList", todos)
}

func getTaskUpdateForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskId, _ := strconv.Atoi(vars["id"])

	task, err := getTaskById(db, taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.ExecuteTemplate(w, "taskUpdateForm", task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	taskname := r.FormValue("task")
	taskdone := r.FormValue("done")

	var doneBool bool

	switch strings.ToLower(taskdone) {
	case "yes", "on":
		doneBool = true
	default:
		doneBool = false
	}

	task := Task{id, taskname, doneBool}

	err := updateTaskById(db, &task)
	if err != nil {
		log.Fatal(err)
	}

	todos, _ := getTasks(db)

	tmpl.ExecuteTemplate(w, "todoList", todos)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	deleteTaskById(db, id)

	tasks, _ := getTasks(db)

	tmpl.ExecuteTemplate(w, "todoList", tasks)
}

func getTasks(dbPointer *sql.DB) ([]Task, error) {
	query := "SELECT id, task, done FROM tasks;"
	rows, err := dbPointer.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var todo Task

		rowErr := rows.Scan(&todo.Id, &todo.TaskName, &todo.Done)
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
	query := "SELECT id, task, done FROM tasks WHERE id = ?;"

	var task Task

	row := dbPointer.QueryRow(query, id)

	err := row.Scan(&task.Id, &task.TaskName, &task.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No task was found with id %d", id)
		}

		return nil, err
	}

	return &task, nil
}

func updateTaskById(dbPointer *sql.DB, task *Task) error {
	query := "UPDATE tasks SET task=?, done=? WHERE id=?;"

	result, err := dbPointer.Exec(query, task.TaskName, task.Done, task.Id)
	if err != nil {
		return err
	}

	rowsChanged, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsChanged == 0 {
		fmt.Println("No rows changed by update")
	}

	return nil
}

func deleteTaskById(dbPointer *sql.DB, id int) error {
	query := "DELETE FROM tasks WHERE id=?;"

	result, err := dbPointer.Exec(query, id)
	if err != nil {
		return err
	}

	rowsChanged, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsChanged == 0 {
		fmt.Println("No rows deleted")
	}

	return nil
}
