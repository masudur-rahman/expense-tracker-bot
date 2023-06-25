package models

type User struct {
	ID               string `db:"id,pk"`
	Name             string
	Email            string `db:"email,uq"`
	ContactInfo      string
	Balance          float64
	LastTxnTimestamp int64
}

//func (u *User) APIFormat() gqtypes.User {
//	return gqtypes.User{
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
