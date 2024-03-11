package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Alfrederson/crebitos/escriba"
	"github.com/Alfrederson/crebitos/templo"
)

func TestCunhagem(t *testing.T) {
	escriba.Cunhar([]string{
		"saldo.dat",
		"extrato.dat",
	}, []int{
		escriba.ApagarTudo,
		escriba.BotarNoFinal,
	}, func(tabuas []*escriba.Tabua) error {
		saldo := tabuas[0]
		extrato := tabuas[1]

		atual := saldo.LerI64()
		ultimo := saldo.LerMomento()

		atual -= 10
		t.Logf("saldo no dia %v: %d", ultimo, atual)
		saldo.Seek(0, 0)
		saldo.CunharI64(atual)
		saldo.CunharMomento(time.Now())

		extrato.CunharU8('C')
		extrato.CunharU64(1)

		extrato.WriteString("comi um cachorro")

		return nil
	})
}

func TestExtrato(t *testing.T) {
	templo.Anotar(
		"1",
		&templo.Transacao{
			Valor:     300,
			Tipo:      templo.TIPO_CREDITO,
			Descricao: "vendeu um ostromelo",
		},
	)
	extrato, _ := templo.ConsultarExtrato("1", 10)
	j, err := json.MarshalIndent(extrato, "", " ")
	t.Log(string(j), err)
}

func TestIniciar(t *testing.T) {
	templo.Inicializar()
}

func TestTimeZones(t *testing.T){
	extrato, _ := templo.ConsultarExtrato("2",5)
	for k,v := range(extrato.UltimasTransacaoes){
		//j, err := json.MarshalIndent(v, "", " ")
		t.Log(k,v.RealizadaEm)	
	}
}