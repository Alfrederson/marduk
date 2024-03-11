package templo

// o templo é onde todas as dívidas e créditos são registrados na suméria antiga
// fazendo desse jeito tosco porque mandar todas as requisições direto pro sacerdote
// memorizar e ir anotando na argila seria um cheatcode

import (
	"log"
	"os"
	"path/filepath"

	devocao "github.com/Alfrederson/crebitos/devocao"
	"github.com/Alfrederson/crebitos/escriba"
	"github.com/Alfrederson/crebitos/tabuas"
)

type Conta struct {
	Limite int64
	Saldo  int64
}

func InterpretarSonho(e error) *tabuas.ErroTransacao {
	erro, ok := e.(tabuas.ErroTransacao)
	if ok {
		return &erro
	}
	return &tabuas.ErroTransacao{}
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

func Anotar(devoto *devocao.Devoto, clienteId string, t *tabuas.Transacao) (*ResultadoTransacao, error) {
	_, tem := contas[clienteId]
	if !tem {
		return nil, tabuas.ErroTransacao{Codigo: tabuas.E_CLIENTE_DESCONHECIDO}
	}

	// leia o registro dos bens e os limites do habitante da aldeia
	resposta, err := devoto.RegistrarTransacao(
		clienteId,
		t.Tipo,
		t.Valor,
		t.Descricao,
	)
	if err != nil {
		return nil, err
	}

	return &ResultadoTransacao{
		Limite: int64(resposta.Limite),
		Saldo:  int64(resposta.Total),
	}, nil
}

func ConsultarExtrato(devoto *devocao.Devoto, clienteId string, num int) (*tabuas.Extrato, error) {
	_, tem := contas[clienteId]
	if !tem {
		return nil, tabuas.ErroTransacao{Codigo: tabuas.E_CLIENTE_DESCONHECIDO}
	}

	extrato, err := devoto.ConsultarExtrato(clienteId, num)
	if err != nil {
		return nil, tabuas.ErroTransacao{Codigo: tabuas.E_SACERDOTE_FOI_AO_BANHEIRO}
	}

	return extrato, nil
}
