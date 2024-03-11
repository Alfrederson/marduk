package main

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	devoção "github.com/Alfrederson/crebitos/devocao"
	"github.com/Alfrederson/crebitos/tabuas"
	"github.com/Alfrederson/crebitos/templo"
)

func main() {
	endereçoSacerdote := os.Args[1]
	port := os.Args[2]

	f := devoção.Devoto{}
	err := f.Start(endereçoSacerdote)
	if err != nil {
		log.Fatalf("o devoto conseguiu escutar o sacerdote porque %v", err)
	}
	defer f.Stop()
	log.Println("o devoto está escutando o sacerdote")
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
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
		resultadoTransacao, err := templo.Anotar(
			&f,
			c.Params("id"),
			&req,
		)
		if err != nil {
			e := templo.InterpretarSonho(err)
			switch e.Codigo {
			case tabuas.E_CLIENTE_DESCONHECIDO:
				return c.Status(404).JSON(e.Error())
			case tabuas.E_LIMITE_INSUFICIENTE:
				return c.Status(422).JSON(e.Error())
			}
		}
		return c.Status(200).JSON(resultadoTransacao)
	})
	app.Get("/clientes/:id/extrato", func(c *fiber.Ctx) error {
		clientId := c.Params("id")
		if extrato, err := templo.ConsultarExtrato(&f, clientId, 10); err != nil {
			return c.Status(404).JSON(err)
		} else {
			return c.Status(200).JSON(extrato)
		}
	})
	log.Fatal(app.Listen("0.0.0.0:" + port))
}
