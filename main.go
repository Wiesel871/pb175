package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "wiesel/pb175/components"
	data "wiesel/pb175/database"
	cmd "wiesel/pb175/cmd"
    state "wiesel/pb175/state"

	_ "github.com/a-h/templ"

)

func main() {
    dbh, err := data.InitDB()
    if err != nil {
        return
    }

    anonym := data.User{
        ID: -1,
        Name: "",
        Email: "",
        Password: []byte(""),
        IsAdmin: false,
        HasPFP: false,
    }
    st := state.GlobalState {
        DBH: dbh,
        Anonym: &anonym,
    }

    mux := http.NewServeMux()
    cmd.SetupUserHandler(mux, &st)
    st.SRV = &http.Server {
        Addr: ":8070",
        Handler: mux,
        ErrorLog: log.Default(),
    }

    users, err := st.DBH.GetUsers()
    if err != nil {
        fmt.Printf("get users: %v\n", err)
        return
    }
    fmt.Printf("%s:\n", dbh.Users)
    for i := range len(users) {
        fmt.Printf("user: %v\n", users[i])
    }
    offers, err := st.DBH.GetOffers()
    if err != nil {
        fmt.Printf("get offers: %v\n", err)
        return
    }
    fmt.Printf("%s:\n", dbh.Offers)
    for i := range len(offers) {
        fmt.Printf("offer: %v\n", offers[i])
    }


    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan

        shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
        defer shutdownRelease()

        if err := st.SRV.Shutdown(shutdownCtx); err != nil {
            log.Fatalf("HTTP shutdown error: %v", err)
        }
    }()

    if err := st.SRV.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        fmt.Printf("err: %v\n", err)
    }
    fmt.Println("Shutdown")
    st.DBH.DB.Close()
}
