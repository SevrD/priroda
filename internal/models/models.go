package models

// 1 - запрос текста объявления
// 2 - запрос фото
// 3 - запрос номера удаляемого объявления.
type StatusCode int64

type AnnouncementInfo struct {
	Text     string
	FileID   string
	ChatID   int64
	Id       int64
	TgID     int64
	PublicID int64
}

type Button struct {
	Text string
	Name string
}
