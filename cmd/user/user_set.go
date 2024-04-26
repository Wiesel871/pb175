package user

import (
	"fmt"
	"net/http"
	"strconv"
    "os"
    "io"

	comp "wiesel/pb175/components"
    ut "wiesel/pb175/cmd/utility"


	_ "github.com/a-h/templ"
)


func ChangeDetails(st ut.GSP) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)

        client, err := st.DBH.GetUserById(id)
        if err != nil {
            w.WriteHeader(403)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
        }
        name := r.FormValue("name")
        email := r.FormValue("email")
        details := r.FormValue("details")

        client.Name = name
        client.Details = details
        

        err = r.ParseMultipartForm(10 << 20) 
        if err != nil {
            return
        }
        hasPfp := false

        file, _, err := r.FormFile("photo")
        if err != nil && err != http.ErrMissingFile {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        if err == nil {
            defer file.Close()
            println("got a file")

            path := "images/" + strconv.Itoa(id) + "/pfp.jpeg"
            f, err := os.OpenFile(
                path,
                os.O_WRONLY | os.O_CREATE,
                0666)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                fmt.Printf("err: %v\n", err)
                return
            }
            defer f.Close()

            _, err = io.Copy(f, file)
            if err != nil {
                w.WriteHeader(400)
                comp.Profile(client, client, err.Error()).Render(r.Context(), w)
                return
            }
            hasPfp = true
        }
        hasPfp = hasPfp || client.HasPFP

        client.HasPFP = hasPfp
        if err = st.DBH.AdjustUser(client, name, email, details, hasPfp); err != nil {
            w.WriteHeader(400)
            fmt.Printf("err: %v\n", err)
        }

        w.Header().Set("HX-Refresh", "true")
        comp.Page(comp.Profile(client, client, ""), client, comp.All).Render(r.Context(), w)
    }
}
