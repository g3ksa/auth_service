package mysql

import (
	"AuthService/pkg/models"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"strings"
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

func (r *UserRepository) Get(email, password string) (*models.User, error) {
	user := new(models.User)
	res := r.db.QueryRow(
		`SELECT u.id, u.password, u.companyCode, GROUP_CONCAT(ur.rightsCode) AS rightsCodes
		FROM users_rights_rights AS ur
		LEFT JOIN Users AS u ON ur.usersId = u.id
		WHERE u.email = ?`, email)

	var userRights string
	err := res.Scan(&user.Id, &user.Password, &user.UserCompany, &userRights)
	user.UserRights = strings.Split(userRights, ",")
	user.Email = email

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}
