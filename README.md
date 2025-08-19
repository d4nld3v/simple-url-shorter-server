# 🔗 URL Shortener Go

Un servicio de acortamiento de URLs rápido y seguro construido con Go, diseñado con arquitectura limpia y mejores prácticas de seguridad.

## ✨ Características

- 🚀 **Alto rendimiento** - Construido con Go nativo sin frameworks pesados
- 🛡️ **Seguridad robusta** - Protección contra SSRF, XSS, inyecciones SQL y ataques de fuerza bruta
- 📊 **Rate limiting** - Control de tasa de requests por IP con limpieza automática de memoria
- 🔒 **Validación estricta** - Validación completa de URLs y prevención de acceso a redes privadas
- 📈 **Tracking de clicks** - Contador de clicks por URL acortada
- 🗄️ **Base de datos SQLite** - Almacenamiento persistente con connection pooling
- 🔐 **Headers de seguridad** - Headers HTTP de seguridad implementados
- 📝 **Respuestas estandarizadas** - API con formato JSON consistente y códigos de error claros

## 🏗️ Arquitectura

El proyecto sigue los principios de **Clean Architecture**:

```
cmd/api/v1/              # Punto de entrada de la aplicación
├── main.go              # Configuración e inicio del servidor

internal/                # Código interno de la aplicación
├── config/              # Configuración y conexión a base de datos
│   ├── config.go        # Carga de configuración
│   └── db.go            # Patrón Singleton para DB con connection pooling
├── handler/             # Controladores HTTP (Presentation Layer)
│   ├── url.go           # Handlers para endpoints de URLs
│   └── response.go      # Estructuras y helpers para respuestas
├── services/            # Lógica de negocio (Business Layer)
│   ├── url.go           # Servicios de acortamiento de URLs
│   └── validator.go     # Validaciones de URLs y seguridad
├── repository/          # Acceso a datos (Data Layer)
│   └── url.go           # Operaciones CRUD con base de datos
└── server/              # Configuración del servidor HTTP
    └── server.go        # Setup de rutas y middlewares

pkg/middleware/          # Middlewares reutilizables
├── api.go               # Headers de seguridad
└── ratelimit.go         # Rate limiting con token bucket algorithm

docs/                    # Documentación
└── PATH_PARAMETERS.md   # Documentación de parámetros
```

## 🔒 Características de Seguridad

### **Protección contra ataques:**

- **SSRF (Server-Side Request Forgery)** - Bloqueo de IPs privadas y metadata servers
- **SQL Injection** - Prepared statements con placeholders
- **XSS** - Validación y sanitización de entrada
- **Rate Limiting** - Límite de requests por IP con algoritmo token bucket
- **Hash Collision** - SHA-256 + timestamp + random bytes en lugar de MD5

### **Validaciones implementadas:**

- URLs malformadas o peligrosas
- Acceso a redes privadas (127.0.0.1, 192.168.x.x, etc.)
- Caracteres de control y potencialmente peligrosos
- Límites de longitud y formato de datos
- Timeouts para prevenir ataques DoS

## 🚀 Inicio Rápido

### **Prerrequisitos**

- Go 1.19 o superior
- SQLite (incluido en el proyecto)

### **Instalación**

```bash
# Clonar el repositorio
git clone https://github.com/tu-usuario/url-shortener-go.git
cd url-shortener-go

# Instalar dependencias
go mod download

# Ejecutar la aplicación
go run cmd/api/v1/main.go
```

El servidor estará disponible en `http://localhost:8080`

## 📚 API Reference

### **Acortar URL**

```http
POST /shorten
Content-Type: application/json

{
  "url": "https://example.com/very/long/url"
}
```

**Respuesta exitosa (201 Created):**

```json
{
  "data": {
    "original_url": "https://example.com/very/long/url",
    "shortened_id": "abc123xy",
    "shortened_url": "http://localhost:8080/short/abc123xy",
    "created_at": "2025-08-19T10:30:00Z"
  },
  "message": "URL shortened successfully",
  "timestamp": "2025-08-19T10:30:00Z"
}
```

### **Redireccionar URL**

```http
GET /short/{shortID}
```

**Respuesta:** Redirección HTTP 302 a la URL original

### **Respuestas de Error**

Todas las respuestas de error siguen el formato estándar:

```json
{
  "error": {
    "code": "INVALID_URL",
    "message": "Invalid URL",
    "details": "only http and https schemes are allowed"
  },
  "timestamp": "2025-08-19T10:30:00Z",
  "path": "/shorten",
  "method": "POST"
}
```

### **Códigos de Error**

| Código                | Descripción                        |
| --------------------- | ---------------------------------- |
| `INVALID_INPUT`       | Datos de entrada inválidos         |
| `INVALID_URL`         | URL con formato incorrecto         |
| `URL_TOO_LONG`        | URL excede el límite de caracteres |
| `URL_NOT_FOUND`       | URL corta no encontrada            |
| `RATE_LIMIT_EXCEEDED` | Límite de requests excedido        |
| `METHOD_NOT_ALLOWED`  | Método HTTP no permitido           |
| `INTERNAL_ERROR`      | Error interno del servidor         |

## ⚙️ Configuración

El servicio se configura mediante variables de entorno:

```bash
# Puerto del servidor
ADDR=:8080

# Rate limiting
RATE_LIMIT=100          # Requests por minuto por IP
BURST_LIMIT=10          # Burst máximo permitido

# Base de datos
DB_PATH=./url_shortener.db
```

## 🧪 Testing

```bash
# Ejecutar tests
go test ./...

# Tests con coverage
go test -cover ./...

# Tests de integración
go test -tags=integration ./...
```

## 📊 Ejemplos de Uso

### **Acortar una URL**

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

### **Usar la URL acortada**

```bash
curl -L http://localhost:8080/short/abc123xy
```

## 🔧 Desarrollo

### **Estructura de commits**

- `feat:` nuevas características
- `fix:` corrección de bugs
- `docs:` cambios en documentación
- `refactor:` refactorización de código
- `test:` añadir o modificar tests

### **Contribuir**

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/amazing-feature`)
3. Commit tus cambios (`git commit -m 'feat: add amazing feature'`)
4. Push a la rama (`git push origin feature/amazing-feature`)
5. Abre un Pull Request

## 📈 Rendimiento

- **Tiempo de respuesta**: < 50ms para acortamiento
- **Redirecciones**: < 10ms
- **Rate limit**: 100 requests/min por IP por defecto
- **Base de datos**: Connection pooling con máximo 25 conexiones

## 🛠️ Tecnologías Utilizadas

- **[Go](https://golang.org/)** - Lenguaje de programación
- **[SQLite](https://sqlite.org/)** - Base de datos embebida
- **[net/http](https://pkg.go.dev/net/http)** - Servidor HTTP nativo
- **Arquitectura limpia** - Separación de capas y responsabilidades

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.

## 👥 Autores

- **Tu Nombre** - [@d4nld3v](https://github.com/d4nld3v)

## 🙏 Agradecimientos

- Inspirado en las mejores prácticas de seguridad de OWASP
- Basado en principios de Clean Architecture
- Implementación de patrones de diseño estándar de Go
