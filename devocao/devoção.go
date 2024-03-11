package devoção

import (
	"context"
	"errors"
	"sync"
	"time"

	pb "github.com/Alfrederson/crebitos/proto"
	"github.com/Alfrederson/crebitos/tabuas"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client duplica os locks do server porque essa instancia pode
// ter mais de uma requisição tentando alterar o mesmo arquivo,
// então a gente tranca "aqui" antes de mandar a requisição
// pro grpc pra lotar a fiação dos containers. não sei se isso
// faz sentido, pra mim não parece fazer, mas o que é isso
// perto de um postgres da vida?
type Devoto struct {
	pb.SacerdoteClient
	locks map[string]*sync.Mutex
	lock  sync.Mutex
	conn  *grpc.ClientConn
}

func (c *Devoto) Start(address string) error {
	var err error
	c.conn, err = grpc.Dial(
		address,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return err
	}
	c.locks = make(map[string]*sync.Mutex)
	c.SacerdoteClient = pb.NewSacerdoteClient(c.conn)
	return nil
}

// faz alguma coisa exclusivamente em um arquivo.
// adaptado da implementação em COBOL.
func (c *Devoto) LockFile(filename string, exclusive func() error) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	mut, tem := c.locks[filename]
	if !tem {
		mut = &sync.Mutex{}
		c.locks[filename] = mut
	}
	mut.Lock()
	defer mut.Unlock()

	res, err := c.SacerdoteClient.Consultar(context.Background(), &pb.Habitante{Id: filename})
	// se a tranca tiver pego, então a gente manda soltar a tranca.
	if err == nil && !res.Error {
		defer func() {
			res, err := c.SacerdoteClient.Sair(context.Background(), &pb.Habitante{Id: filename})
			if res.Error {
				panic("isso nunca vai acontecer.")
			}
			if err != nil {
				panic(err)
			}
		}()
	}
	if err != nil {
		return err
	}
	if res.Error {
		return errors.New("o sacerdote disse algo incompreensível")
	}
	return exclusive()
}

// func (c *Devoto) ConsultarSaldo(clienteId string) (tabuas.Saldo, error) {
// 	r, err := c.SacerdoteClient.ConsultarSaldo(context.Background(), &pb.SaldoConsulta{ClienteId: clienteId})
// 	if err != nil {
// 		return tabuas.Saldo{}, err
// 	}
// 	return tabuas.Saldo{
// 		Limite: r.Limite,
// 		Total:  r.Saldo,
// 	}, nil
// }

func (c *Devoto) ConsultarExtrato(clienteId string, num int) (*tabuas.Extrato, error) {
	r, err := c.SacerdoteClient.ConsultarExtrato(context.Background(), &pb.Habitante{Id: clienteId})
	if err != nil {
		return nil, err
	}
	resultado := tabuas.Extrato{
		Saldo: tabuas.Saldo{
			Limite: r.Limite,
			Total:  r.Saldo,
			Data:   time.Now(),
		},
		UltimasTransacaoes: make([]*tabuas.Transacao, 0, num),
	}
	for _, v := range r.UltimasTransacoes {
		resultado.UltimasTransacaoes = append(resultado.UltimasTransacaoes, &tabuas.Transacao{
			Valor:       v.Valor,
			Descricao:   v.Descricao,
			RealizadaEm: v.RealizadaEm.AsTime(),
			Tipo:        v.Tipo,
		})
	}
	return &resultado, nil
}

func (c *Devoto) RegistrarTransacao(clientId string, tipo string, valor int64, motivo string) (tabuas.Saldo, error) {
	r, err := c.SacerdoteClient.RegistrarTransacao(context.Background(), &pb.PedidoTransacao{
		ClienteId: clientId,
		Tipo:      tipo,
		Descricao: motivo,
		Valor:     valor,
	})

	if err != nil {
		return tabuas.Saldo{}, err
	}
	if r.Erro != 0 {
		return tabuas.Saldo{}, tabuas.ErroTransacao{Codigo: int(r.Erro)}
	}
	return tabuas.Saldo{
		Limite: r.Limite,
		Total:  r.NovoSaldo,
	}, nil
}

func (c *Devoto) Stop() {
	c.conn.Close()
}
