package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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
	selDB, error := db.Query("select * from Employee order by id desc")
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
