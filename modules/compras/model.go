package compras

import (
	// "time"

	"github.com/gobuffalo/nulls"
)

type ModelCadgrupo struct {
	Grupo          string       `db:"grupo"`
	Nome           string       `db:"nome"`
	BalcoTce       nulls.String `db:"id_audesp"`
	Ocultar        string       `db:"ocultar"`
}

type ModelCadsubgr struct {
	Grupo          string       `db:"grupo"`
	Subgrupo       string       `db:"subgrupo"`
	Nome           string       `db:"nome"`
	Ocultar        string       `db:"ocultar"`
}

type ModelCadest struct {
	Grupo 		string       `db:"grupo"`
	Subgrupo 	string       `db:"subgrupo"`
	Codigo 		string       `db:"codigo"`
	Cadpro 		string       `db:"cadpro"`
	Codreduz 	string       `db:"codreduz"`
	Disc1 		string       `db:"disc1"`
	Ocultar 	string       `db:"inativo"`
	Tipopro 	string       `db:"tipopro"`
	Usopro 		string       `db:"usopro"`
	Unid1 		nulls.String  `db:"unidade"`
}

type modelCentroCusto struct {
	Poder 		string       `db:"poder"`
	Orgao 		string       `db:"orgao"`
	Destino 	string       `db:"destino"`
	Descr 		string       `db:"nome"`
	Codccusto 	string       `db:"codccusto"`
	Ccusto 		string       `db:"ccusto"`
}

type ModelCadorc struct {
	IdCadorc    int64        `db:"id_cadorc"`
	Num 	   	string       `db:"num"`
	Ano 	   	string       `db:"ano"`
	Numorc      string       `db:"numorc"`
	Dtorc    	nulls.Time   `db:"dtorc"`
	Descr 	 	string       `db:"descr"`
	Prioridade  string 	 	 `db:"prioridade"`
	Status      string     	 `db:"status"`
	Liberado    string       `db:"liberado"`
	Codccusto   string       `db:"codccusto"`
	LiberadoTela string       `db:"liberado_tela"`
}

type ModelIcadorc struct {
	Numorc   string
	Item   int64 	  `db:"item"`
	Cadpro string       
	Codreduz string 	`db:"codreduz"`
	Qtd  float64     `db:"qtd"`
	Valor float64    `db:"valor"`
	IdCadorc int64        `db:"id_cadorc"`
	Codccusto string 	
	Itemorc int64 
}

type ModelVcadorc struct {
	Numorc string
	Item int64        `db:"item"`
	Codif int64		`db:"codif"`
	Vlruni float64    `db:"valor_unitario"`
	Vlrtot float64    `db:"total"`
	IdCadorc int64        `db:"id_cadorc"`
	Classe string
}

type ModelPedido struct {
	Numped  string `db:"numped"`
	Num    string `db:"num"`
	Ano    string `db:"ano"`
	Codif  string `db:"codif"`
	Datped nulls.Time `db:"datped"`
	Ficha  nulls.Int64 `db:"ficha_id"`
	Codccusto int64 `db:"codccusto"`
	Entrou string 
	Obs     nulls.String `db:"obs"`
	IdCadped int64 `db:"id_cadped"`
	IdContrato nulls.String `db:"id_contrato"`
	Codreduz string `db:"codreduz"`
	Cadpro string
	Qtd   float64 `db:"quantidade"`
	Prcunt float64 `db:"prcunt"`
	Pctot float64 `db:"prctot"`
	Item int64 `db:"item"`
	IdCadorc int64 `db:"id_cadorc"`
}