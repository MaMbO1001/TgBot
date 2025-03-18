package app

import (
	"context"
	"github.com/google/uuid"
)

func (a *App) CreateFolder(ctx context.Context, f Folder) (fold *Folder, err error) {
	fold, err = a.r.CheckExistsNameFolder(ctx, f)
	if err == nil {
		fold, err = a.r.CreateFolder(ctx, f)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return fold, nil
}

func (a *App) DeleteFolder(ctx context.Context, id uuid.UUID) (err error) {
	err = a.r.DeleteFolder(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (a *App) OpenFolder(ctx context.Context, id uuid.UUID) (card []Card, err error) {
	card, err = a.r.OpenFolder(ctx, id)
	if err != nil {
		return nil, err
	}
	return card, nil
}
func (a *App) GetAllFolders(ctx context.Context, id int64) (fold []Folder, err error) {
	fold, err = a.r.GetAllFolders(ctx, id)
	if err != nil {
		return nil, err
	}
	return fold, nil
}
func (a *App) CreateNameCard(ctx context.Context, c Card) (car *Card, err error) {
	car, err = a.r.CheckExistsCardName(ctx, c.Name, c.FolderID)
	if err == nil {
		car, err = a.r.CreateNameCard(ctx, c)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return car, nil
}
func (a *App) CreateTextCard(ctx context.Context, f Card) (err error) {
	err = a.r.CreateTextCard(ctx, f)
	if err != nil {
		return err
	}
	return nil
}
func (a *App) DeleteCard(ctx context.Context, name string) (err error) {
	err = a.r.DeleteCard(ctx, name)
	if err != nil {
		return err
	}
	return nil
}
func (a *App) CreateUser(ctx context.Context, us User) (err error) {
	err = a.r.CreateUser(ctx, us)
	if err != nil {
		return err
	}
	return nil
}
func (a *App) GetCard(ctx context.Context, id uuid.UUID) (card *Card, err error) {
	c, err := a.r.GetCard(ctx, id)
	if err != nil {
		return nil, err
	}
	return c, nil
}
