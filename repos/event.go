package repos

import "github.com/masudur-rahman/expense-tracker-bot/models"

type EventRepository interface {
	AddEvent(event string) error
	ListEvents() ([]models.Event, error)
}
