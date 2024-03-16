package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
    "log"
	//"github.com/a-h/templ"
)

type Templates struct {
    templates *template.Template;
}

func (t* Templates) Render(w io.Writer, name string, data interface{}) error {
    err := t.templates.ExecuteTemplate(w, name, data)
    if err != nil {
        fmt.Printf("render err.Error(): %v\n", err.Error())
    }
    return err
}

func newTemplate() *Templates {
    return &Templates{
        templates: template.Must(template.ParseGlob("views/*.html")),
    }
}

type State struct {
    Count int
    Server *http.Server
    temp *Templates
    Left bool
    Theme string
    DB *sql.DB
    contacts_size uint
}

func (st *State) Plus(w http.ResponseWriter, r *http.Request) {
    st.temp.Render(w, "stount", st)
    st.Count++
}

func (st *State) Index(w http.ResponseWriter, r *http.Request) {
    st.temp.Render(w, "index", st)
}

func (st *State) Buttons(w http.ResponseWriter, r *http.Request) {
    st.Left = !st.Left
    st.temp.Render(w, "buttons", st)
}

func (st *State) SwapTheme(w http.ResponseWriter, r *http.Request) {
    if st.Theme == "dark" {
        st.Theme = "light"
    } else {
        st.Theme = "dark"
    }
    st.temp.Render(w, "theme", st)
}


func (st *State) AddContact(w http.ResponseWriter, r *http.Request) {
    if err := insertContact(st.DB, newContact(st.contacts_size + 1, r.FormValue("name"), r.FormValue("email"))); err != nil {
        fmt.Printf("insert err.Error(): %v on %s %s\n", err.Error(), r.FormValue("name"), r.FormValue("email"))
        return
    }
    st.contacts_size++
    st.temp.Render(w, "contacts", st)
}

func main() {
    st := State { 
        Count: 0,
        temp: newTemplate(),
        Left: true,
        Theme: "dark",
    }
    db, err := initDB()
    if err != nil {
        fmt.Println("init Error")
        return
    }
    st.DB = db

    var size *sql.Rows
    if size, err = db.Query("SELECT COUNT(*) FROM Contacts"); err != nil {
        fmt.Printf("err.Error(): %v\n", err.Error())
        return 
    }
    if !size.Next() {
        return
    }
    if err := size.Scan(&st.contacts_size); err != nil {
        fmt.Printf("err.Error(): %v\n", err.Error())
        return 
    }
    fmt.Printf("st.contacts_size: %v\n", st.contacts_size)
    if st.contacts_size == 0 {
        insertContact(st.DB, newContact(0, "Joe", "Mamam"))
    }
    
    fmt.Println(st.getContacts())

    mux := http.NewServeMux()
    mux.HandleFunc("GET /", st.Index)
    mux.HandleFunc("POST /count", st.Plus)
    mux.HandleFunc("POST /buttons", st.Buttons)
    mux.HandleFunc("POST /theme", st.SwapTheme)
    mux.HandleFunc("POST /contacts", st.AddContact)
    srv := &http.Server {
        Addr: ":8090",
        Handler: mux,
        ErrorLog: log.Default(),
    }
    st.Server = srv
    
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Println("Err")
		}
	}()

    fmt.Println("Started")
    sigChan := make(chan uint)

    http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Received shutdown request")
        sigChan <- 0
        if err := srv.Shutdown(context.Background()); err != nil {
            fmt.Printf("Error shutting down server: %s\n", err)
        } else {
            fmt.Println("Server has been shut down")
        }
    })


    go func() {
        if err := http.ListenAndServe(":9090", nil); err != nil {
            fmt.Printf("Error starting shutdown server: %s\n", err)
        }
    }()
    <-sigChan
}
