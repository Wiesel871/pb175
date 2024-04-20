package offers

import (
    _ "fmt"
    "io"
    "net/http"
    "os"
    "strconv"

    comp "wiesel/pb175/components"
    ut "wiesel/pb175/cmd/utility"
    db "wiesel/pb175/database"
)

func AddOffer(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
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
        id := ut.GetClientID(r)
        client := ut.GetUser(st, id)
        if id == -1 {
            w.WriteHeader(403)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }

        name := r.FormValue("name")
        desc := r.FormValue("description")
        of_id, err := db.SmallestMissingID(st.DBH.DB, st.DBH.Offers)
        if err = st.DBH.InsertOffer(db.NewOffer(of_id, id, name, desc)); err != nil { 

        }


        err = r.ParseMultipartForm(10 << 20) // 10MB max
        if err != nil {
            return
        }

        file, _, err := r.FormFile("photo")
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        defer file.Close()

        f, err := os.OpenFile(
            "images/" + strconv.Itoa(id) + "/offers/" + strconv.Itoa(of_id) + ".jpg",
            os.O_WRONLY | os.O_CREATE,
            0666)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer f.Close()

        _, err = io.Copy(f, file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/profile", http.StatusMovedPermanently)
    }
}
