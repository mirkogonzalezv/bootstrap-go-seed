# Microservice Go

- Semilla base para microservices en Go con Onion Architecture

## Tabla de Contenidos
1. [Descripción](#descripción)
2. [Requisitos](#requisitos)
3. [Estructura](#estructura)
4. [Quick Start](#quick-start)
5. [Comandos](#comandos)
6. [Guía para Nuevos Módulos](#guía-para-nuevos-módulos)
7. [Configuración](#configuración)
8. [Docker](#docker)

## Requisitos
- Go 1.25+
- (Opcional) Air para hot-reload

## Estructura
microservice/
├── cmd/                    # Punto de entrada
│   ├── api/main.go         # Inicio del servidor
│   ├── app/                # Configuración de la app
│   ├── config/             # Variables de entorno
│   ├── container/          # Inyección de dependencias
│   └── routes/             # Registro de rutas
├── internal/
│   └── modules/            # Módulos de negocio
│       └── health/
│           ├── application/     # Casos de uso
│           │   ├── dtos/       # Objetos de transferencia
│           │   └── use_cases/  # Lógica de negocio
│           └── presentation/    # Capa de presentación
│               ├── handler/    # Controladores
│               └── router/     # Definición de rutas
└── pkg/                    # Utilidades compartidas
    ├── env/                # Entornos
    ├── errors/             # Errores personalizados
    ├── http/               # Utilidades HTTP
    ├── logger/             # Logger con Zap
    ├── middlewares/       # Middlewares Gin
    └── version/            # Versionado de API

## Quick Start
```
# Desarrollo con hot-reload
make start-dev

# Producción
make start
```

# Comandos

make start-dev	      Iniciar con hot-reload
make start	          Iniciar aplicación
make test	            Ejecutar tests
make test-coverage	  Tests con coverage
make docs-generate	  Generar documentación OpenAPI
make security-scan	  Escaneo de secretos
make audit	          Auditoría de dependencias

# Configuración
Crear archivo .env.dev basado en .env.example:
PORT=4200
ENVIRONMENT=development
LOG_LEVEL=info

# API

### Health Check
- GET /api/v1/health

Respuesta:
{
  "timestamp": "2024-06-01T12:00:00Z",
  "status": "microservice"
}

Docker

# Build
docker build -t microservice .

# Run
docker run -p 4200:4200 microservice

# Guía: Agregar Nuevo Módulo

