package offers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	ut "wiesel/pb175/cmd/utility"
	comp "wiesel/pb175/components"
	db "wiesel/pb175/database"
)

func AddOffer(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)
        client := ut.GetUser(st, id)
        if client.ID < 0 {
            w.WriteHeader(http.StatusForbidden)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }
        comp.Page(comp.NewOffer(client, ""), client, comp.All).Render(r.Context(), w)
    }
}

func UploadOffer(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)
        client := ut.GetUser(st, id)
        if id == -1 {
            w.WriteHeader(http.StatusForbidden)
            println("invalid id")
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }

        name := r.FormValue("name")
        desc := r.FormValue("description")

        of_id := (id << 32) + time.Now().Unix()
        offer := db.NewOffer(of_id, id, name, desc)

        err := r.ParseMultipartForm(10 << 20) 
        if err != nil {
            return
        }

        file, _, err := r.FormFile("photo")
        if err != nil {
            w.WriteHeader(http.StatusUnprocessableEntity)
            comp.Offer(offer, client, client, "photo missing").Render(r.Context(), w)
            return
        }
        defer file.Close()

        path := "images/" + strconv.FormatInt(id, 10) + "/" + strconv.FormatInt(of_id, 10) + ".jpeg"
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
        http.Redirect(w, r, "/profile/" + strconv.FormatInt(id, 10) + "/offers", http.StatusFound)
    }
}
