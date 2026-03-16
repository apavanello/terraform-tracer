package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/apavanello/terraform-tracer/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed all:dist
var uiDist embed.FS

// Server holds the parsed graph and serves the API + embedded frontend.
type Server struct {
	app   *fiber.App
	graph *models.Graph
}

// NewServer creates a new Server with the given parsed graph.
func NewServer(graph *models.Graph) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	s := &Server{app: app, graph: graph}

	// CORS for dev mode (Vue dev server on :5173)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// API routes
	api := app.Group("/api/v1")
	api.Get("/graph", s.handleGetGraph)
	api.Get("/health", s.handleHealth)

	// Serve embedded frontend (production build)
	distFS, err := fs.Sub(uiDist, "dist")
	if err != nil {
		log.Printf("warning: embedded UI not found (dev mode?): %v", err)
	} else {
		app.Use("/", filesystem.New(filesystem.Config{
			Root:         http.FS(distFS),
			Browse:       false,
			Index:        "index.html",
			NotFoundFile: "index.html", // SPA fallback
		}))
	}

	return s
}

// Listen starts the HTTP server on the given address.
func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) handleGetGraph(c *fiber.Ctx) error {
	return c.JSON(s.graph)
}

func (s *Server) handleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
