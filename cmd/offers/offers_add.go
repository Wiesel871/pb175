package offers

import (
	"fmt"
    "io"
	"net/http"
	"os"
	"strconv"

	ut "wiesel/pb175/cmd/utility"
	comp "wiesel/pb175/components"
	db "wiesel/pb175/database"
)

func AddOffer(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
        println("got add offer get")
        id := ut.GetClientID(r)
        client := ut.GetUser(st, id)
        if client.ID < 0 {
            w.WriteHeader(403)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }
        comp.Page(comp.NewOffer(client), client, comp.All).Render(r.Context(), w)
    }
}

func UploadOffer(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
        println("got upload post")
        id := ut.GetClientID(r)
        client := ut.GetUser(st, id)
        if id == -1 {
            w.WriteHeader(403)
            println("invalid id")
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }

        name := r.FormValue("name")
        desc := r.FormValue("description")
        fmt.Printf("name: %v\n", name)
        fmt.Printf("desc: %v\n", desc)

        of_id, err := db.SmallestMissingID(st.DBH.DB, st.DBH.Offers)

        err = r.ParseMultipartForm(10 << 20) 
        if err != nil {
            return
        }

        file, _, err := r.FormFile("photo")
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        defer file.Close()

        path := "images/" + strconv.Itoa(id) + "/" + strconv.Itoa(of_id) + ".jpeg"
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
            http.Error(w, err.Error(), http.StatusInternalServerError)
            fmt.Printf("err: %v\n", err)
            return
        }
        if err = st.DBH.InsertOffer(db.NewOffer(of_id, id, name, desc)); err != nil { 
            fmt.Printf("uo err: %v\n", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/profile", http.StatusMovedPermanently)
    }
}
