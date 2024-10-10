package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql/generated"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Store struct {
	db         *sql.DB
	config     pkg.Config
	tokenMaker pkg.Maker
}

func NewDB(config pkg.Config, maker pkg.Maker) *Store {
	return &Store{
		config:     config,
		tokenMaker: maker,
	}
}

func (s *Store) Open() error {
	if s.config.DB_DSN == "" {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "dsn is empty")
	}

	var err error

	s.db, err = sql.Open("mysql", s.config.DB_DSN)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to open database: %v", err)
	}

	err = s.db.Ping()
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to ping database: %v", err)
	}

	return s.migration()
}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}

func (s *Store) migration() error {
	if s.config.MIGRATION_PATH == "" {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "migrations directory is empty")
	}

	migration, err := migrate.New(s.config.MIGRATION_PATH, s.config.DB_DSN)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "Failed to load migration: %s", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "Failed to run migrate up: %s", err)
	}

	return nil
}

func (s *Store) execTx(ctx context.Context, fn func(q *generated.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "couldn't being transaction: %v", err)
	}

	q := generated.New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "tx err: %v, rb err: %v", err, rbErr)
		}

		return pkg.Errorf(pkg.INTERNAL_ERROR, "tx err: %v", err)
	}

	return tx.Commit()
}
