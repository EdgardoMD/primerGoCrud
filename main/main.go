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

var tmpl = template.Must(template.ParseGlob("template/*"))

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

func Mostrar(w http.ResponseWriter, r *http.Request) {
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

func Editar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Empleado WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := model.Empleado{}
	for selDB.Next() {
		var id int
		var nombre, apellido, ciudad string
		err = selDB.Scan(&id, &nombre, &apellido, &ciudad)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Nombre = nombre
		emp.Apellido = apellido
		emp.Ciudad = ciudad
	}
	tmpl.ExecuteTemplate(w, "Editar", emp)
	defer db.Close()
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		apellido := r.FormValue("apellido")
		ciudad := r.FormValue("ciudad")
		id := r.FormValue("uid")
		insertarForm, error := db.Prepare("UPDATE Empleado SET nombre=?, apellido=?, ciudad=? WHERE id = ?")
		if error != nil {
			panic(error.Error())
		}
		insertarForm.Exec(nombre, apellido, ciudad, id)
		log.Printf("Se ha actualizado al Empleado: " + nombre + ", " + apellido + "que vive en " + ciudad)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)

}

func Eliminar(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	eliminar, err := db.Prepare("delete from Empleado where id = ?")
	if err != nil {
		panic(err.Error())
	}
	eliminar.Exec(emp)
	log.Println("Eliminado")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Servidor corriendo en puerto: http://localhost:8091")
	http.HandleFunc("/", Index)
	http.HandleFunc("/mostrar", Mostrar)
	http.HandleFunc("/nuevo", Nuevo)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/eliminar", Eliminar)
	http.ListenAndServe(":8091", nil)
}
