package user

import (
	"github.com/masudur-rahman/expense-tracker-bot/models"
	"github.com/masudur-rahman/expense-tracker-bot/repos"
	"github.com/masudur-rahman/expense-tracker-bot/services"
)

type drCrService struct {
	drCrRep repos.DebtorCreditorRepository
}

var _ services.DebtorCreditorService = &drCrService{}

func NewDebtorCreditorService(drCrRepo repos.DebtorCreditorRepository) *drCrService {
	return &drCrService{
		drCrRep: drCrRepo,
	}
}

func (u *drCrService) GetDebtorCreditorByID(id int64) (*models.DebtorsCreditors, error) {
	return u.drCrRep.GetDebtorCreditorByID(id)
}

func (u *drCrService) GetDebtorCreditorByName(userID int64, name string) (*models.DebtorsCreditors, error) {
	return u.drCrRep.GetDebtorCreditorByName(userID, name)
}

func (u *drCrService) ListDebtorCreditors(userID int64) ([]models.DebtorsCreditors, error) {
	return u.drCrRep.ListDebtorCreditors(userID)
}

func (u *drCrService) CreateDebtorCreditor(user *models.DebtorsCreditors) error {
	return u.drCrRep.AddNewDebtorCreditor(user)
}

func (u *drCrService) UpdateDebtorCreditorBalance(id int64, amount float64) error {
	return u.drCrRep.UpdateDebtorCreditorBalance(id, amount)
}

func (u *drCrService) DeleteDebtorCreditor(id int64) error {
	return u.drCrRep.DeleteDebtorCreditor(id)
}
