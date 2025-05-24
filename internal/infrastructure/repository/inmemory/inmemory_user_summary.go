package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"encontro/internal/domain/entity"
	"encontro/internal/domain/repository"
)

type InMemoryUserSummaryRepository struct {
	servers     map[int64]string // serverID -> title
	chats       []*ChatRecord
	serverUsers []*ServerUserRecord
	users       map[int64]*entity.UserInfo
	messages    map[int64][]*entity.Message // chatID -> messages
	mu          sync.RWMutex
}

type ChatRecord struct {
	ID       int64
	ServerID int64
	Title    string
	Type     string // "text" or "voice"
}

type ServerUserRecord struct {
	ServerID int64
	UserID   int64
}

func NewInMemoryUserSummaryRepository(cfg *InMemoryConfig) repository.UserSummaryRepository {
	r := &InMemoryUserSummaryRepository{
		servers:     make(map[int64]string),
		chats:       []*ChatRecord{},
		serverUsers: []*ServerUserRecord{},
		users:       make(map[int64]*entity.UserInfo),
		messages:    make(map[int64][]*entity.Message),
	}

	if cfg != nil && cfg.InitTestData {
		r.generateTestData()
	}
	return r
}

func (r *InMemoryUserSummaryRepository) FetchUserSummary(ctx context.Context, userID int64) (*entity.UserSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	summary := &entity.UserSummary{}

	for _, su := range r.serverUsers {
		if su.UserID != userID {
			continue
		}
		title := r.servers[su.ServerID]
		srv := entity.ServerInfo{Title: title}

		for _, ch := range r.chats {
			if ch.ServerID == su.ServerID && ch.Type == "voice" {
				srv.VoiceChannels = append(srv.VoiceChannels, entity.VoiceChannel{
					ID:    ch.ID,
					Title: ch.Title,
				})
			}
		}

		for _, ch := range r.chats {
			if ch.ServerID == su.ServerID && ch.Type == "text" {
				tc := entity.TextChat{ID: ch.ID, Title: ch.Title}
				msgs := r.messages[ch.ID]
				sort.Slice(msgs, func(i, j int) bool {
					return msgs[i].CreatedAt.After(msgs[j].CreatedAt)
				})
				limit := 50
				if len(msgs) < 50 {
					limit = len(msgs)
				}
				for i := 0; i < limit; i++ {
					tc.Messages = append(tc.Messages, *msgs[i])
				}
				srv.TextChats = append(srv.TextChats, tc)
			}
		}

		for _, su2 := range r.serverUsers {
			if su2.ServerID == su.ServerID {
				if u, ok := r.users[su2.UserID]; ok {
					srv.AssociatedUsers = append(srv.AssociatedUsers, *u)
				}
			}
		}

		summary.Servers = append(summary.Servers, srv)
	}

	return summary, nil
}

func (r *InMemoryUserSummaryRepository) AddUser(ctx context.Context, user *entity.UserInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserSummaryRepository) UpdateUser(ctx context.Context, user *entity.UserInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[user.ID]; !ok {
		return errors.New("user not found")
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserSummaryRepository) DeleteUser(ctx context.Context, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[userID]; !ok {
		return errors.New("user not found")
	}
	delete(r.users, userID)
	return nil
}

func (r *InMemoryUserSummaryRepository) AddServer(ctx context.Context, serverID int64, title string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.servers[serverID] = title
	return nil
}

func (r *InMemoryUserSummaryRepository) AddUserToServer(ctx context.Context, userID, serverID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.serverUsers = append(r.serverUsers, &ServerUserRecord{ServerID: serverID, UserID: userID})
	return nil
}

func (r *InMemoryUserSummaryRepository) AddChat(ctx context.Context, serverID int64, title, chatType string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var nextID int64 = 1
	for _, c := range r.chats {
		if c.ID >= nextID {
			nextID = c.ID + 1
		}
	}
	r.chats = append(r.chats, &ChatRecord{
		ID:       nextID,
		ServerID: serverID,
		Title:    title,
		Type:     chatType,
	})
	return nil
}

func (r *InMemoryUserSummaryRepository) AddMessage(ctx context.Context, chatID int64, content string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	msg := &entity.Message{
		ID:        fmt.Sprintf("msg-%d-%d", chatID, time.Now().UnixNano()),
		Content:   content,
		RoomID:    fmt.Sprintf("%d", chatID),
		SenderID:  "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	r.messages[chatID] = append(r.messages[chatID], msg)
	return nil
}

func (r *InMemoryUserSummaryRepository) UpdateUserAvatar(ctx context.Context, userID int64, newURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[userID]
	if !ok {
		return errors.New("user not found")
	}
	u.AvatarURL = newURL
	return nil
}

func (r *InMemoryUserSummaryRepository) UpdateServerTitle(ctx context.Context, serverID int64, newTitle string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.servers[serverID]; !ok {
		return errors.New("server not found")
	}
	r.servers[serverID] = newTitle
	return nil
}

func (r *InMemoryUserSummaryRepository) generateTestData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.servers[1] = "Test Server 1"
	r.servers[2] = "Test Server 2"

	r.users[100] = &entity.UserInfo{ID: 100, Username: "alice", AvatarURL: "avatar1.png"}
	r.users[101] = &entity.UserInfo{ID: 101, Username: "bob", AvatarURL: "avatar2.png"}

	r.serverUsers = append(r.serverUsers,
		&ServerUserRecord{ServerID: 1, UserID: 100},
		&ServerUserRecord{ServerID: 1, UserID: 101},
		&ServerUserRecord{ServerID: 2, UserID: 100},
	)

	r.chats = append(r.chats,
		&ChatRecord{ID: 1, ServerID: 1, Title: "general", Type: "text"},
		&ChatRecord{ID: 2, ServerID: 1, Title: "voice 1", Type: "voice"},
		&ChatRecord{ID: 4, ServerID: 1, Title: "voice 2", Type: "voice"}, // Added voice channel for server 1
		&ChatRecord{ID: 3, ServerID: 2, Title: "random", Type: "text"},
		&ChatRecord{ID: 5, ServerID: 2, Title: "voice A", Type: "voice"}, // Added voice channel for server 2
	)

	now := time.Now()
	r.messages[1] = []*entity.Message{
		{
			ID:        "msg-1",
			Content:   "Hello World",
			RoomID:    "1",
			SenderID:  "100",
			CreatedAt: now.Add(-1 * time.Minute),
			UpdatedAt: now.Add(-1 * time.Minute),
		},
		{
			ID:        "msg-2",
			Content:   "How are you?",
			RoomID:    "1",
			SenderID:  "101",
			CreatedAt: now.Add(-30 * time.Second),
			UpdatedAt: now.Add(-30 * time.Second),
		},
	}

	r.messages[3] = []*entity.Message{
		{
			ID:        "msg-3",
			Content:   "Welcome to server 2",
			RoomID:    "3",
			SenderID:  "100",
			CreatedAt: now.Add(-10 * time.Minute),
			UpdatedAt: now.Add(-10 * time.Minute),
		},
		{
			ID:        "msg-4",
			Content:   "Second message",
			RoomID:    "3",
			SenderID:  "100",
			CreatedAt: now.Add(-9 * time.Minute),
			UpdatedAt: now.Add(-9 * time.Minute),
		},
	}
}
