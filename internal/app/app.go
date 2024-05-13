package app

import (
	"Server/internal/handlers/auth/login"
	"Server/internal/handlers/auth/register"
	"Server/internal/handlers/chats/addToChat"
	"Server/internal/handlers/chats/create"
	"Server/internal/handlers/chats/getChatList"
	"Server/internal/handlers/chats/goToChatRoom"
	"Server/internal/storage/database"
	"Server/internal/storage/pool"
	"log/slog"
	"net/http"
	"os"
)

type App struct {
	log          *slog.Logger
	storage      *database.Storage
	chatRoomPool *pool.Pool
	server       http.Server
}

func New() *App {
	app := App{
		log:          slog.New(slog.NewTextHandler(os.Stdout, nil)),
		storage:      database.New(),
		chatRoomPool: pool.New(),
		server: http.Server{
			Addr:    "localhost:8080",
			Handler: nil,
		},
	}

	app.setupRoutes()
	return &app
}

func (a *App) ParseCfg() {
	//TODO: implement
}

func (a *App) Run() {
	a.log.Info("starting server")
	err := a.server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (a *App) setupRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", login.New(a.storage, a.log))
	mux.HandleFunc("POST /register", register.New(a.storage, a.log))
	mux.HandleFunc("GET /chats", getChatList.New(a.storage, a.log))
	mux.HandleFunc("POST /chats/create", create.New(a.storage, a.log))
	mux.HandleFunc("GET /chats/{chatid}", goToChatRoom.New(a.storage, a.chatRoomPool, a.log))
	mux.HandleFunc("PATCH /chats/{chatid}/addUser", addToChat.New(a.storage, a.log))
	//mux.HandleFunc("DELETE /chats/{chatid}")
	//mux.HandleFunc("PATCH /chats/{chatid}/removeUser")
	//mux.HandleFunc("GET /chats/{chatid}/messages")
	//mux.HandleFunc("DELETE /chats/{chatid}/message")
	//mux.HandleFunc("POST /chats/{chatid}/message")
	//mux.HandleFunc("PUT /chats/{chatid}/message")
	a.server.Handler = mux
}
