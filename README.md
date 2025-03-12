# Go Echo Backend

Backend básico construido con Go y Echo Framework para crear una api rest.

## Estructura del Proyecto

```
.
├── config/         # Configuración de la aplicación
├── FateBakend/     # almacenamiento temporal de las Tasks
├── handlers/       # Manejadores HTTP
├── routes/         # Definición de rutas
├── types/         # Definición de structut 
├── verification/  # Definición de validadores
├── .env           # Variables de entorno
├── go.mod         # Dependencias
├── main.go        # Punto de entrada
└── README.md      # Este archivo
```

## Requisitos

- Go 1.21 o superior
- Echo Framework
- validator
- testify

## Instalación

1. Clona el repositorio
```bash
git clone https://github.com/Chris-medi/Backend-go.git
```

2. Instala las dependencias:
```bash
go mod tidy
```

## Ejecutar el proyecto

```bash
go run main.go
```

## Ejecutar los test

```bash
go test -v -cover ./...
```

El servidor estará disponible en `http://localhost:8080`

## Endpoints disponibles

- `GET /api/v1/tasks` - Get all Task
- `GET /api/v1/task/:id` - Get Task by id
- `DELETE /api/v1/task/:id` - Delete Task by id
- `PUT /api/v1/task` - Update Task
- `POST /api/v1/task` - Create Task



## Nota:
- Para actualizar una Tasks se carga el id en el body lo hice para que sea mas facil de actualizar
