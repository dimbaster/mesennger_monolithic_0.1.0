package login

import (
	"Server/internal/models"
	"Server/internal/tokens"
	"encoding/json"
	"log/slog"
	"net/http"
)

type database interface {
	GetUser(login string) (models.User, error)
}

type request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type response struct {
	Token string `json:"token"`
}

func New(s database, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqData request
		err := json.NewDecoder(r.Body).Decode(&reqData)
		if err != nil {
			log.Error("failed to parse json from request body", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		user, err := s.GetUser(reqData.Login)
		if err != nil {
			log.Error("failed to get a user from db", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}

		if user.Password != reqData.Password {
			http.Error(w, "wrong login or password", http.StatusUnauthorized)
			return
		}

		tokenString, err := tokens.CreateToken(user.ID, user.Login, user.Password)
		if err != nil {
			log.Error("failed to create token", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		resData := response{
			Token: tokenString,
		}

		res, err := json.Marshal(resData)
		if err != nil {
			log.Error("failed to encode a json", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(res)
		if err != nil {
			log.Error("failed to send a token to the client", err)
		}
	}
}
