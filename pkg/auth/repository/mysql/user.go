package mysql

import (
	"AuthService/pkg/auth"
	"AuthService/pkg/models"
	"database/sql"
	"fmt"
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

func (r *UserRepository) Refresh(userId int, refreshToken string) (*models.User, error) {
	//refreshToken == bd.refreshToken -> new access -> to client

	user := new(models.User)
	query := r.db.QueryRow(
		`SELECT hashedRt FROM Users WHERE id = ?`,
		userId,
	)
	var hashedRt string
	var email string
	err := query.Scan(&hashedRt, &email)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден")
	}
	if !strings.EqualFold(refreshToken, hashedRt) {
		return nil, fmt.Errorf("неверный токен: %s", refreshToken)
	}

	res := r.db.QueryRow(
		`SELECT u.companyCode, GROUP_CONCAT(ur.rightsCode) AS rightsCodes
		FROM users_rights_rights AS ur
		LEFT JOIN Users AS u ON ur.usersId = u.id
		WHERE u.id = ?`, userId)

	var userRights string
	err = res.Scan(&user.UserCompany, &userRights)
	user.UserRights = strings.Split(userRights, ",")
	user.Id = userId

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных пользователя")
	}

	return user, nil
}

func (r *UserRepository) Logout(userId int) error {
	//cookie refresh -> userId -> delete rt db
	_, err := r.db.Exec("UPDATE Users SET hashedRt = ? WHERE id = ?", nil, userId)
	if err != nil {
		return fmt.Errorf("ошибка при запросе к бд: %s", err)
	}
	return nil
}

func (r *UserRepository) Get(params auth.GetParams) (*models.User, error) {
	user := new(models.User)
	res := r.db.QueryRow(
		`SELECT u.id, u.password, u.companyCode, GROUP_CONCAT(ur.rightsCode) AS rightsCodes
		FROM users_rights_rights AS ur
		LEFT JOIN Users AS u ON ur.usersId = u.id
		WHERE u.email = ?`, params.Email)

	var userRights string
	err := res.Scan(&user.Id, &user.Password, &user.UserCompany, &userRights)
	user.UserRights = strings.Split(userRights, ",")
	user.Email = params.Email

	if err != nil {
		return nil, fmt.Errorf("пользователь не найден")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return nil, fmt.Errorf("неверный пароль")
	}

	return user, nil
}
