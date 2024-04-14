package handlers

import (
    "net/http"
	"github.com/a-h/templ"
    data "wiesel/pb175/database"
    comp "wiesel/pb175/components"
)

type GlobalState struct {
    DBH *data.DBHandler
}

func SetupUserHandler(mux *http.ServeMux, st *GlobalState) {
    mux.Handle("GET /", templ.Handler(comp.Index()))
    mux.Handle("GET /signup", templ.Handler(comp.SignUp()))
    mux.HandleFunc("POST /signup", st.AddUser)
    mux.Handle("GET /login", templ.Handler(comp.LogIn()))
    mux.HandleFunc("POST /login", st.Login)
    mux.HandleFunc("GET /profile/{id}", st.Profile)
    mux.HandleFunc("GET /users", st.GetAllUsers)

    fs := http.FileServer(http.Dir("images"))
    mux.Handle("GET /images/", http.StripPrefix("/images/", fs))

}
