package app

import (
	"context"
	"github.com/google/uuid"
)

type Repo interface {
	CreateFolder(ctx context.Context, folder Folder) (*Folder, error)
	DeleteFolder(ctx context.Context, id uuid.UUID) error
	GetAllFolders(ctx context.Context, id int64) ([]Folder, error)
	OpenFolder(ctx context.Context, id uuid.UUID) ([]Card, error)
	CreateNameCard(ctx context.Context, card Card) (*Card, error)
	CreateTextCard(ctx context.Context, card Card) error
	DeleteCard(ctx context.Context, name string) error
	CreateUser(ctx context.Context, user User) error
	GetCard(ctx context.Context, id uuid.UUID) (*Card, error)
	CheckExistsCardName(ctx context.Context, name string, folderID uuid.UUID) (*Card, error)
	CheckExistsNameFolder(ctx context.Context, f Folder) (fold *Folder, err error)
}
