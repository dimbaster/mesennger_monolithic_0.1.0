package create

import (
	"Server/internal/tokens"
	"log/slog"
	"net/http"
)

type database interface {
	CreateChat(ownerid int) error
}

func New(s database, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := tokens.VerifyAndParse(r)
		if err != nil {
			log.Error("wrong token", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		err = s.CreateChat(user.ID)
		if err != nil {
			log.Error("failed to create a chat", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
