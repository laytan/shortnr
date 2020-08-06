package click

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/api/user"
	"github.com/laytan/shortnr/pkg/ratelimit"
	"github.com/laytan/shortnr/pkg/responder"
)

// SetRoutes sets the routes for interacting with clicks
func SetRoutes(r *mux.Router, store Storage) {
	clickR := r.PathPrefix("").Subrouter()
	clickR.Use(ratelimit.Middleware(ratelimit.GeneralRateLimit))
	clickR.Use(user.JwtAuthorization)

	clickR.HandleFunc("/{id}", all(store)).Methods(http.MethodGet)
}

func all(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := user.GetUser(r)
		if err != nil {
			responder.CastAndSend(err, w)
			return
		}

		// FIXME: QUICK AND VERY DIRTY
		linkID := mux.Vars(r)["id"]
		var linkOwnerID uint
		err = store.(MysqlStorage).Conn.QueryRow("SELECT user_id FROM links WHERE id = ?", linkID).Scan(&linkOwnerID)
		if err != nil {
			responder.Err{
				Code: http.StatusInternalServerError,
				Err:  err,
			}.Send(w)
			return
		}
		if linkOwnerID != user.ID {
			responder.Err{
				Code: http.StatusForbidden,
				Err:  errors.New("this is not a link you created"),
			}.Send(w)
			return
		}

		clicks := GetAll(linkID, store)

		responder.Res{
			Data: clicks,
		}.Send(w)
	}
}
