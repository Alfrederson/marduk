package tabuas

import "time"

type Saldo struct {
	Total  int64     `json:"total"`
	Data   time.Time `json:"data_extrato"`
	Limite int64     `json:"limite"`
}

type Transacao struct {
	Valor       int64     `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type Extrato struct {
	Saldo              Saldo        `json:"saldo"`
	UltimasTransacaoes []*Transacao `json:"ultimas_transacoes"`
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
