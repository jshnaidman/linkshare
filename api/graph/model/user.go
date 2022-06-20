package model

type User struct {
	Username  string   `json:"username"`
	FirstName *string  `json:"firstName"`
	Lastname  *string  `json:"lastname"`
	Email     *string  `json:"email"`
	GoogleID  *string  `json:"googleID"`
	PageURLs  []string `json:"PageURLs"`
}
