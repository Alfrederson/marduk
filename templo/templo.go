package templo

// o templo é onde todas as dívidas e créditos são registrados na suméria antiga
// fazendo desse jeito tosco porque mandar todas as requisições direto pro sacerdote
// memorizar e ir anotando na argila seria um cheatcode

import (
	"log"
	"os"
	"path/filepath"
	"time"

	devoção "github.com/Alfrederson/crebitos/devocao"
	"github.com/Alfrederson/crebitos/escriba"
	"github.com/Alfrederson/crebitos/tabuas"
)

type Conta struct {
	Limite int64
	Saldo  int64
}

const (
	E_CLIENTE_DESCONHECIDO      = 1
	E_LIMITE_INSUFICIENTE       = 2
	E_SACERDOTE_FOI_AO_BANHEIRO = 3
)

type ErroTransacao struct {
	Codigo int
}

func (e ErroTransacao) Error() string {
	switch e.Codigo {
	case E_CLIENTE_DESCONHECIDO:
		return "cliente desconhecido na aldeia"
	case E_LIMITE_INSUFICIENTE:
		return "cliente sem crédito"
	case E_SACERDOTE_FOI_AO_BANHEIRO:
		return "sacerdote foi ao banheiro"
	default:
		return "não sei o que aconteceu"
	}
}

func InterpretarSonho(e error) *ErroTransacao {
	erro, ok := e.(ErroTransacao)
	if ok {
		return &erro
	}
	return &ErroTransacao{}
}

type ResultadoTransacao struct {
	Limite int64 `json:"limite"`
	Saldo  int64 `json:"saldo"`
}

// isso não era pra estar aqui, mas whateva
var contas map[string]*Conta = map[string]*Conta{
	"1": {
		Limite: 100000,
	},
	"2": {
		Limite: 80000,
	},
	"3": {
		Limite: 1000000,
	},
	"4": {
		Limite: 10000000,
	},
	"5": {
		Limite: 500000,
	},
}

func Inicializar() {
	_, err := os.Stat(filepath.Join("tabuas", "zigurat"))
	if os.IsNotExist(err) {
		log.Println("o escriba prepara as tábuas de argila")
		file, err := os.Create(filepath.Join("tabuas", "zigurat"))
		if err != nil {
			panic(err)
		}
		defer file.Close()
	} else if err != nil {
		panic(err)
	} else {
		log.Println("o escriba percebe que as tábuas de argila estão preparadas")
		return
	}

	for clienteId, valor := range contas {
		dirPath := filepath.Join("tabuas", clienteId)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				panic(err)
			}
		}
		os.Remove(filepath.Join("tabuas", clienteId, "saldo"))
		os.Remove(filepath.Join("tabuas", clienteId, "extrato"))

		escriba.Cunhar([]string{
			filepath.Join("tabuas", clienteId, "saldo"),
			filepath.Join("tabuas", clienteId, "extrato"),
		}, []int{
			escriba.ApagarTudo,
			escriba.ApagarTudo,
		}, func(tabuas []*escriba.Tabua) error {
			conta := tabuas[0]
			conta.CunharI64(valor.Limite)
			conta.CunharI64(0)
			tabuas[1].CunharU8(0)
			return nil
		})
	}
}

func Anotar(devoto *devoção.Devoto, clienteId string, t *tabuas.Transacao) (*ResultadoTransacao, error) {
	_, tem := contas[clienteId]
	if !tem {
		return nil, ErroTransacao{E_CLIENTE_DESCONHECIDO}
	}

	var saldo int64
	var limite int64
	// leia o registro dos bens e os limites do habitante da aldeia
	resposta, err := devoto.ConsultarSaldo(clienteId)
	if err != nil {
		return nil, err
	}
	saldo = resposta.Saldo
	limite = resposta.Limite
	switch t.Tipo {
	case TIPO_CREDITO:
		saldo += t.Valor
	case TIPO_DEBITO:
		if saldo+limite-t.Valor < 0 {
			return nil, ErroTransacao{E_LIMITE_INSUFICIENTE}
		}
		saldo -= t.Valor
	}
	devoto.MudarSaldo(clienteId, saldo, t.Descricao)

	if err != nil {
		return nil, err
	}
	return &ResultadoTransacao{
		Limite: int64(limite),
		Saldo:  int64(saldo),
	}, nil
}

type Saldo struct {
	Total  int64     `json:"total"`
	Data   time.Time `json:"data_extrato"`
	Limite int64     `json:"limite"`
}

type Extrato struct {
	Saldo              Saldo               `json:"saldo"`
	UltimasTransacaoes []*tabuas.Transacao `json:"ultimas_transacoes"`
}

func ConsultarExtrato(devoto *devoção.Devoto, clienteId string, num int) (*Extrato, error) {
	_, tem := contas[clienteId]
	if !tem {
		return nil, ErroTransacao{E_CLIENTE_DESCONHECIDO}
	}
	saldo, err := devoto.ConsultarSaldo(clienteId)
	if err != nil {
		return nil, ErroTransacao{E_SACERDOTE_FOI_AO_BANHEIRO}
	}
	resultado := &Extrato{
		Saldo: Saldo{
			Total:  saldo.Saldo,
			Limite: saldo.Limite,
			Data:   time.Now(),
		},
		UltimasTransacaoes: make([]*tabuas.Transacao, 0),
	}

	ultimas, err := devoto.UltimasTransacaoes(clienteId, num)
	if err != nil {
		return nil, err
	}
	resultado.UltimasTransacaoes = ultimas

	return resultado, nil

}
