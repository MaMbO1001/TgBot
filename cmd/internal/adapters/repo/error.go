package repo

import (
	"database/sql"
	"errors"
	"tgbot/cmd/internal/app"
)

func convertErr(err error) error {

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return app.ErrNotFound
	default:
		return err
	}
}
