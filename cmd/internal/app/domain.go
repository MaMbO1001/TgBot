package app

import "github.com/google/uuid"

type (
	User struct {
		ID       uuid.UUID
		Name     string
		TgUserID int64
	}
	Folder struct {
		ID       uuid.UUID
		Name     string
		TgUserID int64
	}
	Card struct {
		ID       uuid.UUID
		Name     string
		FolderID uuid.UUID
		Text     string
	}
)
