package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"primerGoCrud/model"
	"text/template"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "sasa"
	dbName := "goblog"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, error := db.Query("select * from Empleado order by id desc")
	if error != nil {
		panic(error.Error())
	}
	emp := model.Empleado{}
	res := []model.Empleado{}
	for selDB.Next() {
		var id int
		var nombre, apellido, ciudad string
		error = selDB.Scan(&id, &nombre, &apellido, &ciudad)
		if error != nil {
			panic(error.Error())
		}

		emp.Id = id
		emp.Nombre = nombre
		emp.Apellido = apellido
		emp.Ciudad = ciudad
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func mostrar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, error := db.Query("select * from Empleado where id=?", nId)
	if error != nil {
		panic(error.Error())
	}
	emp := model.Empleado{}
	for selDB.Next() {
		var id int
		var nombre, apellido, ciudad string
		error = selDB.Scan(&id, &nombre, &apellido, &ciudad)
		if error != nil {
			panic(error.Error())
		}

		emp.Id = id
		emp.Nombre = nombre
		emp.Apellido = apellido
		emp.Ciudad = ciudad
	}
	tmpl.ExecuteTemplate(w, "Mostrar", emp)
	defer db.Close()

}

func Nuevo(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Nuevo", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		apellido := r.FormValue("apellido")
		ciudad := r.FormValue("ciudad")
		insertarForm, error := db.Prepare("insert into Empleado(nombre, apellido, ciudad) values(?,?,?)")
		if error != nil {
			panic(error.Error())
		}
		insertarForm.Exec(nombre, apellido, ciudad)
		log.Printf("Se ha insertado Empleado: " + nombre + ", " + apellido + "que vive en " + ciudad)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}
