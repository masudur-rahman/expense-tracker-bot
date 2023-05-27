package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StatusError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (err StatusError) Error() string {
	jsonBytes, _ := json.Marshal(err)
	return string(jsonBytes)
}

func ParseStatusError(err error) (int, string) {
	serr := StatusError{}
	if perr := json.Unmarshal([]byte(err.Error()), &serr); perr != nil {
		return http.StatusInternalServerError, err.Error()
	}

	if serr.Status == 0 {
		serr.Status = http.StatusInternalServerError
	}
	return serr.Status, serr.Message
}

type ErrUserNotAuthenticated struct{}

func (err ErrUserNotAuthenticated) Error() string {
	return StatusError{
		Status:  http.StatusUnauthorized,
		Message: "user must be authenticated",
	}.Error()
}

type ErrUserNotFound struct {
	ID       string
	Username string
	Email    string
}

func (err ErrUserNotFound) Error() string {
	return StatusError{
		Status:  http.StatusNotFound,
		Message: fmt.Sprintf("user [id: %v, username: %v, email: %v] doesn't exist", err.ID, err.Username, err.Email),
	}.Error()
}

type ErrUserAlreadyExist struct {
	Username string
	Email    string
}

func (err ErrUserAlreadyExist) Error() string {
	return StatusError{
		Status:  http.StatusConflict,
		Message: fmt.Sprintf("user [username: %v, email: %v] already exist", err.Username, err.Email),
	}.Error()
}

type ErrUserPasswordMismatch struct{}

func (ErrUserPasswordMismatch) Error() string {
	return StatusError{
		Status:  http.StatusBadRequest,
		Message: "username or password is invalid",
	}.Error()
}

type ErrShelterNotFound struct {
	ID   string
	Name string
}

func (err ErrShelterNotFound) Error() string {
	return StatusError{
		Status:  http.StatusNotFound,
		Message: fmt.Sprintf("shelter [id: %v, name: %v] doesn't exist", err.ID, err.Name),
	}.Error()
}

type ErrShelterAlreadyExist struct {
	ID   string
	Name string
}

func (err ErrShelterAlreadyExist) Error() string {
	return StatusError{
		Status:  http.StatusConflict,
		Message: fmt.Sprintf("shelter [id: %v, name: %v] already exist", err.ID, err.Name),
	}.Error()
}

type ErrPetNotFound struct {
	ID   string
	Name string
}

func (err ErrPetNotFound) Error() string {
	return StatusError{
		Status:  http.StatusNotFound,
		Message: fmt.Sprintf("pet [id: %v, name: %v] doesn't exist", err.ID, err.Name),
	}.Error()
}

type ErrPetAdoptionNotFound struct {
	PetID string
}

func (err ErrPetAdoptionNotFound) Error() string {
	return StatusError{
		Status:  http.StatusNotFound,
		Message: fmt.Sprintf("pet adoption status not found for pet: %v", err.PetID),
	}.Error()
}

func IsErrNotFound(err error) bool {
	switch err.(type) {
	case ErrUserNotFound:
		return true
	case ErrShelterNotFound:
		return true
	case ErrPetNotFound:
		return true
	case ErrPetAdoptionNotFound:
		return true
	default:
		return false
	}
}

func IsErrConflict(err error) bool {
	switch err.(type) {
	case ErrUserAlreadyExist:
		return true
	case ErrShelterAlreadyExist:
		return true
	default:
		return false
	}
}

func IsErrBadRequest(err error) bool {
	switch err.(type) {
	case ErrUserPasswordMismatch:
		return true
	default:
		return false
	}
}
