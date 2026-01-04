package response

import "errors"

var (
	ErrUserIDNotInContext    = errors.New("identidad del usuario no disponible")
	ErrMissingToken          = errors.New("falta el token de autenticación")
	ErrMissingMetadata       = errors.New("falta la metadata")
	ErrInvalidTokenFormat    = errors.New("formato de token inválido")
	ErrExpiredToken          = errors.New("token expirado")
	ErrInvalidToken          = errors.New("token inválido")
	ErrTokenValidationFailed = errors.New("error al validar token")
	ErrPermissionDenied      = errors.New("permiso denegado")

	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUserBlocked        = errors.New("usuario inactivo o bloqueado")
)

var (
	ErrRegistryIDRequired   = errors.New("ID de registro es requerido")
	ErrPlateRequired        = errors.New("la placa es requerida")
	ErrPlateAndTypeRequired = errors.New("placa y tipo de vehículo son requeridos")
	ErrUserCreateValidation = errors.New("el nombre de usuario, contraseña y rol son requeridos")
	ErrInvalidID            = errors.New("ID de usuario inválido o ausente")
	ErrInvalidRole          = errors.New("rol de usuario inválido. Los roles permitidos son 'admin' y 'common'")
	ErrChangeOwnRole        = errors.New("no puedes cambiar tu propio rol")
	ErrOwnDelete            = errors.New("no puedes eliminarte a ti mismo")
	ErrUpdateValidation     = errors.New("al menos un campo (username, rol, is_active) debe ser proporcionado para la actualización")
)

var (
	ErrMethodNotAllowed = errors.New("método no permitido")
	ErrTimeout          = errors.New("tiempo de espera excedido")
	ErrInternalError    = errors.New("error interno del servidor")
)
