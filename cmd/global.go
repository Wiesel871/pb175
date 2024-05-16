package cmd

import (
	_ "fmt"
	"net/http"
    "log"
    "context"
    "time"

	comp "wiesel/pb175/components"
    ut "wiesel/pb175/cmd/utility"
    st "wiesel/pb175/state"
    ct "wiesel/pb175/cmd/client"
    us "wiesel/pb175/cmd/user"
    of "wiesel/pb175/cmd/offers"

	_"github.com/a-h/templ"
)



func Home(st ut.GSP) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)
        user := ut.GetUser(st, id)
        comp.Page(comp.IndexBody(), user, comp.HomeN).Render(r.Context(), w)
    }
}

func Shutdown(st ut.GSP) ut.Response {
    return func(w http.ResponseWriter, r *http.Request) {
        id := ut.GetClientID(r)
        user := ut.GetUser(st, id)
        if !user.IsAdmin {
            comp.Page(comp.Forbidden(), user, comp.All).Render(r.Context(), w)
            return
        }

        shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 3*time.Second)
        defer shutdownRelease()

        if err := st.SRV.Shutdown(shutdownCtx); err != nil {
            log.Fatalf("HTTP shutdown error: %v", err)
        }
    }
}

func SetupUserHandler(mux *http.ServeMux, st *st.GlobalState) {
    mux.HandleFunc("/", Home(st))
    mux.HandleFunc("GET /home", Home(st))

    fs := http.FileServer(http.Dir("images"))
    mux.HandleFunc("GET /images/", http.StripPrefix("/images/", fs).ServeHTTP)

    fs = http.FileServer(http.Dir("scripts"))
    mux.HandleFunc("GET /scripts/", http.StripPrefix("/scripts/", fs).ServeHTTP)

    mux.HandleFunc("POST /logout", ct.LogOut)

    mux.HandleFunc("GET /signup", ct.SignUpForm(st))
    mux.HandleFunc("POST /signup", ct.SignUp(st))

    mux.HandleFunc("GET /login", ct.LogInForm(st))
    mux.HandleFunc("POST /login", ct.LogIn(st))

    mux.HandleFunc("/profile", ct.RedirectToUser(st))

    mux.HandleFunc("GET /profile/{id}", us.Profile(st))
    mux.HandleFunc("GET /users", us.GetAllUsers(st))

    mux.HandleFunc("POST /change_details" , us.ChangeDetails(st))

    mux.HandleFunc("GET /profile/{id}/offers", of.GetUserOffers(st))
    mux.HandleFunc("GET /profile/{id_owner}/offers/{id}", of.GetOffer(st))

    mux.HandleFunc("GET /offers", of.GetOffers(st))
    mux.HandleFunc("GET /offers/{by}/{sc}/", of.GetOffers(st))
    mux.HandleFunc("GET /offers/{by}/{sc}/{fil}", of.GetOffers(st))
    mux.HandleFunc("POST /offers/{by}/{sc}", of.FilterOffers(st))

    mux.HandleFunc("GET /add_offer", of.AddOffer(st))
    mux.HandleFunc("POST /add_offer", of.UploadOffer(st))

    mux.HandleFunc("POST /change_offer/{id}", of.ChangeOffer(st))

    mux.HandleFunc("DELETE /profile/{id}", us.DeleteUser(st))

    mux.HandleFunc("DELETE /profile/{owner}/offers/{id}", of.DeleteOffer(st))

    mux.HandleFunc("POST /promote/{id}", us.Mote(st, st.DBH.Promote))
    mux.HandleFunc("POST /demote/{id}", us.Mote(st, st.DBH.Demote))

    mux.HandleFunc("POST /shutdown", Shutdown(st))
}
