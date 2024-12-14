package middleware

import (
	"context"
	"errors"
	"os"
	"reflect"
	"reservation/internal/storage/db"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type MQ struct {
	db.Querier
	users []db.User
}

func (m MQ) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return db.User{}, errors.New("user not found")
}
func Test_middleware_Authorize(t *testing.T) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	mockQ := MQ{nil, make([]db.User, 0)}
	middleware := InitMiddleware(logger, mockQ)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &middleware{
				logger:  tt.fields.logger,
				Querier: tt.fields.Querier,
			}
			if got := a.Authorize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("middleware.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}
