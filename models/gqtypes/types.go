package gqtypes

type RegisterParams struct {
	FirstName string `json:"firstName" form:"firstName"`
	LastName  string `json:"lastName" form:"lastName"`

	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LoginParams struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
