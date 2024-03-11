package tabuas

import "time"

type Saldo struct {
	Limite int64
	Saldo  int64
}

type Transacao struct {
	Valor       int64     `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}
