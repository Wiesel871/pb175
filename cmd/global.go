package handlers

import (
	"fmt"
	_ "fmt"
	"net/http"
	"strconv"
	"time"
	comp "wiesel/pb175/components"
	data "wiesel/pb175/database"

	_"github.com/a-h/templ"
)

type GlobalState struct {
    DBH *data.DBHandler
    Anonym *data.User
}

const (
    session     = "NeUplneBazosSKSessionIdentifierCookie"
    sessionTimeout = 24 * time.Hour 
    cookieMaxAge   = int(sessionTimeout / time.Second)
)

func NewSession(id int) *http.Cookie {
    expiration := time.Now().Add(sessionTimeout)
    return &http.Cookie{
        Name:    session,
        Value:   strconv.Itoa(id),
        Expires: expiration,
        MaxAge:  cookieMaxAge,
        Path:    "/",
        HttpOnly: true,
    }
}

func GetClientID(r *http.Request) int {
    cookie, err := r.Cookie(session)
    if err != nil {
        return -1
    }
    id, err := strconv.Atoi(cookie.Value)
    if err != nil {
        return -1
    }
    return id
}

func LogOut(w http.ResponseWriter, r *http.Request) {
    expiration := time.Now().Add(-time.Hour)
    cookie := http.Cookie{
        Name:    session,
        Value:   "-1",
        Expires: expiration,
        MaxAge:  -1,
        Path:    "/",
    }
    http.SetCookie(w, &cookie)

    http.Redirect(w, r, "/", http.StatusFound)
}

func (st *GlobalState) SignUpForm(w http.ResponseWriter, r *http.Request) {
    id := GetClientID(r)
    fmt.Printf("id: %v\n", id)
    if id > -1 {
        http.SetCookie(w, NewSession(id))
        //http.Redirect(w, r, "/home", http.StatusMovedPermanently)
        println("already signed up")
        return
    }
    comp.Page(comp.SignUpForm(""), st.Anonym, comp.HomeN).Render(r.Context(), w)
}

func (st *GlobalState) LogInForm(w http.ResponseWriter, r *http.Request) {
    id := GetClientID(r)
    fmt.Printf("id: %v\n", id)
    if id > -1 {
        http.SetCookie(w, NewSession(id))
        //http.Redirect(w, r, "/home", http.StatusMovedPermanently)
        println("already signed up")
        return
    }
    comp.Page(comp.LogInForm(""), st.Anonym, comp.HomeN).Render(r.Context(), w)
}

func (st *GlobalState) Home(w http.ResponseWriter, r *http.Request) {
    id := GetClientID(r)
    user, err := st.DBH.GetUserById(id)
    if err != nil {
        user = st.Anonym
    }
    comp.Page(comp.IndexBody(), user, comp.HomeN).Render(r.Context(), w)
}


func SetupUserHandler(mux *http.ServeMux, st *GlobalState) {
    fmt.Printf("st.Anonym: %v\n", st.Anonym)
    mux.HandleFunc("/", st.Home)

    mux.HandleFunc("/signup", st.SignUpForm)
    mux.HandleFunc("POST /signup", st.SignUp)

    mux.HandleFunc("/login", st.LogInForm)
    mux.HandleFunc("POST /login", st.Login)

    mux.HandleFunc("GET /offers", st.GetOffers)
    mux.HandleFunc("GET /home", st.Home)


    mux.HandleFunc("GET /profile/{id}", st.Profile)
    mux.HandleFunc("GET /users", st.GetAllUsers)

    fs := http.FileServer(http.Dir("images"))
    mux.Handle("GET /images/", http.StripPrefix("/images/", fs))

    mux.HandleFunc("GET /profile/{id}/offers", st.GetUserOffers)
    mux.HandleFunc("GET /profile/{id_owner}/offers/{id}", st.GetOffer)

    mux.HandleFunc("POST /logout", LogOut)
}
