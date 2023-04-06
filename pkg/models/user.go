package models

type User struct {
	Id          int      `json:"id"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	UserCompany string   `json:"userCompany"`
	UserRights  []string `json:"userRights"`
}
