package user_model

import "database/sql"

type User struct {
	UserId      string
	Email       sql.NullString
	Phone       sql.NullString
	Password    string
	Name        string
	ImageUrl    string
	AccessToken string
}
