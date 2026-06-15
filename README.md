# xintecgo — Microservice Seed

Seed base para microservicios en Go con **Clean Architecture**.  
Diseñada para que te preocupes solo de los módulos de negocio en `internal/`.

## Requisitos

- Go 1.25+
- (Opcional) Air para hot-reload → `make install-air`

## Quick Start

```bash
cp .env.example .env.dev
# completar vars DB en .env.dev si tienes PostgreSQL
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
├── container/              # Inyección de dependencias (composition root)
└── routes/                 # Registro de rutas

internal/                   # ← Módulos de negocio (aquí trabajas)
└── <modulo>/
    ├── application/        # Lógica de negocio
    │   ├── dto/            # Objetos de transferencia
    │   └── usecase/        # Casos de uso
    ├── domain/             # Entidades e interfaces de repositorio
    ├── infrastructure/     # Implementaciones concretas (DB, APIs externas)
    └── presentation/       # Capa HTTP
        ├── handler/        # Controladores Gin
        └── router/         # Definición de rutas

migrations/                 # SQL de migraciones (ejecutar a mano)
pkg/                        # Utilidades compartidas
├── apperrors/              # Errores de dominio tipados
├── database/               # Pool de PostgreSQL portable
├── env/                    # Tipos de entorno (dev, qa, prod)
├── httputil/               # Mapper errores de dominio → HTTP
├── logger/                 # Logger con Zap
└── middleware/             # Middlewares Gin (error handler, security)
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
    ID     string `json:"id"`
    Nombre string `json:"nombre"`
}

func NewTareaResponse(id, nombre string) *TareaResponse {
    return &TareaResponse{ID: id, Nombre: nombre}
}
```

### 3. Definir la interfaz del repositorio en domain

`internal/tareas/domain/repository.go`:

```go
package domain

import "context"

type TareaRepository interface {
    FindByID(ctx context.Context, id string) (*Tarea, error)
}
```

### 4. Crear el caso de uso

`internal/tareas/application/usecase/obtener_tarea.go`:

```go
package usecase

import (
    "context"
    "microservice/internal/tareas/application/dto"
    "microservice/internal/tareas/domain"
)

type ObtenerTareaUseCase struct {
    repo domain.TareaRepository
}

func NewObtenerTareaUseCase(repo domain.TareaRepository) *ObtenerTareaUseCase {
    return &ObtenerTareaUseCase{repo: repo}
}

func (uc *ObtenerTareaUseCase) Execute(ctx context.Context) (*dto.TareaResponse, error) {
    tarea, err := uc.repo.FindByID(ctx, "1")
    if err != nil {
        return nil, err
    }
    return dto.NewTareaResponse(tarea.ID, tarea.Nombre), nil
}
```

### 5. Crear la implementación del repositorio

`internal/tareas/infrastructure/postgres_repo.go`:

```go
package infrastructure

import (
    "context"
    "microservice/internal/tareas/domain"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTareaRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresTareaRepository(pool *pgxpool.Pool) *PostgresTareaRepository {
    return &PostgresTareaRepository{pool: pool}
}

func (r *PostgresTareaRepository) FindByID(ctx context.Context, id string) (*domain.Tarea, error) {
    // query a la DB...
}
```

### 6. Crear el handler HTTP

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
    res, err := h.uc.Execute(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, res)
}
```

### 7. Definir las rutas

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

### 8. Registrar en el container de dependencias

`cmd/container/container.go`:

```go
import (
    tareaDomain "microservice/internal/tareas/domain"
    tareaInfra "microservice/internal/tareas/infrastructure"
    tareaUsecase "microservice/internal/tareas/application/usecase"
    tareaHandler "microservice/internal/tareas/presentation/handler"
)

type Container struct {
    DB            *pgxpool.Pool
    HealthHandler *handler.HealthHandler
    TareaHandler  *tareaHandler.TareaHandler   // ← agregar
    Log           *zap.Logger
}

func NewContainer(cfg *config.Configuration, log *zap.Logger) (*Container, error) {
    // ... pool, health ...

    var tareaRepo tareaDomain.TareaRepository = tareaInfra.NewPostgresTareaRepository(pool)
    tareaUC := tareaUsecase.NewObtenerTareaUseCase(tareaRepo)
    tareaHdlr := tareaHandler.NewTareaHandler(tareaUC)

    return &Container{
        DB:            pool,
        HealthHandler: healthHdlr,
        TareaHandler:  tareaHdlr,              // ← agregar
        Log:           log,
    }, nil
}
```

### 9. Registrar la ruta

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

## Base de datos con Repository Pattern

El proyecto usa PostgreSQL con `pgx/v5` y sigue el patrón **Repository**:

```
Handler (HTTP)
  → UseCase (lógica de negocio)
    → Repository (interfaz en domain/)
      → PostgresRepository (implementación en infrastructure/)
        → pgxpool (pool de conexiones)
