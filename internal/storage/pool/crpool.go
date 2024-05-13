package pool

import (
	"Server/internal/chatroom"
	"errors"
	"sync"
)

var ErrCRAlreadyExists = errors.New("cr with this id already exist")
var ErrNoCr = errors.New("cr with this id does not exist")

type Pool struct {
	mu           sync.RWMutex
	ChatRoomPool map[int]*chatroom.ChatRoom
}

func New() *Pool {
	return &Pool{
		ChatRoomPool: make(map[int]*chatroom.ChatRoom),
	}
}

func (p *Pool) GetChatRoom(crid int) (*chatroom.ChatRoom, error) {
	p.mu.RLock()
	cr, ok := p.ChatRoomPool[crid]
	p.mu.RUnlock()
	if !ok {
		return nil, ErrNoCr
	}

	return cr, nil
}

func (p *Pool) CreateChatRoom(crid int) error {
	_, ok := p.ChatRoomPool[crid]
	if ok {
		return ErrCRAlreadyExists
	}

	cr := chatroom.New()
	p.mu.Lock()
	p.ChatRoomPool[crid] = cr
	p.mu.Unlock()

	go func() {
		<-cr.GetDoneChan()
		p.RemChatRoom(crid)
	}()

	return nil
}

func (p *Pool) RemChatRoom(crid int) error {
	_, ok := p.ChatRoomPool[crid]
	if !ok {
		return ErrNoCr
	}

	p.mu.Lock()
	delete(p.ChatRoomPool, crid)
	p.mu.Unlock()
	return nil
}
