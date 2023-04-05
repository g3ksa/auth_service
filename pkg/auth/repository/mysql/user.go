package mysql

import (
	"AuthService/pkg/models"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Refresh(userId int, refreshToken string) (string, error) {
	//refreshToken == bd.refreshToken -> new access -> to client
	return "", nil
}

func (r *UserRepository) Logout() error {
	//coock refresh -> userId -> delete rt db
	return nil
}

func (r *UserRepository) GetUser(email, password string) (*models.User, error) {
	user := new(models.User)
	res, _ := r.db.Query(fmt.Sprintf("SELECT * FROM Users;"))

	for res.Next() {
		err := res.Scan(&user.Id, &user.CreatedAt, &user.ImageUrl, &user.Email, &user.Password, &user.HashedRt)

		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
