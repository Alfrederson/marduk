package escriba

import (
	"os"
	"path/filepath"
	"time"

	"github.com/Alfrederson/crebitos/tabuas"
)

type Conta struct {
	Saldo  int64
	Limite int64
}

type Tabua struct {
	*os.File
}

func (t *Tabua) CunharU8(n uint8) {
	t.Write([]byte{
		n,
	})
}
func (t *Tabua) LerU8() uint8 {
	buf := make([]byte, 1)
	t.Read(buf)
	return buf[0]
}

func (t *Tabua) CunharU64(n uint64) {
	t.Write([]byte{
		uint8(n >> 56 & 0xFF),
		uint8(n >> 48 & 0xFF),
		uint8(n >> 40 & 0xFF),
		uint8(n >> 32 & 0xFF),
		uint8(n >> 24 & 0xFF),
		uint8(n >> 16 & 0xFF),
		uint8(n >> 8 & 0xFF),
		uint8(n & 0xFF),
	})
}
func (t *Tabua) LerU64() uint64 {
	b := make([]byte, 8)
	t.Read(b)
	return uint64(b[0])<<56 |
		uint64(b[1])<<48 |
		uint64(b[2])<<40 |
		uint64(b[3])<<32 |
		uint64(b[4])<<24 |
		uint64(b[5])<<16 |
		uint64(b[6])<<8 |
		uint64(b[7])
}

func (t *Tabua) CunharI64(n int64) {
	t.CunharU64(uint64(n))
}

func (t *Tabua) LerI64() int64 {
	return int64(t.LerU64())
}

func (t *Tabua) CunharMomento(m time.Time) {
	agora, _ := m.MarshalBinary()
	t.Write(agora)
}
func (t *Tabua) LerMomento() time.Time {
	escrito := make([]byte, 15)
	t.Read(escrito)
	momento := time.Time{}
	momento.UnmarshalBinary(escrito)
	return momento
}

func (t *Tabua) CunharTexto(str string) {
	if len(str) > 255 {
		panic("string comprida demais")
	}
	t.CunharU8(uint8(len(str)))
	t.File.WriteString(str)
}
func (t *Tabua) LerTexto() string {
	tamanho := t.LerU8()
	buf := make([]byte, tamanho)
	t.File.Read(buf)
	return string(buf)
}

func (t *Tabua) Posicao() int64 {
	r, _ := t.Seek(0, 1)
	return r
}

type cunhagem func([]*Tabua) error

const (
	ApenasLer    = os.O_RDONLY
	ApagarTudo   = os.O_CREATE | os.O_RDWR
	BotarNoFinal = os.O_CREATE | os.O_RDWR | os.O_APPEND
)

// Efetua operações de cunhagem em uma lista de tábuas de argila nomeadas
func Cunhar(tabuas []string, modos []int, op cunhagem) error {
	abertos := make([]*Tabua, len(tabuas))
	for i, arquivo := range tabuas {
		esseArquivo, err := os.OpenFile(
			arquivo, modos[i], 0644,
		)
		if err != nil {
			return err
		}
		defer esseArquivo.Close()
		abertos[i] = &Tabua{esseArquivo}
	}
	return op(abertos)
}

func LerConta(clienteId string) *Conta {
	c := &Conta{}
	Cunhar(
		[]string{
			filepath.Join("tabuas", clienteId, "saldo"),
		},
		[]int{
			ApenasLer,
		}, func(t []*Tabua) error {
			c.Limite = t[0].LerI64()
			c.Saldo = t[0].LerI64()
			return nil
		})
	return c
}

func MudarConta(clienteId string, c *Conta) {
	Cunhar(
		[]string{
			filepath.Join("tabuas", clienteId, "saldo"),
		},
		[]int{
			ApagarTudo,
		}, func(t []*Tabua) error {
			t[0].Seek(0, 0)
			t[0].CunharI64(c.Limite)
			t[0].CunharI64(c.Saldo)
			return nil
		},
	)
}

func AnotarTransacao(clienteId string, tipo string, valor int64, descricao string) {
	Cunhar(
		[]string{
			filepath.Join("tabuas", clienteId, "extrato"),
		},
		[]int{
			BotarNoFinal,
		},
		func(t []*Tabua) error {
			tabua := t[0]

			tabua.CunharU8(tipo[0])
			tabua.CunharI64(valor)
			tabua.CunharMomento(time.Now())
			tabua.CunharTexto(descricao)
			tabua.CunharU8(uint8(1 + 8 + 15 + 1 + len(descricao)))
			return nil
		},
	)
}

func lerTransacaoDaTabua(tabua *Tabua) *tabuas.Transacao {
	return &tabuas.Transacao{
		Tipo:        string(tabua.LerU8()),
		Valor:       tabua.LerI64(),
		RealizadaEm: tabua.LerMomento(),
		Descricao:   tabua.LerTexto(),
	}
}

func LerUltimasTransacoes(clienteId string, quantas int) ([]*tabuas.Transacao, error) {
	resultado := make([]*tabuas.Transacao, 0, quantas)
	err := Cunhar([]string{
		filepath.Join("tabuas", clienteId, "extrato"),
	}, []int{
		ApenasLer,
	}, func(t []*Tabua) error {
		extrato := t[0]
		extrato.Seek(-1, 2)
		comprimento := int64(extrato.LerU8())
		if comprimento == 0 {
			return nil
		}
		for i := 0; i < quantas; i++ {
			extrato.Seek(-comprimento-1, 1)
			resultado = append(resultado, lerTransacaoDaTabua(extrato))
			// o registro anterior não existe.
			if extrato.Posicao()-comprimento-1 <= 0 {
				break
			}
			extrato.Seek(-comprimento-1, 1)
			comprimento = int64(extrato.LerU8())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resultado, nil
}
