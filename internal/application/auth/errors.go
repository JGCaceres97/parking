package auth

import "errors"

var ErrInvalidToken = errors.New("token JWT inválido o mal formado")
var ErrExpiredToken = errors.New("token JWT expirado")
