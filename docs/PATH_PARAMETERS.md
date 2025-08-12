# URL Shortener API - Path Parameters

## Cambios Realizados

Se ha actualizado la API para usar **path parameters** en lugar de query parameters para las redirecciones.

## Endpoints

### 1. Acortar URL

- **Método:** `POST`
- **Endpoint:** `/shorten`
- **Content-Type:** `application/json`

**Request Body:**

```json
{
  "url": "https://example.com/very/long/url"
}
```

**Response:**

```json
{
  "original_url": "https://example.com/very/long/url",
  "shortened_id": "abc12345",
  "shortened_url": "http://localhost:8080/abc12345",
  "message": "URL shortened successfully"
}
```

### 2. Redirección (NUEVO - Path Parameter)

- **Método:** `GET`
- **Endpoint:** `/{shortID}`
- **Ejemplo:** `http://localhost:8080/abc12345`

**Comportamiento:**

- Redirige automáticamente a la URL original con status `302 Found`
- Si el shortID no existe, retorna `404 Not Found`

### 3. Información de la API

- **Método:** `GET`
- **Endpoint:** `/`

**Response:**

```json
{
  "message": "URL Shortener API",
  "usage": "POST /shorten with JSON body {\"url\": \"your-url\"} to shorten, GET /{shortID} to redirect"
}
```

## Cambios en el Código

### Antes (Query Parameter):

```
GET /?id=abc12345
```

### Después (Path Parameter):

```
GET /abc12345
```

### Ventajas del Path Parameter:

1. **URLs más limpias:** `example.com/abc123` vs `example.com/?id=abc123`
2. **SEO friendly:** Los path parameters son mejor interpretados por motores de búsqueda
3. **Estándar de la industria:** La mayoría de URL shorteners usan este patrón
4. **Menos verboso:** URLs más cortas y fáciles de compartir
5. **Mejor UX:** URLs más intuitivas para los usuarios

## Ejemplos de Uso

### Usando curl:

```bash
# Acortar URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'

# Respuesta:
# {
#   "original_url": "https://www.google.com",
#   "shortened_id": "xyz789",
#   "shortened_url": "http://localhost:8080/xyz789",
#   "message": "URL shortened successfully"
# }

# Acceder a URL acortada (redirige automáticamente)
curl -L http://localhost:8080/xyz789
```

### Usando el ejemplo de Go:

```bash
# Ejecutar el servidor
go run ./cmd

# En otra terminal, ejecutar el ejemplo
go run examples/api_usage.go
```

## Validaciones y Errores

- **URL inválida:** Retorna `400 Bad Request` con mensaje de error
- **ShortID no encontrado:** Retorna `404 Not Found`
- **Método no permitido:** Retorna `405 Method Not Allowed`
- **Body JSON inválido:** Retorna `400 Bad Request`

## Consideraciones de Seguridad

- Los path parameters son naturalmente más seguros contra SQL injection que query parameters
- Se mantienen todas las validaciones de URL existentes
- El shortID se extrae directamente del path usando `r.URL.Path[1:]`
