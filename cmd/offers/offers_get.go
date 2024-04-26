package offers

import (
    "fmt"
    _ "fmt"
    "net/http"
    "strconv"

    ut "wiesel/pb175/cmd/utility"
    comp "wiesel/pb175/components"
)

func GetUserOffers(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
        user := ut.GetClientID(r)
        if user > -1 {
            http.SetCookie(w, ut.NewSession(user))
        }

        client := ut.GetUser(st, user)

        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            w.WriteHeader(404)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
            println("wrong id in path")
            return
        }
        owner, err := st.DBH.GetUserById(id)
        if err != nil {
            w.WriteHeader(404)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
            println("wrong user")
            return
        }
        offers, err := st.DBH.GetOffersByOwner(id)
        if err != nil {
            fmt.Printf("err: %v\n", err)
            w.WriteHeader(500)
            return
        }
        comp.Page(
            comp.Offers(offers, owner, client), 
            client, 
            comp.All,
        ).Render(r.Context(), w)
    }
}

func GetOffer(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
        user := ut.GetClientID(r)
        if user > -1 {
            http.SetCookie(w, ut.NewSession(user))
        }

        client := ut.GetUser(st, user)

        id_owner, err := strconv.Atoi(r.PathValue("id_owner"))
        if err != nil {
            w.WriteHeader(404)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
            return
        }

        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            w.WriteHeader(404)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
            return
        }

        owner, err := st.DBH.GetUserById(id_owner)
        if err != nil {
            w.WriteHeader(404)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
            return
        }


        offer, err := st.DBH.GetOfferById(id)
        if err != nil || offer.ID_owner != id_owner {
            w.WriteHeader(404)
            comp.Page(comp.NotFound(), client, comp.All).Render(r.Context(), w)
            return
        }
        comp.Page(comp.Offer(offer, owner, client, ""), client, comp.All).Render(r.Context(), w)
    }
}


func GetOffers(st ut.GSP) ut.Response {
    return func (w http.ResponseWriter, r *http.Request) {
        offers, _ := st.DBH.GetOffers()
        id := ut.GetClientID(r)
        if id > -1 {
            http.SetCookie(w, ut.NewSession(id))
        }
        client := ut.GetUser(st, id)

        comp.Page(
            comp.Offers(offers, st.Anonym, client), 
            client, 
            comp.OffersN,
        ).Render(r.Context(), w)
    }
}
