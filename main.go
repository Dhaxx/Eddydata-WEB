package main

import (
	"Eddydata-WEB/modules"
	"Eddydata-WEB/modules/compras"
	// "Eddydata-WEB/modules/patrimonio"
	"sync"

	"github.com/vbauerster/mpb"
)

func main() {
	var wg sync.WaitGroup

	// COMPRAS
	wg.Add(1)
	pc := mpb.New()

	go func() {
		defer wg.Done()
		modules.LimpaCompras()
		compras.Cadunimedida(pc)
		// compras.GrupoSubgrupo(pc)
		// compras.Cadest(pc)
		// compras.Destino(pc)
		// compras.CentroCusto(pc)

		// modules.LimpaLicitacoes()
		// compras.Cadlic(pc)
		// compras.Prolics(pc)
		// compras.Cadprolic(pc)
		// compras.CadproProposta(pc)
		// compras.Aditamento(pc)

		// compras.Requisicoes(pc)
		// compras.Cadped(pc)
		// compras.CadproSaldoAnt(pc)
		// modules.ContratosAdit(pc)
	}()

	// PATRIMONIO
	wg.Add(1)
	go func() {
		// pp := mpb.New()
		// patrimonio.TipoMov(pp)
		// patrimonio.Cadajuste(pp)
		// patrimonio.Cadbai(pp)
		// patrimonio.Cadsit(pp)
		// patrimonio.Cadtip(pp)
		// patrimonio.Cadpatd(pp)
		// patrimonio.Cadpats(pp)
		// patrimonio.Cadpatg(pp)

		// patrimonio.Cadpat(pp)
		// patrimonio.Movbem(pp)
	}()

	wg.Wait()
}
