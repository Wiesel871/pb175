package offers

import (
	"fmt"
    "io"
	"net/http"
	"os"
	"strconv"

	ut "wiesel/pb175/cmd/utility"
	comp "wiesel/pb175/components"
)


func ChangeOffer(st ut.GSP) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)

        client, err := st.DBH.GetUserById(id)
        if err != nil {
            w.WriteHeader(http.StatusForbidden)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
        }
        of_id_str := r.PathValue("id")
        of_id, err := strconv.ParseInt(of_id_str, 10, 64)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
        }

        name := r.FormValue("name")
        desc := r.FormValue("desc")

        offer, err := st.DBH.GetOfferById(of_id)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
        }

        err = r.ParseMultipartForm(10 << 20) 
        if err != nil {
            return
        }

        file, _, err := r.FormFile("photo")
        if err != nil && err != http.ErrMissingFile {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        if err == nil {
            defer file.Close()

            path := "images/" + strconv.FormatInt(id, 10) + "/" + of_id_str + ".jpeg"
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
        }
        if err = st.DBH.AdjustOffer(offer, name, desc); err != nil {
            w.WriteHeader(400)
            return
        }

        offer.Description = desc
        offer.Name = name
        
        w.Header().Set("HX-Refresh", "true")
        comp.Page(comp.Offer(offer, client, client, ""), client, comp.All).Render(r.Context(), w)
    }
}
