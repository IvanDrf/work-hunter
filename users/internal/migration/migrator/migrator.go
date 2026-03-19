package migrator

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func MigrateUp(databaseURL string) error {
	m, err := migrate.New("file://internal/migration/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create mugrator: %w", err)
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	return nil
}

func MigrateDown(databaseURL string) error {
	m, err := migrate.New("file://internal/migration/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create mugrator: %w", err)
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to migrate down: %w", err)
	}

	return nil
}
