// a função da viga é distribuir a carga em 2 pilares.
package main

import (
	"log"
	"os"

	"github.com/Alfrederson/crebitos/templo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func main() {
	templo.Inicializar()
	pilares := os.Args[1:]
	if len(pilares) < 1 {
		panic("a viga sem pilares não tem apoio")
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
