package addToChat

import (
	"Server/internal/models"
	"Server/internal/tokens"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type database interface {
	GetChat(chatid int) (models.Chat, error)
	AddUserToChat(chatid int, userid int) error
}

type request struct {
	ID int `json:"ID"`
}

func New(s database, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestOwner, err := tokens.VerifyAndParse(r)
		if err != nil {
			log.Error("invalid token", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		chatid, err := strconv.Atoi(r.PathValue("chatid"))
		if err != nil {
			http.Error(w, "chat id is not an int", http.StatusBadRequest)
			return
		}

		chat, err := s.GetChat(chatid)
		if err != nil {
			log.Error("failed to get chat")
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}
		if chat.Owner != requestOwner.ID {
			http.Error(w, "You are not the owner, you can`t add users", http.StatusUnauthorized)
			return
		}

		var req request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode json", err)
			http.Error(w, "something went wrong", http.StatusBadRequest)
			return
		}

		err = s.AddUserToChat(chatid, req.ID)
		if err != nil {
			log.Error("failed to add user to a chat", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
