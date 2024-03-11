package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"

	devoção "github.com/Alfrederson/crebitos/devocao"
	"github.com/Alfrederson/crebitos/tabuas"
	"github.com/Alfrederson/crebitos/templo"
)

func main() {
	flockerUrl := os.Args[1]
	port := os.Args[2]

	f := devoção.Devoto{}
	err := f.Start(flockerUrl)
	if err != nil {
		log.Fatalf("não consegui escutar o sacerdote porque %v", err)
	}

	log.Println("estou escutando o sacerdote")
	f.LockFile("zigurat", func() error {
		log.Println("")
		templo.Inicializar()
		return nil
	})
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// app.Use(logger.New(logger.Config{
	// 	Output: os.Stdout,
	// }))

	app.Post("/clientes/:id/transacoes", func(c *fiber.Ctx) error {
		req := tabuas.Transacao{}
		if err := c.BodyParser(&req); err != nil {
			return c.SendStatus(422)
		}
		if req.Valor == 0 {
			return c.SendStatus(422)
		}
		req.Tipo = strings.ToLower(req.Tipo)
		if req.Tipo != "c" && req.Tipo != "d" {
			return c.SendStatus(422)
		}
		if len(req.Descricao) < 1 || len(req.Descricao) > 10 {
			return c.SendStatus(422)
		}
		var resultadoTransacao *templo.ResultadoTransacao
		err := f.LockFile(c.Params("id"), func() error {
			resultadoTransacao, err = templo.Anotar(
				&f,
				c.Params("id"),
				&req,
			)
			return err
		})
		if err != nil {
			e := templo.InterpretarSonho(err)
			switch e.Codigo {
			case templo.E_CLIENTE_DESCONHECIDO:
				return c.Status(404).JSON(e.Error())
			case templo.E_LIMITE_INSUFICIENTE:
				return c.Status(422).JSON(e.Error())
			}
		}
		return c.Status(200).JSON(resultadoTransacao)
	})
	app.Get("/clientes/:id/extrato", func(c *fiber.Ctx) error {
		clientId := c.Params("id")
		var resultado *templo.Extrato
		err := f.LockFile(clientId, func() error {
			if extrato, err := templo.ConsultarExtrato(&f, clientId, 10); err != nil {
				return err
			} else {
				resultado = extrato
				return nil
			}
		})
		if err != nil {
			log.Println(err)
			return c.Status(404).JSON(err)
		}
		return c.Status(200).JSON(resultado)
	})
	app.Get("/reset", func(c *fiber.Ctx) error {
		os.Remove(filepath.Join("tabuas", "zigurat"))
		templo.Inicializar()
		return c.Status(200).SendString("ok")
	})
	log.Fatal(app.Listen("0.0.0.0:" + port))
}
