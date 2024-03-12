package devoção

import (
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/Alfrederson/crebitos/escriba"
	"github.com/Alfrederson/crebitos/tabuas"
)

type Devoto struct {
	locks map[string]*sync.Mutex
}

func (c *Devoto) RegistrarTransacao(clientId string, tipo string, valor int64, motivo string) (tabuas.Saldo, error) {
	resultado := tabuas.Saldo{}
	err := escriba.Cunhar(
		[]string{
			filepath.Join("tabuas", clientId, "saldo"),
			filepath.Join("tabuas", clientId, "extrato"),
		},
		[]int{
			escriba.ApagarTudo,
			escriba.BotarNoFinal,
		},
		func(t []*escriba.Tabua) error {
			s := escriba.Conta{}
			escriba.LerContaDaTabua(t[0], &s)
			if tipo == "c" {
				s.Saldo += valor
			} else if tipo == "d" {
				if (s.Saldo + s.Limite - valor) < 0 {
					return tabuas.ErroTransacao{Codigo: tabuas.E_LIMITE_INSUFICIENTE}
				}
				s.Saldo -= valor
			}
			resultado.Total = s.Saldo
			resultado.Limite = s.Limite

			escriba.EscreverContaNaTabua(t[0], &s)
			escriba.AnotarTransacaoNaTabua(t[1], tipo, valor, motivo)
			return nil
		},
	)
	return resultado, err
}

func (c *Devoto) ConsultarExtrato(clientId string, num int) (*tabuas.Extrato, error) {
	resultado := tabuas.Extrato{
		Saldo: tabuas.Saldo{
			Data: time.Now(),
		},
	}
	err := escriba.Cunhar(
		[]string{
			filepath.Join("tabuas", clientId, "saldo"),
			filepath.Join("tabuas", clientId, "extrato"),
		},
		[]int{
			escriba.ApenasLer,
			escriba.ApenasLer,
		},
		func(t []*escriba.Tabua) error {
			s := escriba.Conta{}
			escriba.LerContaDaTabua(t[0], &s)
			resultado.Saldo.Total = s.Saldo
			resultado.Saldo.Limite = s.Limite

			transacoes, err := escriba.LerUltimasTransacoesDaTabua(t[1], num)
			if err != nil {
				return err
			}
			resultado.UltimasTransacaoes = transacoes
			return nil
		},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &resultado, nil
}
