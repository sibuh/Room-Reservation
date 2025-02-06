package initiator

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
)

func DoMigration(connString, filePath string) {
	m, err := migrate.New(fmt.Sprintf("file://%s", filePath),
		"cockroachdb://"+strings.Split(connString, "//")[1])
	if err != nil {
		log.Fatal("failed to create migration instance", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal("failed to do migration: ", err)
		}
	}
}
