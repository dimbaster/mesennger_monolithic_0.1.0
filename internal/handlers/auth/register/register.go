package register

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type database interface {
	CreateUser(login string, password string) error
}

type request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func New(s database, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("decoding request body into json went wrong", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		err = s.CreateUser(req.Login, req.Password)
		if err != nil {
			log.Error("failed to create a user", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
