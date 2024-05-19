package user

import (
    _ "fmt"
    "net/http"
    "strconv"
    "os"

    comp "wiesel/pb175/components"
    ut "wiesel/pb175/cmd/utility"

    _ "github.com/a-h/templ"
)

func DeleteUser(st ut.GSP) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        cl_id := ut.GetClientID(r)

        client := ut.GetUser(st, cl_id)

        target_str := r.PathValue("id")
        target, err := strconv.ParseInt(target_str, 10, 64)
        if err != nil {
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
            return
        }

        if (client.ID != target && !client.IsAdmin) || client.ID == 0 {
            comp.Page(comp.Forbidden(), client, comp.All).Render(r.Context(), w)
            return
        }

        err = st.DBH.DeleteUser(target)
        if err != nil {
            w.WriteHeader(500)
            return
        }

        os.RemoveAll("images/" + target_str)
        if client.IsAdmin && client.ID != target {
            http.Redirect(w, r, "/users", http.StatusFound)
            return
        }
        http.Redirect(w, r, "/logout", http.StatusFound)
    }
}
