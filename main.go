package main

import (
	"Eddydata-WEB/modules"
	"Eddydata-WEB/modules/compras"

	"github.com/vbauerster/mpb"
)

func main() {
	pc := mpb.New()
	modules.LimpaCompras()
	compras.Cadunimedida(pc)
	compras.GrupoSubgrupo(pc)
	compras.Cadest(pc)
	compras.Destino(pc)
	compras.CentroCusto(pc)

	compras.Cadorc(pc)
	compras.Icadorc(pc)
	compras.Vcadorc(pc)

	compras.Cadped(pc)
}