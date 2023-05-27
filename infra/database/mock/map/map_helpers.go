package _map

import (
	"net/http"

	"github.com/masudur-rahman/expense-tracker-bot/models"
)

func checkEntityNameNonEmpty(entity string) error {
	if entity == "" {
		return models.StatusError{
			Message: "entity name must be set",
		}
	}
	return nil
}

func checkIDNonEmpty(id string) error {
	if id == "" {
		return models.StatusError{
			Status:  http.StatusBadRequest,
			Message: "must provide document id",
		}
	}
	return nil
}

func checkIdOrFilterNonEmpty(id string, filter interface{}) error {
	if id == "" && filter == nil {
		return models.StatusError{
			Status:  http.StatusBadRequest,
			Message: "must provide id and/or filter",
		}
	}
	return nil
}
