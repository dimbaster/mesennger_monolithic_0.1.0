package getChatList

import (
	"Server/internal/models"
	"Server/internal/tokens"
	"encoding/json"
	"log/slog"
	"net/http"
)

type database interface {
	GetChats(userid int) ([]models.Chat, error)
}

type response struct {
	Chats []models.Chat `json:"chats"`
}

func New(s database, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := tokens.VerifyAndParse(r)
		if err != nil {
			log.Error("invalid token", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		chats, err := s.GetChats(user.ID)
		if err != nil {
			log.Error("failed to get chats from db", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		res := response{Chats: chats}
		data, err := json.Marshal(res)
		if err != nil {
			log.Error("failed to marshall json", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if err != nil {
			log.Error("failed to write to a conn", err)
		}
	}
}
