package domain

import "errors"

var (
	ErrAdminProtected     = errors.New("no puedes hacer cambios sobre 'admin'")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUserInactive       = errors.New("usuario bloqueado o inactivo")
)

var (
	ErrUserNotFound                 = errors.New("usuario no encontrado")
	ErrVehicleTypeNotFound          = errors.New("tipo de vehículo no encontrado")
	ErrParkingRecordNotFound        = errors.New("registro de estacionamiento no encontrado")
	ErrActiveParkingNotFound        = errors.New("no se encontró un registro de entrada activo para esta placa")
	ErrUsernameAlreadyExists        = errors.New("nombre de usuario ya existe")
	ErrVehicleTypeNameAlreadyExists = errors.New("nombre de tipo de vehículo ya existe")
	ErrActiveParkingAlreadyExists   = errors.New("ya existe un registro de estacionamiento abierto para esta placa")
	ErrVehicleTypeInUse             = errors.New("tipo de vehículo está actualmente en uso")
)
