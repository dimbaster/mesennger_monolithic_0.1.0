package goToChatRoom

import (
	"Server/internal/chatroom"
	"Server/internal/models"
	"Server/internal/tokens"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type database interface {
	GetChats(userid int) ([]models.Chat, error)
}

type pool interface {
	GetChatRoom(crid int) (*chatroom.ChatRoom, error)
	CreateChatRoom(crid int) error
}

var upgrader = websocket.Upgrader{}

func New(s database, p pool, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := tokens.VerifyAndParse(r)
		if err != nil {
			log.Error("invalid token", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		chatid, err := strconv.Atoi(r.PathValue("chatid"))
		if err != nil {
			log.Error("failed to convert string to int(wtf?)", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		chats, err := s.GetChats(user.ID)
		if err != nil {
			log.Error("failed to get user chats", err)
			http.Error(w, "There is no such chat", http.StatusBadRequest)
			return
		}

		for i := 0; i < len(chats); i++ {
			if chats[i].ID == chatid {
				break
			}
			if i == len(chats)-1 {
				http.Error(w, "you are NOT allowed here", http.StatusUnauthorized)
				return
			}
		}

		cr, err := p.GetChatRoom(chatid)
		if err != nil {
			err = p.CreateChatRoom(chatid)
			if err != nil {
				http.Error(w, "something went wrong", http.StatusInternalServerError)
				return
			}
			cr, _ = p.GetChatRoom(chatid)
		}

		w.WriteHeader(http.StatusSwitchingProtocols)
		succesHeaders := http.Header{}
		succesHeaders.Add("Connection", "Upgrade")
		succesHeaders.Add("Upgrade", "websocket")

		conn, err := upgrader.Upgrade(w, r, succesHeaders)
		if err != nil {
			log.Error("failed to upgrade a user", err)
			http.Error(w, "something went wrong", http.StatusGone)
			return
		}

		cr.AddClient(user.ID, conn)
	}
}
