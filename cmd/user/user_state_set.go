package user

import (
	"fmt"
	_ "fmt"
	"net/http"
	_ "strconv"

	comp "wiesel/pb175/components"
    ut "wiesel/pb175/cmd/utility"


	_ "github.com/a-h/templ"
)


func ChangeDetails(st ut.GSP) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        client_id := ut.GetClientID(r)

        client, err := st.DBH.GetUserById(client_id)
        if err != nil {
            w.WriteHeader(403)
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
        }
        name := r.FormValue("name")
        details := r.FormValue("details")
        if err = st.DBH.AdjustUser(client, name, details); err != nil {
            w.WriteHeader(400)
            fmt.Printf("err: %v\n", err)
        }
        client.Name = name
        client.Details = details
        comp.ChangeDetails(client).Render(r.Context(), w)
    }
}
