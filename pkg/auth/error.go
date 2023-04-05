package auth

import "errors"

var ErrInvalidAccessToken = errors.New("Некорректный токен авторизаци")
var ErrUserDoesNotExist = errors.New("Такого польхователя не существует")
var ErrUserAlreadyExists = errors.New("Пользователь с такими данными уже существует")
