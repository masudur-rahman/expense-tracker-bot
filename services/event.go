package services

import "github.com/masudur-rahman/expense-tracker-bot/models"

type EventService interface {
	AddEvent(event string) error
	ListEvents() ([]models.Event, error)
}
