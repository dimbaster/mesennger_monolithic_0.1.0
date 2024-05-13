package database

import (
	"Server/internal/models"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	DB *sql.DB
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

func openDB(driverName, storagePath string) *sql.DB {
	db, err := sql.Open(driverName, storagePath)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func New() *Storage {
	return &Storage{
		DB: openDB("pgx", "postgres://postgres:root@localhost:5432/messenger"),
	}
}

func (s *Storage) GetUser(login string) (models.User, error) {
	var res models.User
	row := s.DB.QueryRow("SELECT * FROM users WHERE login = $1", login)
	err := row.Scan(&res.ID, &res.Login, &res.Password)
	if err != nil {
		return models.User{}, err
	}

	return res, nil
}

func (s *Storage) CreateUser(login string, password string) error {
	res, err := s.DB.Exec(`INSERT INTO users(login, password) VALUES($1, $2) ON CONFLICT(login, password) DO NOTHING;`, login, password)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return ErrUserAlreadyExists
	}

	return nil

}

func (s *Storage) DeleteUser() {
	//TODO: implement
}

func (s *Storage) CreateChat(ownerid int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO chats(owner) values ($1)", ownerid)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	row := tx.QueryRow("SELECT id FROM chats WHERE owner = $1", ownerid)
	var chatid int64
	row.Scan(&chatid)

	_, err = tx.Exec("INSERT INTO chats_users(userid, chatid) VALUES ($1, $2)", ownerid, chatid)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err

	}

	return nil
}

func (s *Storage) getUserChatIds(userid int) ([]int, error) {
	rows, err := s.DB.Query("SELECT chatid FROM chats_users WHERE userid = $1", userid)
	if err != nil {
		return nil, err
	}

	var res []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		res = append(res, int(id))
	}

	return res, nil
}

func (s *Storage) GetChats(userid int) ([]models.Chat, error) {
	chatids, err := s.getUserChatIds(userid)
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query("SELECT * FROM chats WHERE id = ANY($1)", chatids)
	if err != nil {
		return nil, err
	}

	var res []models.Chat
	for rows.Next() {
		var c models.Chat
		err = rows.Scan(&c.ID, &c.Owner)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}

	return res, nil
}

func (s *Storage) GetChat(chatid int) (models.Chat, error) {
	row := s.DB.QueryRow("SELECT * FROM chats WHERE id = $1", chatid)
	var chat models.Chat
	err := row.Scan(&chat.ID, &chat.Owner)
	if err != nil {
		return models.Chat{}, err
	}

	return chat, nil
}

func (s *Storage) AddUserToChat(userid int, chatid int) error {
	_, err := s.DB.Exec("INSERT INTO chats_users(userid, chatid) VALUES ($1, $2) ON CONFLICT (userid, chatid) DO NOTHING", userid, chatid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteChat() {
	//TODO: implement
}

func (s *Storage) SaveMessage(chatid int, senderid int, text string) {
	s.DB.Exec("INSERT INTO messages (userid, chatid, message) VALUES ($1, $2, $3)", senderid, chatid, text)

}

func (s *Storage) GetMessages(chatid int) ([]models.Message, error) {
	rows, err := s.DB.Query("SELECT * FROM messages WHERE chatid = $1", chatid)
	if err != nil {
		return nil, err
	}

	var res []models.Message
	for rows.Next() {
		var msg models.Message
		err = rows.Scan(&msg.ID, &msg.UserID, &msg.ChatID, &msg.Text)
		if err != nil {
			return nil, err
		}
		res = append(res, msg)
	}

	return res, nil
}

func (s *Storage) DeleteMessage() {
	//TODO: implement
}
