package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
    "database/sql"
    "os"
    "os/signal"
    "syscall"
    "time"

	"github.com/a-h/templ"
    comp "wiesel/pb175/components"
    data "wiesel/pb175/database"
)

type GlobalState struct {
    db *sql.DB 
    db_s int
}

func (st *GlobalState) AddUser(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("name")
    email := r.FormValue("email")
    password := r.FormValue("psw")
    println(name + email + password)

    user, err := data.NewUser(st.db_s + 1, name, email, password)
    if err != nil {
        w.WriteHeader(500)
        comp.SignUp(err.Error()).Render(r.Context(), w)
        return
    }
    err = data.InsertUser(st.db, user)
    if err != nil {
        w.WriteHeader(500)
        comp.SignUp(err.Error()).Render(r.Context(), w)
        return
    }

    comp.Users(data.GetUsers(st.db)).Render(r.Context(), w)

}

func main() {
    db, error := data.InitDB()
    if error != nil {
        return
    }

    st := GlobalState {
        db: db,
    }

    mux := http.NewServeMux()
    mux.Handle("GET /", templ.Handler(comp.Index()))
    mux.Handle("GET /signup", templ.Handler(comp.SignUp("")))
    mux.HandleFunc("POST /signup", st.AddUser)
    srv := &http.Server {
        Addr: ":9080",
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
