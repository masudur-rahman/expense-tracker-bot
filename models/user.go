package models

type DebtorsCreditors struct {
	ID               int64  `db:"id,pk"`
	UserID           int64  `db:",uqs"`
	NickName         string `db:",uqs"`
	FullName         string
	Email            string `db:"email,uqs"`
	ContactInfo      string
	Balance          float64
	LastTxnTimestamp int64
}

type User struct {
	ID         int64  `db:"id,pk"`
	TelegramID int64  `db:",uq"`
	Username   string `db:",uq"`
	FirstName  string
	LastName   string
}

//func (u *DebtorsCreditors) APIFormat() gqtypes.DebtorsCreditors {
//	return gqtypes.DebtorsCreditors{
//		ID:        u.ID,
//		Username:  u.Username,
//		Email:     u.Email,
//		FirstName: u.FirstName,
//		LastName:  u.LastName,
//		//FullName: fmt.Sprintf("%s %s", u.FirstName, u.LastName),
//		Bio:      u.Bio,
//		Location: u.Location,
//		Avatar:   u.Avatar,
//		IsActive: u.IsActive,
//		IsAdmin:  u.IsAdmin,
//	}
//}
