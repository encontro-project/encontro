package entity

type UserSummary struct {
	Servers []ServerInfo `json:"servers"`
}

type ServerInfo struct {
	Title           string         `json:"title"`
	VoiceChannels   []VoiceChannel `json:"voice_channels"`
	TextChats       []TextChat     `json:"text_chats"`
	AssociatedUsers []UserInfo     `json:"accociated_users"`
}

type VoiceChannel struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type TextChat struct {
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	Messages []Message `json:"messages"`
}

type UserInfo struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}
