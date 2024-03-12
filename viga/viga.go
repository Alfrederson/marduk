package viga

import (
	"log"

	"github.com/Alfrederson/crebitos/templo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func Viga(pilares []string) {
	templo.Inicializar()
	if len(pilares) < 1 {
		panic("a viga precisa de pelo menos dois pilares, do contrário cai.")
	}
	for i, v := range pilares {
		log.Println("pilar", i, "=", v)
	}
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	// app.Use(logger.New(logger.Config{
	// 	Output: os.Stdout,
	// }))
	// app.Use(func(c *fiber.Ctx) error {
	// 	start := time.Now()
	// 	c.Next()
	// 	duration := time.Since(start)
	// 	log.Printf("%s %s %v", c.Method(), c.Path(), duration)
	// 	return nil
	// })
	app.Use(proxy.Balancer(proxy.Config{
		Servers: pilares,
	}))
	log.Fatal(app.Listen("0.0.0.0:9999"))
}
