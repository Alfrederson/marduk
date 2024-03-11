// servidor do file locker.

package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Alfrederson/crebitos/escriba"
	pb "github.com/Alfrederson/crebitos/proto"
	"github.com/Alfrederson/crebitos/tabuas"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Extrato = []*tabuas.Transacao

var contas = make(map[string]*escriba.Conta)
var extratos = make(map[string]*tabuas.RingBuffer)

type AlteracaoSaldo struct {
	ContaId     string
	AntigoValor int64
	NovoValor   int64
	Motivo      string
	Limite      int64
}

type sacerdote struct {
	pb.SacerdoteServer
	locks      map[string]*sync.Mutex
	master     sync.Mutex
	fila       sync.Mutex
	alteracoes chan (*AlteracaoSaldo)
}

func (f *sacerdote) init() *sacerdote {
	f.locks = make(map[string]*sync.Mutex, 10)
	f.alteracoes = make(chan *AlteracaoSaldo)

	go func() {
		for alteracao := range f.alteracoes {
			// anota o valor novo
			escriba.MudarConta(
				alteracao.ContaId,
				&escriba.Conta{
					Saldo:  alteracao.NovoValor,
					Limite: alteracao.Limite,
				},
			)
			// anota a transação
			var tipo string
			var delta int64 = alteracao.AntigoValor - alteracao.NovoValor
			if alteracao.AntigoValor > alteracao.NovoValor {
				tipo = "d"
				delta *= -1
			} else {
				tipo = "c"
			}
			escriba.AnotarTransacao(
				alteracao.ContaId,
				tipo,
				delta,
				alteracao.Motivo,
			)
		}
	}()

	return f
}

func (f *sacerdote) ConsultarExtrato(ctx context.Context, req *pb.Habitante) (*pb.Extrato, error) {
	conta, tem := contas[req.Id]
	f.fila.Lock()
	defer f.fila.Unlock()

	if !tem {
		log.Println("sacerdote pede ao escriba para ler uma tábia ", req.Id)
		tabua := escriba.LerConta(req.Id)
		contas[req.Id] = tabua
		conta = tabua
		f.initRingBuffer(req.Id)
	}

	ultimasTransacoes := make([]*pb.Transacao, 0, 10)

	for _, v := range extratos[req.Id].Unroll() {
		t := v.(*tabuas.Transacao)
		ultimasTransacoes = append(ultimasTransacoes, &pb.Transacao{
			Tipo:        t.Tipo,
			Valor:       t.Valor,
			RealizadaEm: timestamppb.New(t.RealizadaEm),
			Descricao:   t.Descricao,
		})
	}

	return &pb.Extrato{
		Saldo:             conta.Saldo,
		Limite:            conta.Limite,
		UltimasTransacoes: ultimasTransacoes,
	}, nil
}

func (f *sacerdote) initRingBuffer(clienteId string) {
	buf := tabuas.NewRingBuffer(10)
	extratos[clienteId] = buf

	ultimasTransacoes, err := escriba.LerUltimasTransacoes(clienteId, 10)
	if err != nil {
		panic(err.Error())
	}
	for _, t := range ultimasTransacoes {
		buf.Add(t)
	}
	log.Println("o escriba leu o extrato de ", clienteId)
}

func (f *sacerdote) initConta(clienteId string) *escriba.Conta {
	tabua := escriba.LerConta(clienteId)
	contas[clienteId] = tabua
	log.Println("o escriba leu a conta. saldo = ", tabua.Saldo, " limite = ", tabua.Limite)
	return tabua
}

func (f *sacerdote) ConsultarSaldo(ctx context.Context, req *pb.SaldoConsulta) (*pb.SaldoCliente, error) {
	conta, tem := contas[req.ClienteId]
	f.fila.Lock()
	defer f.fila.Unlock()

	// o sacerdote diz o número que ele tiver memorizado.
	if !tem {
		conta = f.initConta(req.ClienteId)
		f.initRingBuffer(req.ClienteId)
	}
	return &pb.SaldoCliente{
		Saldo:  conta.Saldo,
		Limite: conta.Limite,
	}, nil
}

func (f *sacerdote) MudarSaldo(ctx context.Context, req *pb.SaldoAtualizacao) (*pb.SaldoCliente, error) {
	conta, tem := contas[req.ClienteId]
	f.fila.Lock()
	defer f.fila.Unlock()

	if !tem {
		log.Println("sacerdote pede ao escriba para ler uma tábua ", req.ClienteId)
		conta = f.initConta(req.ClienteId)
		f.initRingBuffer(req.ClienteId)
	}

	f.alteracoes <- &AlteracaoSaldo{
		ContaId:     req.ClienteId,
		AntigoValor: conta.Saldo,
		NovoValor:   req.NovoSaldo,
		Motivo:      req.Motivo,
		Limite:      conta.Limite,
	}

	var tipo string
	var valor int64
	if req.NovoSaldo > conta.Saldo {
		tipo = "c"
		valor = req.NovoSaldo - conta.Saldo
	} else {
		tipo = "d"
		valor = conta.Saldo - req.NovoSaldo
	}

	extratos[req.ClienteId].Add(&tabuas.Transacao{
		Valor:       valor,
		Tipo:        tipo,
		Descricao:   req.Motivo,
		RealizadaEm: time.Now(),
	})

	conta.Saldo = req.NovoSaldo

	return &pb.SaldoCliente{
		Saldo:  conta.Saldo,
		Limite: conta.Limite,
	}, nil

}

// começa a ler a tábua de argila da pessoa
func (f *sacerdote) Consultar(ctx context.Context, req *pb.Habitante) (*pb.ConsultaStatus, error) {
	f.master.Lock()
	mut, tem := f.locks[req.Id]
	if !tem {
		mut = &sync.Mutex{}
		f.locks[req.Id] = mut
	}
	f.master.Unlock()

	defer mut.Lock()
	return &pb.ConsultaStatus{Error: false}, nil
}

func (f *sacerdote) Sair(ctx context.Context, req *pb.Habitante) (*pb.SaidaStatus, error) {
	f.master.Lock()
	defer f.master.Unlock()
	mut, tem := f.locks[req.Id]
	if !tem {
		return &pb.SaidaStatus{Error: true}, nil
	}

	defer mut.Unlock()
	return &pb.SaidaStatus{Error: false}, nil
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("sacerdote não está no templo porque %v", err)
	}
	s := grpc.NewServer()
	f := sacerdote{}
	f.init()
	pb.RegisterSacerdoteServer(s, &f)
	log.Println("sacerdote escutando na porta 50051")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("sacerdote não pode entrar em contato com marduk porque %v", err)
	}
}
