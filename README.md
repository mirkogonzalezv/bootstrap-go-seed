# xintecgo — Microservice Seed

Seed base para microservicios en Go con **Clean Architecture**.  
Diseñada para que te preocupes solo de los módulos de negocio en `internal/`.

## Requisitos

- Go 1.25+
- (Opcional) Air para hot-reload → `make install-air`

## Quick Start

```bash
cp .env.example .env.dev
make start-dev          # hot-reload
# o
make start              # modo normal
```

Documentación Swagger: `http://localhost:3200/docs/index.html`

---

## Estructura del proyecto

```
cmd/                        # Bootstrap de la app (no tocar)
├── api/main.go             # Punto de entrada
├── app/                    # Init, Run, Shutdown (lifecycle)
├── config/                 # Variables de entorno
├── container/              # Inyección de dependencias
└── routes/                 # Registro de rutas

internal/                   # ← Módulos de negocio (aquí trabajas)
└── <modulo>/
    ├── application/        # Lógica de negocio
    │   ├── dto/            # Objetos de transferencia
    │   └── usecase/        # Casos de uso
    ├── domain/             # Entidades e interfaces de repositorio
    ├── infrastructure/     # Implementaciones concretas (DB, APIs)
    └── presentation/       # Capa HTTP
        ├── handler/        # Controladores Gin
        └── router/         # Definición de rutas

pkg/                        # Utilidades compartidas
├── apperrors/              # Errores de dominio tipados
├── env/                    # Tipos de entorno (dev, qa, prod)
├── httputil/               # Mapper errores de dominio → HTTP
├── logger/                 # Logger con Zap
└── middleware/             # Middlewares Gin (error handler)
```

---

## Cómo crear un nuevo módulo

Ejemplo: módulo `tareas` con endpoint `GET /api/v1/tareas`.

### 1. Crear la estructura de carpetas

```bash
mkdir -p internal/tareas/application/{dto,usecase}
mkdir -p internal/tareas/domain
mkdir -p internal/tareas/infrastructure
mkdir -p internal/tareas/presentation/{handler,router}
```

### 2. Definir el DTO de respuesta

`internal/tareas/application/dto/tarea_response.go`:

```go
package dto

type TareaResponse struct {
    ID    string `json:"id"`
    Nombre string `json:"nombre"`
}

func NewTareaResponse(id, nombre string) *TareaResponse {
    return &TareaResponse{ID: id, Nombre: nombre}
}
```

### 3. Crear el caso de uso

`internal/tareas/application/usecase/obtener_tarea.go`:

```go
package usecase

import "microservice/internal/tareas/application/dto"

type ObtenerTareaUseCase struct{}

func NewObtenerTareaUseCase() *ObtenerTareaUseCase {
    return &ObtenerTareaUseCase{}
}

func (uc *ObtenerTareaUseCase) Execute() (*dto.TareaResponse, error) {
    return dto.NewTareaResponse("1", "Mi tarea"), nil
}
```

### 4. Crear el handler HTTP

`internal/tareas/presentation/handler/tarea.go`:

```go
package handler

import (
    "net/http"

    "microservice/internal/tareas/application/usecase"

    "github.com/gin-gonic/gin"
)

type TareaHandler struct {
    uc *usecase.ObtenerTareaUseCase
}

func NewTareaHandler(uc *usecase.ObtenerTareaUseCase) *TareaHandler {
    return &TareaHandler{uc: uc}
}

func (h *TareaHandler) GetTarea(c *gin.Context) {
    res, err := h.uc.Execute()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, res)
}
```

### 5. Definir las rutas

`internal/tareas/presentation/router/tarea.go`:

```go
package router

import (
    "microservice/internal/tareas/presentation/handler"

    "github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.TareaHandler) {
    v1 := api.Group("/v1/tareas")
    v1.GET("", h.GetTarea)
}
```

### 6. Registrar en el container de dependencias

`cmd/container/container.go`:

```go
import tareaUsecase "microservice/internal/tareas/application/usecase"
import tareaHandler "microservice/internal/tareas/presentation/handler"

type Container struct {
    HealthHandler *healthHandler.HealthHandler
    TareaHandler  *tareaHandler.TareaHandler   // ← agregar
    Log           *zap.Logger
}

func NewContainer(cfg *config.Configuration, log *zap.Logger) *Container {
    // ...
    tareaUC := tareaUsecase.NewObtenerTareaUseCase()
    tareaHdlr := tareaHandler.NewTareaHandler(tareaUC)

    return &Container{
        HealthHandler: healthHdlr,
        TareaHandler:  tareaHdlr,              // ← agregar
        Log:           log,
    }
}
```

