# Sistema de Gestión de Estacionamiento

Este proyecto implementa una API para gestionar la entrada, salida y cobro de tarifas en un
estacionamiento, aplicando reglas de negocio específicas para el cálculo de tiempo y tarifas por
tipo de vehículo.

# Índice

- [Modelo de Datos](#-modelo-de-datos-esquema-mysql)
  - [Relaciones Clave](#relaciones-clave)
  - [Tabla: USERS](#tabla-users)
  - [Tabla: VEHICLE_TYPES](#tabla-vehicle-types)
  - [Tabla: PARKING_RECORDS](#tabla-parking_records-transaccional)
- [Reglas de Negocio para el Cálculo de Tarifas](#-reglas-de-negocio-para-el-cálculo-de-tarifas)
- [Dependencias](#-dependencias)
  - [Entorno de Desarrollo](#entorno-de-desarrollo)
  - [Dependencias de Go](#dependencias-de-go)
- [Ejecución del Proyecto](#-ejecucion-del-proyecto)
  - [Clonación del repositorio](#0-clonación-del-repositorio)
  - [Preparación de Archivos](#1-preparación-de-archivos)
  - [Levantar los Servicios](#2-levantar-los-servicios)

## 💾 Modelo de Datos (Esquema MySQL)

La base de datos se compone de tres tablas principales que gestionan la información de usuarios,
tipos de vehículos y los registros transaccionales de estacionamiento.

### Relaciones Clave

- USERS ⬅️ PARKING_RECORDS: Un usuario (quien opera el sistema) puede generar múltiples registros de
  estacionamiento.
- VEHICLE_TYPES ⬅️ PARKING_RECORDS: Un tipo de vehículo se asocia a múltiples registros para aplicar
  su tarifa correspondiente.

### Tabla: USERS

Almacena la información de autenticación y el rol de los operadores del sistema.

| Columna       | Tipo de Dato            | Clave | Restricciones             | Propósito                        |
| ------------- | ----------------------- | ----- | ------------------------- | -------------------------------- |
| id            | VARCHAR(26)             | PK    | NOT NULL, ULID            | Identificador único del usuario. |
| username      | VARCHAR(255)            |       | UNIQUE, NOT NULL          | Nombre de usuario                |
| password_hash | VARCHAR(255)            |       | NOT NULL                  | Hash seguro de contraseña.       |
| role          | ENUM('admin', 'common') |       | NOT NULL                  | Permisos de usuario.             |
| is_active     | BOOLEAN                 |       | DEFAULT TRUE              | Estado del usuario.              |
| created_at    | TIMESTAMP               |       | DEFAULT CURRENT_TIMESTAMP | Fecha de creación                |

### Tabla: VEHICLE_TYPES

Define las categorías de vehículos y las tarifas horarias que rigen el cálculo del cobro.

| Columna     | Tipo de Dato   | Clave | Restricciones    | Propósito                                         |
| ----------- | -------------- | ----- | ---------------- | ------------------------------------------------- |
| id          | VARCHAR(26)    | PK    | NOT NULL, ULID   | Identificador único del tipo de vehículo.         |
| name        | VARCHAR(50)    |       | UNIQUE, NOT NULL | Nombre del tipo (ej., 'Motocicleta', 'Especial'). |
| hourly_rate | DECIMAL(10, 2) |       | NOT NULL         | Tarifa por hora (ej., 15.00, 5.00 o 0.00).        |
| description | VARCHAR(255)   |       |                  | Descripción opcional del tipo.                    |

### Tabla: PARKING_RECORDS (Transaccional)

Contiene el registro de cada estadía, incluyendo el cálculo final del cargo.

| Columna          | Tipo de Dato   | Clave | Restricciones                | Propósito                                                 |
| ---------------- | -------------- | ----- | ---------------------------- | --------------------------------------------------------- |
| id               | VARCHAR(26)    | PK    | NOT NULL, ULID               | ID único del registro de estacionamiento.                 |
| user_id          | VARCHAR(26)    | FK    | NOT NULL, Ref: USERS         | Usuario que registró la entrada.                          |
| vehicle_type_id  | VARCHAR(26)    | FK    | NOT NULL, Ref: VEHICLE_TYPES | Tipo de vehículo para determinar la tarifa.               |
| license_plate    | VARCHAR(10)    |       | NOT NULL                     | Placa del vehículo.                                       |
| entry_time       | DATETIME       |       | NOT NULL                     | Hora y fecha de entrada (en UTC).                         |
| exit_time        | DATETIME       |       | NULL                         | Hora y fecha de salida. NULL si el vehículo sigue dentro. |
| total_charge     | DECIMAL(10, 2) |       | NULL                         | Cargo total calculado al momento de la salida.            |
| calculated_hours | INT            |       | NULL                         | Horas cobradas aplicando la lógica de redondeo.           |

## 💸 Reglas de Negocio para el Cálculo de Tarifas

El cálculo de las tarifas se basa en el tiempo transcurrido entre el timpo de entrada y de salida,
siguiendo estas directrices:

- Tarifa Base: $15 USD por hora.
- Tarifa Especial: $5 USD por hora (para "Vehículos Especiales").
- Exentos de Pago: Motocicletas ($0 USD por hora).
- Mínimo de Cobro: Toda estadía (no exenta) tiene un cargo mínimo de 1 hora.
- Regla de Redondeo: A partir de la segunda hora, cualquier fracción de tiempo igual o superior a 30
  minutos se redondea a la hora completa siguiente. (Ej: 1h 29m = 1h, 1h 30m = 2h).

## 📦 Dependencias

El proyecto está construido en Go y requiere las siguientes dependencias externas y herramientas:

### Entorno de Desarrollo

- Go: Versión 1.24 o superior.
- MySQL: Base de datos relacional para persistencia de datos.
- Frontend: React 19, TypeScript y Vite para la interfaz de usuario.
- Docker y Docker Compose: Esenciales para levantar el servicio del API y la base de datos MySQL.

### Dependencias de Go

- [github.com/go-chi/chi/v5](https://github.com/go-chi/chi/v5): Router HTTP ligero y modular.
- [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql): Driver para conexión a
  MySQL.
- [github.com/oklog/ulid/v2](https://github.com/oklog/ulid/v2): Generación de identificadores únicos
  (ULID).
- [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt/v5): Manejo de JSON Web Tokens
  (JWT) para autenticación.
- [golang.org/x/crypto](https://golang.org/x/crypto): Funcionalidades de criptografía (ej. hashing
  de contraseñas).
- [github.com/joho/godotenv](https://github.com/joho/godotenv): Carga de variables de entorno desde
  archivos .env.

## 🚀 Ejecución del Proyecto

Siga los siguientes pasos para levantar y ejecutar la API y la base de datos utilizando Docker
Compose:

### 0. Clonación del repositorio

Primero, clone el repositorio a su maquina local y navegue al directorio del proyecto:

```bash
git clone https://github.com/jgcaceres97/parking
cd parking
```

### 1. Preparación de archivos

Asegúrese de tener configurados los siguientes archivos en la raíz del proyecto:

- Dockerfile (para construir la imagen de la API Go).
- docker-compose.yml (para definir y conectar los servicios de API y MySQL).
- .env (para las variables de entorno, incluyendo las credenciales de MySQL y el secreto JWT).

### 2. Ejecutar los servicios

Ejecute el siguiente comando para construir las imágenes (si es necesario) y levantar los
contenedores de la API y la Base de Datos:

```bash
docker compose up --build
```

El servicio estará disponible en el puerto expuesto por Docker Compose, típicamente
http://localhost:3000.

### 3. Primer inicio de sesión

Al iniciar el servicio por primera vez, el sistema crea un usuario `administrador` por defecto con las siguientes credenciales:
```bash
username: admin
password: admin
```
