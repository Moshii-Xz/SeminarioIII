# Guía de Despliegue con Docker (VenciTrack)

Esta guía explica cómo levantar todo el entorno de desarrollo (Backend + Base de Datos) utilizando Docker. Esta es la forma más recomendada para asegurar que todos trabajen en el mismo entorno.

## 1. Prerrequisitos

*   **Docker Desktop**: Asegúrate de tenerlo instalado y corriendo. [Descargar aquí](https://www.docker.com/products/docker-desktop/).

## 2. Iniciar la Aplicación

Abre una terminal en la carpeta raíz del proyecto y ejecuta:

```bash
docker-compose up --build
```

Este comando hará lo siguiente automáticamente:
1.  Compilará el código de Go en un contenedor aislado.
2.  Descargará e iniciará una base de datos PostgreSQL.
3.  Creará la base de datos `vencitrack` y ejecutará automáticamente el script de migración (`migrations/AS_BD.sql`) para crear las tablas.
4.  Conectará el backend con la base de datos.
5.  Expondrá la API en `http://localhost:8081`.

## 3. Verificar que funciona

Una vez que veas logs indicando que el servidor ha iniciado, abre tu navegador o Postman y visita:

`http://localhost:8081/health`

Deberías recibir:
```json
{
    "message": "dabase is ok",
    "status": "healthy"
}
```

## 4. Comandos Útiles

**Detener los contenedores:**
Presiona `Ctrl + C` en la terminal donde corre el log, o ejecuta en otra terminal:
```bash
docker-compose down
```

**Limpiar todo (incluyendo datos de la BD):**
Si quieres reiniciar la base de datos desde cero (borrar todos los datos):
```bash
docker-compose down -v
```

**Ver logs en segundo plano:**
Si prefieres correrlo en modo "detached" (sin ocupar la terminal):
```bash
docker-compose up -d
```
Y para ver los logs después:
```bash
docker-compose logs -f
```

## 5. Solución de Problemas

**Error: "port is already allocated"**
Si te dice que el puerto 5432 o 8080 está ocupado:
1.  Asegúrate de no tener otro Postgres local corriendo en tu máquina.
2.  Asegúrate de no tener otra instancia de la app corriendo (ej. con `go run`).

**Error de conexión a BD al iniciar**
A veces el backend inicia más rápido que la base de datos. Docker Compose intenta manejar esto, pero si falla la primera vez, espera unos segundos; el contenedor de la app debería reiniciarse automáticamente e intentar conectar de nuevo.
