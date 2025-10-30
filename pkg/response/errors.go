package response

import "errors"

var ErrRegistryIDRequired = errors.New("ID de registro es requerido")
var ErrPlateRequired = errors.New("la placa es requerida")
var ErrPlateAndTypeRequired = errors.New("placa y tipo de vehículo son requeridos")
var ErrUserCreateValidation = errors.New("el nombre de usuario, contraseña y rol son requeridos")
var ErrInvalidID = errors.New("ID de usuario inválido o ausente")
var ErrInvalidRole = errors.New("rol de usuario inválido. Los roles permitidos son 'admin' y 'common'")
var ErrChangeOwnRole = errors.New("no puedes cambiar tu propio rol")
var ErrOwnDelete = errors.New("no puedes eliminarte a ti mismo")
var ErrUpdateValidation = errors.New("al menos un campo (username, rol, is_active) debe ser proporcionado para la actualización")

var ErrUserIDNotInContext = errors.New("identidad del usuario no disponible")
var ErrMissingToken = errors.New("falta el token de autenticación")
var ErrInvalidTokenFormat = errors.New("formato de token inválido")
var ErrTokenExpired = errors.New("token expirado")
var ErrInvalidToken = errors.New("token inválido")
var ErrTokenValidationFailed = errors.New("error al validar token")
var ErrPermissionDenied = errors.New("permiso denegado")

var ErrInvalidJSON = errors.New("solicitud JSON inválida")
var ErrInvalidCredentials = errors.New("credenciales inválidas")
var ErrUserBlocked = errors.New("usuario inactivo o bloqueado")
var ErrTimeout = errors.New("tiempo de espera excedido")
var ErrInternalError = errors.New("error interno del servidor")
