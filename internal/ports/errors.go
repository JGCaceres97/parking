package ports

import "errors"

// Errores de autenticación y lógica de negocio (Service)
var ErrInvalidCredentials = errors.New("credenciales inválidas")
var ErrUserBlocked = errors.New("usuario bloqueado o inactivo")
var ErrInvalidToken = errors.New("token JWT inválido o mal formado")
var ErrTokenExpired = errors.New("token JWT expirado")
var ErrRequiredUpdateFields = errors.New("los campos de salida del registro no deben ser nulos al actualizar")

// Errores de persistencia (Repository)
var ErrUserNotFound = errors.New("usuario no encontrado")
var ErrUsernameExists = errors.New("nombre de usuario ya existe")
var ErrVehicleTypeNotFound = errors.New("tipo de vehículo no encontrado")
var ErrVehicleTypeNameExists = errors.New("nombre de tipo de vehículo ya existe")
var ErrVehicleTypeInUse = errors.New("tipo de vehículo está actualmente en uso")
var ErrParkingRecordNotFound = errors.New("registro de estacionamiento no encontrado")

var ErrActiveParkingExists = errors.New("ya existe un registro de estacionamiento abierto para esta placa")
var ErrActiveParkingNotFound = errors.New("no se encontró un registro de entrada activo para esta placa")
