package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	// Update DSN with your MySQL credentials
	db, err = sql.Open("mysql", "root:jspv@tcp(127.0.0.1:3306)/crud_demo")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Database is not reachable:", err)
	}
}

func addTask(w http.ResponseWriter, r *http.Request) {
	task := r.FormValue("_task")
	_, err := db.Query("insert into todolist(task) values(?)", task)
	if err != nil {
		fmt.Println("Error occured")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	fmt.Println("Successfully Printed")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func retriveTask(w http.ResponseWriter, r *http.Request) {
	var tasks []map[string]string
	var task string
	row, err := db.Query("select * from todolist")
	if err != nil {
		fmt.Println("Error Occured")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	defer row.Close()
	for row.Next() {
		row.Scan(&task)
		tasks = append(tasks, map[string]string{"Task": task})
	}
	temp := template.Must(template.ParseFiles("homepage.html"))
	temp.Execute(w, tasks)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	task := r.FormValue("_task")
	_, err := db.Query("delete from todolist where task=?", task)
	if err != nil {
		fmt.Println("Error Occurred")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	fmt.Println("Deleted successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getUpdateData(w http.ResponseWriter, r *http.Request) {

	value := r.FormValue("_task")
	temp := template.Must(template.ParseFiles("updateTask.html"))
	temp.Execute(w, map[string]string{"Task": value})
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	otask := r.FormValue("_otask")
	task := r.FormValue("_task")

	_, err := db.Query("update todolist set task=? where task=?", task, otask)
	if err != nil {
		fmt.Println("Error Occurred")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	fmt.Println("Updated Successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {

	http.HandleFunc("/", retriveTask)
	http.HandleFunc("/task.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "task.html")
	})
	http.HandleFunc("/add", addTask)
	http.HandleFunc("/delete", deleteTask)
	http.HandleFunc("/data", getUpdateData)
	http.HandleFunc("/update", updateTask)
	http.ListenAndServe(":8080", nil)
}
