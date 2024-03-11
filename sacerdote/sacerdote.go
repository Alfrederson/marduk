// servidor do file locker.

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Alfrederson/crebitos/escriba"
	pb "github.com/Alfrederson/crebitos/proto"
	"github.com/Alfrederson/crebitos/tabuas"
	"github.com/Alfrederson/crebitos/templo"
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
	alteracoes chan (*AlteracaoSaldo)
}

func (f *sacerdote) init() *sacerdote {
	f.locks = make(map[string]*sync.Mutex, 10)
	f.alteracoes = make(chan *AlteracaoSaldo)

	templo.Inicializar()
	for i := 1; i < 6; i++ {
		f.initRingBuffer(fmt.Sprintf("%d", i))
		f.initConta(fmt.Sprintf("%d", i))
	}

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
	// tranca a tábua de argila
	m := f.trancar(req.Id)
	defer m.Unlock()
	conta, tem := contas[req.Id]
	if !tem {
		log.Println("sacerdote pede ao escriba para ler uma tábia ", req.Id)
		conta = f.initConta(req.Id)
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

func (f *sacerdote) RegistrarTransacao(ctx context.Context, req *pb.PedidoTransacao) (*pb.ResultadoTransacao, error) {
	m := f.trancar(req.ClienteId)
	defer m.Unlock()
	conta, tem := contas[req.ClienteId]
	if !tem {
		log.Println("sacerdote pede ao escriba para ler uma tábua ", req.ClienteId)
		conta = f.initConta(req.ClienteId)
		f.initRingBuffer(req.ClienteId)
	}
	saldoAnterior := conta.Saldo
	if req.Tipo == "c" {
		conta.Saldo += req.Valor
	}
	if req.Tipo == "d" {
		if (conta.Saldo+conta.Limite)-req.Valor < 0 {
			return &pb.ResultadoTransacao{
				Erro: tabuas.E_LIMITE_INSUFICIENTE,
			}, nil
		}
		conta.Saldo -= req.Valor
	}
	f.alteracoes <- &AlteracaoSaldo{
		ContaId:     req.ClienteId,
		AntigoValor: saldoAnterior,
		NovoValor:   conta.Saldo,
		Motivo:      req.Descricao,
		Limite:      conta.Limite,
	}
	extratos[req.ClienteId].Add(&tabuas.Transacao{
		Valor:       req.Valor,
		Tipo:        req.Tipo,
		Descricao:   req.Descricao,
		RealizadaEm: time.Now(),
	})
	return &pb.ResultadoTransacao{
		NovoSaldo: conta.Saldo,
		Limite:    conta.Limite,
	}, nil
}

func (f *sacerdote) trancar(clienteId string) *sync.Mutex {
	f.master.Lock()
	defer f.master.Unlock()
	mut, tem := f.locks[clienteId]
	if !tem {
		mut = &sync.Mutex{}
		f.locks[clienteId] = mut
	}
	mut.Lock()
	return mut
}

// começa a ler a tábua de argila da pessoa
// func (f *sacerdote) Consultar(ctx context.Context, req *pb.Habitante) (*pb.ConsultaStatus, error) {
// 	f.master.Lock()
// 	mut, tem := f.locks[req.Id]
// 	if !tem {
// 		mut = &sync.Mutex{}
// 		f.locks[req.Id] = mut
// 	}
// 	f.master.Unlock()

// 	defer mut.Lock()
// 	return &pb.ConsultaStatus{Error: false}, nil
// }

// func (f *sacerdote) Sair(ctx context.Context, req *pb.Habitante) (*pb.SaidaStatus, error) {
// 	f.master.Lock()
// 	defer f.master.Unlock()
// 	mut, tem := f.locks[req.Id]
// 	if !tem {
// 		return &pb.SaidaStatus{Error: true}, nil
// 	}

// 	defer mut.Unlock()
// 	return &pb.SaidaStatus{Error: false}, nil
// }

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
