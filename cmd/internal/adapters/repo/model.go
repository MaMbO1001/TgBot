package repo

import (
	"github.com/google/uuid"
	"tgbot/cmd/internal/app"
)

type (
	user struct {
		ID       uuid.UUID `db:"id" json:"id"`
		Name     string    `db:"name" json:"name"`
		TgUserID int64     `db:"tg_user_id" json:"tg_user_id"`
	}
	folder struct {
		ID       uuid.UUID `db:"id" json:"id"`
		Name     string    `db:"name" json:"name"`
		TgUserID int64     `db:"tg_user_id" json:"tg_user_id"`
	}
	card struct {
		ID       uuid.UUID `db:"id" json:"id"`
		Name     string    `db:"name" json:"name"`
		FolderID uuid.UUID `db:"folder_id" json:"folder_id"`
		Text     string    `db:"text" json:"text"`
	}
)

func convert(u app.User) *user {
	return &user{
		ID:       u.ID,
		Name:     u.Name,
		TgUserID: u.TgUserID,
	}
}
func (u user) convert() *app.User {
	return &app.User{
		ID:       u.ID,
		Name:     u.Name,
		TgUserID: u.TgUserID,
	}
}

func convertFolder(f app.Folder) *folder {
	return &folder{
		ID:       f.ID,
		Name:     f.Name,
		TgUserID: f.TgUserID,
	}
}
func (f folder) convertFolder() *app.Folder {
	return &app.Folder{
		ID:       f.ID,
		Name:     f.Name,
		TgUserID: f.TgUserID,
	}
}

func convertCard(c app.Card) *card {
	return &card{
		ID:       c.ID,
		Name:     c.Name,
		FolderID: c.FolderID,
		Text:     c.Text,
	}
}
func (c card) convertCard() *app.Card {
	return &app.Card{
		ID:       c.ID,
		Name:     c.Name,
		FolderID: c.FolderID,
		Text:     c.Text,
	}
}
