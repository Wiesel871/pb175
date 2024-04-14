package handlers

import (
    "net/http"
	"github.com/a-h/templ"
    "time"
    "strconv"
    "fmt"
    data "wiesel/pb175/database"
    comp "wiesel/pb175/components"
)

type GlobalState struct {
    DBH *data.DBHandler
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

func LogedInOwned(r *http.Request, u *data.User) (int, bool) {
    cookie, err := r.Cookie(session)
    if err == nil {
        id, err := strconv.Atoi(cookie.Value)
        fmt.Printf("err: %v\n", err)
        fmt.Printf("id: %v\n", id)
        return id, err == nil && u != nil && id == u.ID
    }
    println("no cookie")
    return -1, false

}


func Redirect(w http.ResponseWriter, r *http.Request) {
    id, _ := LogedInOwned(r, nil)
    if id > -1 {
        http.Redirect(w, r, "/home", http.StatusPermanentRedirect)
        return
    }
    http.SetCookie(w, NewSession(id))
    comp.Index().Render(r.Context(), w)
}



func SetupUserHandler(mux *http.ServeMux, st *GlobalState) {
    mux.HandleFunc("GET /", Redirect)
    mux.HandleFunc("GET /home", st.GetOffers)
    mux.Handle("GET /signup", templ.Handler(comp.SignUp()))
    mux.HandleFunc("POST /signup", st.AddUser)
    mux.Handle("GET /login", templ.Handler(comp.LogIn()))
    mux.HandleFunc("POST /login", st.Login)
    mux.HandleFunc("GET /profile/{id}", st.Profile)
    mux.HandleFunc("GET /users", st.GetAllUsers)

    fs := http.FileServer(http.Dir("images"))
    mux.Handle("GET /images/", http.StripPrefix("/images/", fs))

    mux.HandleFunc("GET /profile/{id}/offers", st.GetUserOffers)
    mux.HandleFunc("GET /profile/{id_owner}/offers/{id}", st.GetOffer)
}
