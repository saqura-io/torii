package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/saqura-io/torii/internal/auth"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Route struct {
	Path    string `yaml:"path"`
	Service string `yaml:"service"`
}

type RouteConfig struct {
	Routes []Route `yaml:"routes"`
}

func LoadConfig(filename string) (*RouteConfig, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var config RouteConfig
	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

/*
func SetupRoutes(app *fiber.App, config *RouteConfig) {
	for _, route := range config.Routes {
		r := route // Create a new 'r' variable to hold the 'route'

		// Adding the authorization middleware
		app.Use(r.Path, auth.New())

		// List of all HTTP methods
		methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "CONNECT", "TRACE"}

		// Custom handler for each route and method
		for _, method := range methods {
			app.Add(method, r.Path+"/*", func(c *fiber.Ctx) error {
				// Strip the prefix from the path
				strippedPath := strings.TrimPrefix(c.Path(), r.Path)

				// Create a new URL for the request
				newURL, err := url.Parse(r.Service + strippedPath)
				if err != nil {
					return c.Status(http.StatusInternalServerError).SendString("Internal server error")
				}

				// Create a new HTTP request
				newReq, err := http.NewRequest(c.Method(), newURL.String(), bytes.NewReader(c.Body()))
				if err != nil {
					return c.Status(http.StatusInternalServerError).SendString("Internal server error")
				}

				// Copy the original headers
				c.Request().Header.VisitAll(func(k []byte, v []byte) {
					newReq.Header.Add(string(k), string(v))
				})

				// Send the new HTTP request
				res, err := http.DefaultClient.Do(newReq)
				if err != nil {
					return c.Status(http.StatusInternalServerError).SendString("Internal server error")
				}
				defer res.Body.Close()

				// Relay the response back
				body, _ := io.ReadAll(res.Body)
				for k, v := range res.Header {
					c.Response().Header.Set(k, strings.Join(v, ","))
				}
				c.Response().Header.SetStatusCode(res.StatusCode)
				c.Response().SetBody(body)

				return nil
			})
		}
	}
}*/

func SetupRoutes(app *fiber.App, config *RouteConfig) {
	for _, route := range config.Routes {
		r := route

		// Adding the authorization middleware
		app.Use(r.Path, auth.New())

		// Proxy the request to the internal service
		app.All(r.Path+"/*", func(c *fiber.Ctx) error {
			// Strip the prefix from the path
			c.Path(strings.TrimPrefix(c.Path(), r.Path))

			// Call the proxy middleware
			return proxy.Balancer(proxy.Config{
				Servers: []string{r.Service},
			})(c)
		})
	}
}
