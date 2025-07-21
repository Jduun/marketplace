package migrations

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func Migrate(DBURL string) {
	m, err := migrate.New(
		"file://migrations",
		DBURL,
	)
	if err != nil {
		log.Fatalf("Migration error: %s", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migration error: %s", err)
	}
}
