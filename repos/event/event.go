package event

import (
	"github.com/masudur-rahman/expense-tracker-bot/infra/logr"
	"github.com/masudur-rahman/expense-tracker-bot/models"

	isql "github.com/masudur-rahman/styx/sql"
)

type SQLEventRepository struct {
	db     isql.Engine
	logger logr.Logger
}

func NewSQLEventRepository(db isql.Engine, logger logr.Logger) *SQLEventRepository {
	return &SQLEventRepository{
		db:     db.Table("event"),
		logger: logger,
	}
}

func (e *SQLEventRepository) AddEvent(event string) error {
	//TODO implement me
	panic("implement me")
}

func (e *SQLEventRepository) ListEvents() ([]models.Event, error) {
	//TODO implement me
	panic("implement me")
}
