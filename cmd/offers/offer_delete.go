package offers

import (
	"net/http"
	"os"
	"strconv"

	ut "wiesel/pb175/cmd/utility"
	comp "wiesel/pb175/components"
)


func DeleteOffer(st ut.GSP) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        client_id := ut.GetClientID(r)
        client := ut.GetUser(st, client_id)

        owner_id_str := r.PathValue("owner")
        id_str := r.PathValue("id")
        owner_id, err := strconv.ParseInt(owner_id_str, 10, 64)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
            return
        }

        if owner_id != client_id && !client.IsAdmin {
            w.WriteHeader(http.StatusForbidden)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }

        id, err := strconv.ParseInt(id_str, 10, 64)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
            return
        }

        err = st.DBH.DeleteOffer(id)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)            
            return
        }

        os.Remove("/images/" + owner_id_str + "/" + id_str + ".jpeg")
    }
}
