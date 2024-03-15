package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
    "io"

	//"github.com/a-h/templ"
)

type Templates struct {
    templates *template.Template;
}

func (t* Templates) Render(w io.Writer, name string, data interface{}) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
    return &Templates{
        templates: template.Must(template.ParseGlob("views/*.html")),
    }
}

type State struct {
    Count int
    temp *Templates
}

func (c *State) Plus(w http.ResponseWriter, r *http.Request) {
    c.temp.Render(w, "count", c)
    c.Count++
}

func (c *State) Index(w http.ResponseWriter, r *http.Request) {
    c.temp.Render(w, "index", c)
}

func main() {
    st := State { 
        Count: 0,
        temp: newTemplate(),
    }

    mux := http.NewServeMux()
    mux.HandleFunc("GET /", st.Index)
    mux.HandleFunc("POST /count", st.Plus)
    srv := &http.Server {
        Addr: ":8090",
        Handler: mux,
    }
    
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
