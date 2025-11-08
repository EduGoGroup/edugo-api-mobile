package router

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/EduGoGroup/edugo-api-mobile/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ConfigureSwaggerHost configura dinámicamente el host de Swagger basado en el puerto del servidor.
// Esto actualiza el SwaggerInfo.Host en tiempo de ejecución para que las peticiones de Swagger UI
// apunten al puerto correcto.
// Si el host es 0.0.0.0, usa localhost para que sea accesible desde el navegador.
func ConfigureSwaggerHost(host string, port int) {
	// 0.0.0.0 no es válido para clientes del navegador, usar localhost
	if host == "0.0.0.0" || host == "" {
		host = "localhost"
	}
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", host, port)
}

// SetupSwaggerUI configura Swagger UI con detección dinámica de host.
// Inyecta JavaScript que detecta window.location.host y lo usa para construir
// las URLs de las peticiones API, eliminando la necesidad de hardcodear el host.
func SetupSwaggerUI(r *gin.Engine) {
	// Servir Swagger UI con handler personalizado que intercepta index.html
	r.GET("/swagger/*any", func(c *gin.Context) {
		// Si la ruta es /swagger/index.html o /swagger/, servir nuestro HTML personalizado
		path := c.Param("any")
		if path == "/index.html" || path == "/" || path == "" {
			c.Header("Content-Type", "text/html; charset=utf-8")
			tmpl := template.Must(template.New("swagger").Parse(swaggerTemplate))
			if err := tmpl.Execute(c.Writer, nil); err != nil {
				c.String(http.StatusInternalServerError, "Error rendering Swagger UI")
				return
			}
			return
		}

		// Para el resto de archivos, usar el handler estándar de Swagger
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})
}

// swaggerTemplate es el HTML personalizado para Swagger UI con detección dinámica de host.
// El JavaScript detecta window.location.host y lo inyecta en la configuración de Swagger,
// permitiendo que las peticiones API usen el host y puerto correctos automáticamente.
const swaggerTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>EduGo API Mobile - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="./swagger-ui.css">
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>

    <script src="./swagger-ui-bundle.js" charset="UTF-8"></script>
    <script src="./swagger-ui-standalone-preset.js" charset="UTF-8"></script>
    <script>
        window.onload = function() {
            // Detectar el host dinámicamente desde la URL del navegador
            const dynamicHost = window.location.host;
            const protocol = window.location.protocol;
            
            console.log('Swagger UI - Detected host:', dynamicHost);
            console.log('Swagger UI - Protocol:', protocol);

            // Configurar Swagger UI
            const ui = SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                // Interceptor para inyectar el host dinámico en todas las peticiones
                requestInterceptor: function(request) {
                    // Si la URL es relativa o no tiene host, usar el host detectado
                    if (request.url.startsWith('/')) {
                        request.url = protocol + '//' + dynamicHost + request.url;
                    }
                    console.log('Swagger UI - Request URL:', request.url);
                    return request;
                },
                onComplete: function() {
                    console.log('Swagger UI - Initialized with dynamic host:', dynamicHost);
                }
            });

            window.ui = ui;
        };
    </script>
</body>
</html>
`
