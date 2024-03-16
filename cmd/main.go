package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
    "log"
    "os"
    "path/filepath"
    "strings"
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

func ParseTemplates() *template.Template {
    templ := template.New("")
    err := filepath.Walk("./views", func(path string, info os.FileInfo, err error) error {
        if strings.Contains(path, ".html") {
            _, err = templ.ParseFiles(path)
            if err != nil {
                log.Println(err)
            }
        }
        return err
    })

    if err != nil {
        panic(err)
    }

    return templ
}

func newTemplate() *Templates {
    return &Templates{
        templates: ParseTemplates(),
    }
}

type Index struct {
    Count int
    Server *http.Server
    temp *Templates
    Left bool
    Theme string
    DB *sql.DB
    contacts_size uint
}

func (st *Index) Plus(w http.ResponseWriter, r *http.Request) {
    st.temp.Render(w, "count", st)
    st.Count++
}

func (st *Index) Index(w http.ResponseWriter, r *http.Request) {
    st.temp.Render(w, "index", st)
}

func (st *Index) postIndex(w http.ResponseWriter, r *http.Request) {
    st.temp.Render(w, "index_body", st)
}

func (st *Index) Buttons(w http.ResponseWriter, r *http.Request) {
    st.Left = !st.Left
    st.temp.Render(w, "buttons", st)
}

func (st *Index) SwapTheme(w http.ResponseWriter, r *http.Request) {
    if st.Theme == "dark" {
        st.Theme = "light"
    } else {
        st.Theme = "dark"
    }
    st.temp.Render(w, "theme", st)
}


func (st *Index) AddContact(w http.ResponseWriter, r *http.Request) {
    if err := insertContact(st.DB, newContact(st.contacts_size + 1, r.FormValue("name"), r.FormValue("email"))); err == nil {
        st.contacts_size += 1
    } else {
        fmt.Printf("insert err.Error(): %v on %s %s\n", err.Error(), r.FormValue("name"), r.FormValue("email"))
    }
    st.temp.Render(w, "contact_range", st)
}


func (st *Index) indexContacts(w http.ResponseWriter, r *http.Request) {
    st.temp.Render(w, "offers", st)
}

func (st *Index) bodyContacts(w http.ResponseWriter, r *http.Request) {
    st.temp.Render(w, "offers_body", st)
}

func main() {
    st := Index { 
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

    if err = db.QueryRow("SELECT COUNT(*) FROM Contacts").Scan(&st.contacts_size); err != nil {
        fmt.Printf("err.Error(): %v\n", err.Error())
        return 
    }
    fmt.Printf("st.contacts_size: %v\n", st.contacts_size)
    
    mux := http.NewServeMux()
    mux.HandleFunc("GET /", st.Index)
    mux.HandleFunc("POST /", st.postIndex)

    mux.HandleFunc("POST /count", st.Plus)
    mux.HandleFunc("POST /buttons", st.Buttons)
    mux.HandleFunc("POST /theme", st.SwapTheme)
    mux.HandleFunc("GET /offers", st.indexContacts)
    mux.HandleFunc("POST /offers", st.bodyContacts)
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
