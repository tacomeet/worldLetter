package main

import (
	"github.com/gorilla/mux"
	"github.com/tacomea/worldLetter/database"
	"github.com/tacomea/worldLetter/repository"
	"github.com/tacomea/worldLetter/usecase"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tpl *template.Template

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tpl = template.Must(template.ParseGlob("templates/*html"))
}

func main() {
	db := database.NewPostgresDB()
	ur := repository.NewUserRepositoryPG(db)
	sr := repository.NewSessionRepositoryPG(db)
	lr := repository.NewLetterRepositoryPG(db)

	uu := usecase.NewUserUsecase(ur)
	su := usecase.NewSessionUsecase(sr)
	lu := usecase.NewLetterUsecase(lr)

	// handlers
	h := newHandler(uu, su, lu)
	r := mux.NewRouter()

	//private routes
	r.HandleFunc("/", h.jwtAuth(h.indexHandler)).Methods("GET")
	r.HandleFunc("/create", h.jwtAuth(h.createHandler)).Methods("GET")
	r.HandleFunc("/send", h.jwtAuth(h.sendHandler)).Methods("POST")
	r.HandleFunc("/letter/received", h.jwtAuth(h.letterReceivedHandler)).Methods("GET")
	r.HandleFunc("/letter/sent", h.jwtAuth(h.letterSentHandler)).Methods("GET")
	r.HandleFunc("/admin", h.jwtAuth(h.adminHandler)).Methods("GET")

	// public routes
	r.HandleFunc("/signin", h.signinHandler).Methods("GET")
	r.HandleFunc("/signup", h.signupHandler).Methods("GET")
	r.HandleFunc("/register", h.registerHandler).Methods("POST")
	r.HandleFunc("/login", h.loginHandler).Methods("POST")
	r.HandleFunc("/logout", h.logoutHandler).Methods("POST")

	r.Handle("/favicon.ico", http.NotFoundHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatalln(http.ListenAndServe(":"+port, r))
}
