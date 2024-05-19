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

	cmd "wiesel/pb175/cmd"
	_ "wiesel/pb175/components"
	data "wiesel/pb175/database"
	gst "wiesel/pb175/state"

	_ "github.com/a-h/templ"
)

func main() {
    dbh, err := data.InitDB("bazos.db")
    if err != nil {
        return
    }


    st := gst.GlobalState {
        DBH: dbh,
        Anonym: gst.GetAnonym(),
    }

    mux := http.NewServeMux()
    cmd.SetupUserHandler(mux, &st)
    st.SRV = &http.Server {
        Addr: ":8090",
        Handler: mux,
        ErrorLog: log.Default(),
    }


    id := int64(0)
    admin, _ := data.NewUser(id, "@name@", "@email@", "@password@")
    admin.IsAdmin = true
    st.DBH.InsertUser(admin)
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
