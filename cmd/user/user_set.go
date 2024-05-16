package user

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	ut "wiesel/pb175/cmd/utility"
	comp "wiesel/pb175/components"

	_ "github.com/a-h/templ"
)


func ChangeDetails(st ut.GSP) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)

        client, err := st.DBH.GetUserById(id)
        if err != nil {
            w.WriteHeader(http.StatusForbidden)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
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

            path := "images/" + strconv.FormatInt(id, 10) + "/pfp.jpeg"
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
                w.WriteHeader(http.StatusNotFound)
                comp.Profile(client, client, err.Error()).Render(r.Context(), w)
                return
            }
            hasPfp = true
        }
        hasPfp = hasPfp || client.HasPFP

        client.HasPFP = hasPfp
        if err = st.DBH.AdjustUser(client, name, email, details, hasPfp); err != nil {
            w.WriteHeader(http.StatusNotFound)
            fmt.Printf("err: %v\n", err)
            return
        }

        http.Redirect(w, r, "/profile/" + strconv.FormatInt(id, 10), http.StatusFound)
    }
}

func Mote(st ut.GSP, mote func(int64) error) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)

        client, err := st.DBH.GetUserById(id)
        if err != nil || !client.IsAdmin {
            w.WriteHeader(http.StatusForbidden)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }
        
        target_id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
        if err != nil || target_id == 0 || target_id == id {
            w.WriteHeader(http.StatusForbidden)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }
        err = mote(target_id)
        if err != nil {
            fmt.Printf("err %v\n", err.Error())
            w.WriteHeader(http.StatusForbidden)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }

        target, err := st.DBH.GetUserById(target_id)
        if err != nil {
            w.WriteHeader(http.StatusForbidden)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }

        comp.PromDem(target).Render(r.Context(), w)
    }
}
