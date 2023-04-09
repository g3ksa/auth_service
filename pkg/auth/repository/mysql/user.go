package mysql

import (
	"AuthService/pkg/auth"
	"AuthService/pkg/hasher"
	"AuthService/pkg/models"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
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

func (r *UserRepository) Token(userId int, refreshToken string) (*models.User, error) {
	//refreshToken == bd.refreshToken -> new access -> to client

	user := new(models.User)
	query := r.db.QueryRow(
		`SELECT hashedRt FROM Users WHERE id = ?`,
		userId,
	)
	var hashedRt string
	err := query.Scan(&hashedRt)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден")
	}
	rtSecret, _ := os.LookupEnv("refreshSecret")
	if !hasher.DoHashMatch(hashedRt, refreshToken, []byte(rtSecret)) {
		return nil, fmt.Errorf("неверный токен: %s", refreshToken)
	}

	res := r.db.QueryRow(
		`SELECT u.companyCode, u.roleCode, GROUP_CONCAT(rr.rightsCode) AS rightsCodes
		FROM roles_rights_rights AS rr
		LEFT JOIN Users AS u ON rr.rolesCode = u.roleCode
		WHERE u.id = ?`, userId)

	var userRights string
	err = res.Scan(&user.UserCompany, &user.UserRole, &userRights)
	user.UserRights = strings.Split(userRights, ",")
	user.Id = userId

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("ошибка при получении данных пользователя")
	}

	return user, nil
}

func (r *UserRepository) Delete(userId int) error {
	//cookie refresh -> userId -> delete rt db
	_, err := r.db.Exec("UPDATE Users SET hashedRt = ? WHERE id = ?", nil, userId)

	return err
}

func (r *UserRepository) Get(params auth.GetParams) (*models.User, error) {
	user := new(models.User)
	res := r.db.QueryRow(
		`SELECT u.id, u.password, u.companyCode, u.roleCode, GROUP_CONCAT(rr.rightsCode) AS rightsCodes
			FROM roles_rights_rights AS rr
			LEFT JOIN Users AS u ON rr.rolesCode = u.roleCode
			WHERE u.email = ?`, params.Email)

	fmt.Println(res)

	var userRights string
	err := res.Scan(&user.Id, &user.Password, &user.UserCompany, &user.UserRole, &userRights)
	user.UserRights = strings.Split(userRights, ",")
	user.Email = params.Email
	fmt.Println(err)

	if err != nil {
		return nil, auth.ErrUserDoesNotExist
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return nil, auth.ErrIncorrectPassword
	}

	return user, nil
}

func (r *UserRepository) PutToken(refreshToken string, userId int) error {
	fmt.Println(userId)
	fmt.Println(refreshToken)
	_, err := r.db.Exec("UPDATE Users SET hashedRt = ? WHERE id = ?", refreshToken, userId)

	return err
}
