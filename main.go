package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	comp "wiesel/pb175/components"
	data "wiesel/pb175/database"

	"github.com/a-h/templ"
)

type GlobalState struct {
    dbh *data.DBHandler
}


func (st *GlobalState) AddUser(w http.ResponseWriter, r *http.Request) {
    ID, err := data.SmallestMissingID(st.dbh.DB, st.dbh.Users)
    if err != nil {
        w.WriteHeader(500)
        comp.SignUp(err.Error()).Render(r.Context(), w)
        return
    }
    name := r.FormValue("name")
    email := r.FormValue("email")
    password := r.FormValue("psw")
    println(strconv.Itoa(ID) + name + email + password)

    user, err := data.NewUser(ID, name, email, password)
    if err != nil {
        w.WriteHeader(500)
        comp.SignUp(err.Error()).Render(r.Context(), w)
        return
    }
    err = st.dbh.InsertUser(user)
    if err != nil {
        w.WriteHeader(500)
        comp.SignUp(err.Error()).Render(r.Context(), w)
        return
    }

    users, err := st.dbh.GetUsers()
    if err != nil {
        w.WriteHeader(500)
        comp.SignUp(err.Error()).Render(r.Context(), w)
        return
    }
    comp.Users(users).Render(r.Context(), w)
}

func main() {
    dbh, err := data.InitDB()
    if err != nil {
        return
    }

    st := GlobalState {
        dbh: dbh,
    }

    mux := http.NewServeMux()
    mux.Handle("GET /", templ.Handler(comp.Index()))
    mux.Handle("GET /signup", templ.Handler(comp.SignUp("")))
    mux.HandleFunc("POST /signup", st.AddUser)
    srv := &http.Server {
        Addr: ":8080",
        Handler: mux,
        ErrorLog: log.Default(),
    }


    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan

        shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
        defer shutdownRelease()

        if err := srv.Shutdown(shutdownCtx); err != nil {
            log.Fatalf("HTTP shutdown error: %v", err)
        }
    }()

    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        fmt.Printf("err: %v\n", err)
    }
    fmt.Println("Shutdown")
}