```

### Migraciones

Las migraciones son **SQL puro**, ejecutadas a mano. No hay ORM, no hay auto-migrate:

```bash
psql -d microservice -f migrations/001_create_health_check.up.sql
```

### Pool de conexiones

Se crea en `cmd/container/container.go` usando `pkg/database/postgres.go`:

```go
pool, err := database.NewPool(ctx, database.PostgresConfig{
    Host: cfg.DBHost, Port: cfg.DBPort, User: cfg.DBUser,
    Password: cfg.DBPassword, DBName: cfg.DBName, SSLMode: cfg.DBSSLMode,
})
```

El pool se **cierra automáticamente** en el shutdown de la app.

---

## Tests

Los tests de use cases usan **mocks manuales** (sin frameworks externos):

```go
type mockHealthRepository struct{}

func (m *mockHealthRepository) Ping(ctx context.Context) error {
    return nil
}

func TestHealthUseCase(t *testing.T) {
    mockRepo := &mockHealthRepository{}
    uc := NewHealthUseCase(mockRepo)

    res, err := uc.Execute(context.Background())

    if err != nil {
        t.Fatalf("Execute() returned unexpected error: %v", err)
    }

    if res.Database != "connected" {
        t.Errorf("Execute().Database = %q, want %q", res.Database, "connected")
    }
}
```

**Patrón**: mock manual → inyectar vía constructor → testear solo lógica de negocio, sin DB real.

Ejecutar:

```bash
make test
```

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

apperrors.NewNotFoundError("USER_001", "Usuario no encontrado")
apperrors.NewValidationError("EMAIL_001", "Email inválido")
apperrors.NewUnauthorizedError("AUTH_001", "Token expirado")
apperrors.NewInternalError("DB_001", "Error de conexión")
```

El middleware global los mapea automáticamente al HTTP status code correcto:

| Error | HTTP Status |
|---|---|
| `Validation` | 400 |
| `NotFound` | 404 |
| `Conflict` | 409 |
| `Unauthorized` | 401 |
| `Forbidden` | 403 |
| `Internal` | 500 |

---

## Security Headers

Usa [`github.com/goddtriffin/helmet`](https://github.com/goddtriffin/helmet), el port a Go de HelmetJS.

Se activa automáticamente en todas las rutas vía middleware:

```go
router.Use(middleware.SecurityHeaders())
```

**Cabeceras incluidas por defecto:**

| Header | Valor |
|---|---|
| `X-Frame-Options` | `SAMEORIGIN` |
| `X-Content-Type-Options` | `nosniff` |
| `X-XSS-Protection` | `1; mode=block` |
| `X-DNS-Prefetch-Control` | `off` |
| `X-Download-Options` | `noopen` |
| `Strict-Transport-Security` | `max-age=5184000; includeSubDomains` |
| `X-Powered-By` | removido |

---

## Configuración

Las variables de entorno se cargan automáticamente desde `.env.{environment}` (dev/qa) o del sistema (prod).

| Variable | Default | Descripción |
|---|---|---|
| `PORT` | `3200` | Puerto del servidor |
| `ENVIRONMENT` | `dev` | Entorno: `dev`, `qa` o `prod` |
| `LOG_LEVEL` | `info` | Nivel de log |
| `SHUTDOWN_TIMEOUT` | `30` | Timeout de graceful shutdown (seg) |
| `DB_HOST` | `localhost` | Host de PostgreSQL |
| `DB_PORT` | `5432` | Puerto de PostgreSQL |
| `DB_USER` | `postgres` | Usuario de PostgreSQL |
| `DB_PASSWORD` | `postgres` | Contraseña de PostgreSQL |
| `DB_NAME` | `microservice` | Nombre de la base de datos |
| `DB_SSL_MODE` | `disable` | Modo SSL (`disable`, `require`, `verify-full`) |

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
| Interfaz en domain, implementación en infrastructure | `domain.HealthRepository` → `infrastructure.PostgresHealthRepository` |