### 7. Registrar la ruta

`cmd/routes/register.go`:

```go
import tareaRouter "microservice/internal/tareas/presentation/router"

func ConfigureModule(router *gin.Engine, c *container.Container) {
    api := router.Group("/api")
    health.RegisterRoutes(api, c.HealthHandler)
    tareaRouter.RegisterRoutes(api, c.TareaHandler)   // ← agregar
}
```

✅ **Listo** — tu endpoint ya responde en `GET /api/v1/tareas`.

---

## Logger

Usa `go.uber.org/zap`. Dos formas de usarlo:

**En módulos de negocio (`internal/`)** — helpers simples, sin acoplamiento:

```go
import "microservice/pkg/logger"

logger.Success("todo ok")
logger.Error("algo falló: %s", err)
logger.General("mensaje informativo")
```

**En bootstrap e infraestructura** — logger inyectado vía DI:

```go
func ErrorHandlerMiddleware(log *zap.Logger) gin.HandlerFunc {
    // usa log.Error(), log.Warn(), etc.
}
```

---

## Manejo de errores

Usa errores de dominio tipados en `pkg/apperrors/`:

```go
import "microservice/pkg/apperrors"

// Según el tipo, el middleware mapea automáticamente al HTTP status code correcto
apperrors.NewNotFoundError("USER_001", "Usuario no encontrado")
apperrors.NewValidationError("EMAIL_001", "Email inválido")
apperrors.NewUnauthorizedError("AUTH_001", "Token expirado")
apperrors.NewInternalError("DB_001", "Error de conexión")
```

| Error | HTTP Status |
|---|---|
| `Validation` | 400 |
| `NotFound` | 404 |
| `Conflict` | 409 |
| `Unauthorized` | 401 |
| `Forbidden` | 403 |
| `Internal` | 500 |

---

## Configuración

Las variables de entorno se cargan automáticamente desde `.env.{environment}` (dev/qa) o del sistema (prod).

| Variable | Default | Descripción |
|---|---|---|
| `PORT` | `3200` | Puerto del servidor |
| `ENVIRONMENT` | `dev` | Entorno: `dev`, `qa` o `prod` |
| `LOG_LEVEL` | `info` | Nivel de log |
| `SHUTDOWN_TIMEOUT` | `30` | Timeout de graceful shutdown (seg) |

```bash
cp .env.example .env.dev
# editar .env.dev con tus valores
```

---

## Comandos disponibles

```bash
make help              # Lista todos los comandos
make start-dev         # Iniciar con hot-reload (Air)
make start             # Iniciar en modo normal
make test              # Ejecutar tests de internal/
make test-coverage     # Tests con reporte de cobertura
make docs-generate     # Generar documentación Swagger
make docs-serve        # Iniciar con documentación ReDoc
make audit             # Auditoría de dependencias (govulncheck)
make security-scan     # Escaneo de secretos (gitleaks)
make deps-check        # Verificar y limpiar dependencias
```

---

## Documentación API (Swagger)

La documentación se genera desde las anotaciones en los handlers:

```go
// @Summary Obtener tarea
// @Description Retorna una tarea por ID
// @Tags Tareas
// @Success 200 {object} dto.TareaResponse
// @Router /tareas [get]
```

Generar y visualizar:

```bash
make docs-generate
# Abrir http://localhost:3200/docs/index.html (solo dev/qa)
```

---

## Docker

```bash
docker build -t microservice .
docker run -p 3200:3200 microservice
```

**Características de la imagen:**
- Base: `gcr.io/distroless/static-debian12:nonroot`
- Tamaño: ~15 MB
- Corre como non-root (mejores prácticas)
- Certificados CA incluidos (HTTPS listo)
- Sin shell, sin package manager (mínima superficie de ataque)

Compatible con: **Cloud Run, AWS ECS/Fargate, Railway, Azure Container Apps**.

---

## Convenciones Go usadas

| Regla | Ejemplo |
|---|---|
| Sin shadow de stdlib | `pkg/apperrors`, `pkg/httputil` |
| Packages en singular | `middleware`, `usecase`, `handler` |
| Sin prefijo `Get` en getters | `Version()` en vez de `GetVersion()` |
| Archivos en snake_case | `error_handler.go`, `health_response.go` |
| Package name = nombre del directorio | `package dto` dentro de `dto/` |
