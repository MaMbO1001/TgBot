package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
	"github.com/sipki-tech/database"
	"github.com/sipki-tech/database/connectors"
	"github.com/sipki-tech/database/migrations"
	"tgbot/cmd/internal/app"
)

type Repo struct {
	sql *database.SQL
}
type (
	// Config provide connection info for database.
	Config struct {
		Postgres   connectors.Raw
		MigrateDir string
		Driver     string
	}
)

func New(ctx context.Context, cfg Config) (*Repo, error) {
	const subsystem = "repo"

	// List of app.Err… returned by Repo methods.
	returnErrs := []error{
		app.ErrNotFound,
		app.ErrTelegramExist,
	}

	migrates, err := migrations.Parse(cfg.MigrateDir)
	if err != nil {
		return nil, fmt.Errorf("migrations.Parse: %w", err)
	}

	err = migrations.Run(ctx, cfg.Driver, &cfg.Postgres, migrations.Up, migrates)
	if err != nil {
		return nil, fmt.Errorf("migrations.Run: %w", err)
	}

	conn, err := database.NewSQL(ctx, cfg.Driver, database.SQLConfig{
		ReturnErrs: returnErrs,
	}, &cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("librepo.NewCockroach: %w", err)
	}

	return &Repo{
		sql: conn,
	}, nil
}

func (r *Repo) CheckExistsNameFolder(ctx context.Context, f app.Folder) (fold *app.Folder, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from folder where name = $1 and tg_user_id = $2`
		var folderr folder
		err = db.GetContext(ctx, &folderr, query, f.Name, f.TgUserID)
		if err == nil {
			return fmt.Errorf("db.Getconext: %w", convertErr(err))
			fold = folderr.convertFolder()
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fold, nil
}

func (r *Repo) CreateFolder(ctx context.Context, f app.Folder) (fold *app.Folder, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		newFolder := convertFolder(f)
		const query = `
INSERT INTO
folder
(name, tg_user_id)
values 
    ($1, $2)
    returning *`
		var res folder
		err := db.GetContext(ctx, &res, query, newFolder.Name, newFolder.TgUserID)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}
		fold = res.convertFolder()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fold, nil
}
func (r *Repo) DeleteFolder(ctx context.Context, id uuid.UUID) (err error) {
	return r.sql.NoTx(func(db *sqlx.DB) error {
		const qury = `delete
from card
where folder_id = $1 returning *`
		err = db.GetContext(ctx, &card{}, qury, id)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}
		const query = `
delete 
from folder
where id = $1 returning *`
		err := db.GetContext(ctx, &folder{}, query, id)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}
		return nil
	})
}
func (r *Repo) GetAllFolders(ctx context.Context, id int64) (f []app.Folder, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from folder where tg_user_id = $1`
		res := make([]folder, 0)
		err := db.SelectContext(ctx, &res, query, id)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}
		f = make([]app.Folder, len(res))
		for i := range res {
			f[i] = *res[i].convertFolder()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (r *Repo) OpenFolder(ctx context.Context, id uuid.UUID) (c []app.Card, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from card where folder_id = $1`
		res := make([]card, 0)
		err := db.SelectContext(ctx, &res, query, id)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}
		c = make([]app.Card, len(res))
		for i := range res {
			c[i] = *res[i].convertCard()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *Repo) CheckExistsCardName(ctx context.Context, name string, folderID uuid.UUID) (c *app.Card, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from card where name = $1 and folder_id = $2`
		var res card
		err = db.GetContext(ctx, &res, query, name, folderID)
		if err == nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}
		c = res.convertCard()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *Repo) CreateNameCard(ctx context.Context, c app.Card) (car *app.Card, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `insert into card (name, folder_id) values ($1, $2) returning *`
		var res card
		err = db.GetContext(ctx, &res, query, c.Name, c.FolderID)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}
		car = res.convertCard()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return car, nil
}
func (r *Repo) CreateTextCard(ctx context.Context, c app.Card) (err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `UPDATE card set text = $1 where id = $2 returning *`
		var res card
		err = db.GetContext(ctx, &res, query, c.Text, c.ID)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) DeleteCard(ctx context.Context, name string) (err error) {
	return r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `delete from card where name = $1 returning *`
		var res card
		err = db.GetContext(ctx, &res, query, name)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}
		return nil
	})
}

func (r *Repo) CreateUser(ctx context.Context, u app.User) (err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `insert into users_table (name, tg_user_id) values ($1, $2) returning *`
		us := convert(u)
		var res user
		err = db.GetContext(ctx, &res, query, us.Name, us.TgUserID)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetCard(ctx context.Context, id uuid.UUID) (c *app.Card, err error) {
	err = r.sql.NoTx(func(db *sqlx.DB) error {
		const query = `select * from card where id = $1`
		var res card
		err = db.GetContext(ctx, &res, query, id)
		if err != nil {
			return fmt.Errorf("db.GetContext: %w", convertErr(err))
		}
		c = res.convertCard()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

//func (r *Repo) RandomCard(ctx context.Context, id uuid.UUID) ([]app.Card, error) {
//}
