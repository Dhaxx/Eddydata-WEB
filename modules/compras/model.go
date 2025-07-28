package compras

import (
	// "time"

	"github.com/gobuffalo/nulls"
)

type ModelCadgrupo struct {
	Grupo          string       `db:"GRUPO"`
	Nome           string       `db:"NOME"`
	BalcoTce       nulls.String       `db:"ID_AUDESP"`
	Ocultar        string       `db:"OCULTAR"`
}

type ModelCadsubgr struct {
	Grupo          string       `db:"GRUPO"`
	Subgrupo       string       `db:"SUBGRUPO"`
	Nome           string       `db:"NOME"`
	Ocultar        string       `db:"OCULTAR"`
}

type ModelCadest struct {
	Grupo 		string       `db:"GRUPO"`
	Subgrupo 	string       `db:"SUBGRUPO"`
	Codigo 		string       `db:"CODIGO"`
	Cadpro 		string       `db:"CADPRO"`
	Codreduz 	string       `db:"CODREDUZ"`
	Disc1 		string       `db:"DISC1"`
	Ocultar 	string       `db:"INATIVO"`
	Tipopro 	string       `db:"TIPOPRO"`
	Usopro 		string       `db:"USOPRO"`
	Unid1 		nulls.String  `db:"UNIDADE"`
}

type modelCentroCusto struct {
	Poder 		string       `db:"PODER"`
	Orgao 		string       `db:"ORGAO"`
	Destino 	string       `db:"DESTINO"`
	Descr 		string       `db:"DESCR"`
	Codccusto 	string       `db:"CODCCUSTO"`
	Ccusto 		string       `db:"CCUSTO"`
}

type ModelCadorc struct {
	IdCadorc    int64        `db:"ID_CADORC"`
	Num 	   	string       `db:"NUM"`
	Ano 	   	string       `db:"ANO"`
	Numorc      string       `db:"NUMORC"`
	Dtorc    	nulls.Time   `db:"DTORC"`
	Descr 	 	string `db:"DESCR"`
	Prioridade  string 	 	 `db:"PRIORIDADE"`
	Status      string     	 `db:"STATUS"`
	Liberado    string       `db:"LIBERADO"`
	Codccusto   string       `db:"CODCCUSTO"`
	LiberadoTela string       `db:"LIBERADO_TELA"`
}

type ModelIcadorc struct {
	Numorc   string
	Item   int64 	  `db:"ITEM"`
	Cadpro string       
	Codreduz string 	`db:"CODREDUZ"`
	Qtd  float64     `db:"QTD"`
	Valor float64    `db:"VALOR"`
	IdCadorc int64        `db:"ID_CADORC"`
	Codccusto string 	
	Itemorc int64 
}

type ModelVcadorc struct {
	Numorc string
	Item int64        `db:"ITEM"`
	Codif int64		`db:"CODIF"`
	Vlruni float64    `db:"VALOR_UNITARIO"`
	Vlrtot float64    `db:"VALOR_TOTAL"`
	IdCadorc int64        `db:"ID_CADORC"`
	Classe string
}

type ModelCadlic struct {
	Licit           string       
	Numpro          nulls.Int    `db:"NUMPRO"`
	Datae           nulls.Time   `db:"DATAE"`
	Dtpub           nulls.Time   `db:"DTPUB"`
	Dtenc           nulls.Time   `db:"DTENC"`
	Horenc          nulls.String `db:"HORENC"`
	Horabe          nulls.String `db:"HORABE"`
	Discr           nulls.String `db:"DISCR"`
	Discr7          nulls.String `db:"DISCR7"`
	Modlic          string       
	Dthom           nulls.Time   `db:"DTHOM"`
	Dtadj           nulls.Time   `db:"DTADJ"`
	Comp            nulls.String `db:"COMP"`
	Numero          nulls.String `db:"NUMERO"`
	Ano             nulls.String `db:"ANO"`
	Valor           nulls.Float64 `db:"VALOR"`
	Tipopubl        nulls.String `db:"TIPOPUBL"`
	Detalhe         nulls.String `db:"DETALHE"`
	Horreal         nulls.String `db:"HORREAL"`
	Local           nulls.String `db:"LOCAL"`
	Proclic         string       `db:"PROCLIC"`
	Numlic          int64        `db:"NUMLIC"`
	Liberacompra    nulls.String `db:"LIBERACOMPRA"`
	Microempresa    nulls.String `db:"MICROEMPRESA"`
	Licnova         nulls.String `db:"LICNOVA"`
	Codtce          nulls.String `db:"CODTCE"`
	ProcessoData    nulls.Time   `db:"PROCESSO_DATA"`
	Codmod          int64        
	Anomod          nulls.String `db:"ANOMOD"`
	Registropreco   nulls.String `db:"REGISTROPRECO"`
	Empresa         int64        `db:"EMPRESA"`
	Modalidade      int        `db:"PKANT_ID_MODALIDADE"`
	Processo 		nulls.String `db:"PROCESSO"`
	ProcessoAno    nulls.String `db:"PROCESSO_ANO"`
	Dtreal 	   nulls.Time   `db:"DTREAL"`
	ItensAgrup       nulls.String `db:"ITENS_AGRUP"`
	CodTce 		 nulls.String `db:"COD_TCE"`
}

type ModelCadprolic struct {
	Numlic 		int64        `db:"NUMLIC"`
	Item 		int64        `db:"ITEM"`
	ItemMask    int64		 `db:"ITEM_MASK"`
	Microempresa string       `db:"MICROEMPRESA"`
	Cadpro 		string       `db:"CADPRO"`
	Quan1 		float64      `db:"QUAN1"`
	Vamed1 		float64      `db:"VAMED1"`
	Valor 		float64      `db:"VALOR"`
	Codccusto   string       `db:"CODCCUSTO"`
	Reduz 		string       `db:"REDUZ"`
	Tlance      string 	 	 `db:"TLANCE"`
}

type ModelProlics struct {
	Sessao  int64        `db:"SESSAO"`
	Codif  string       `db:"CODIF"`
	Numlic int64        `db:"NUMLIC"`
	Habilitado nulls.String `db:"HABILITADO"`
	Status nulls.String `db:"STATUS"`
	Nome  nulls.String `db:"NOME"`
	Representante nulls.String `db:"NOME"`
}

type ModelProposta struct {
	Sessao   int64        `db:"SESSAO"`
	Codif  string       `db:"CODIF"`
	Item    int64        `db:"ITEM"`
	Itemp  int64        `db:"ITEMP"`
	Quan1  nulls.Float64 `db:"QUAN1"`
	Vaun1 nulls.Float64 `db:"VAUN1"`
	Vato1 nulls.Float64 `db:"VATO1"`
	Numlic int64        `db:"NUMLIC"`
	Status  nulls.String `db:"STATUS"`
	Subem nulls.String `db:"SUBEM"`
	Marca nulls.String `db:"MARCA"`
	ItemLance string        `db:"ITEMLANCE"`
}

type ModelRequi struct {
	IdRequi    int64         `db:"ID_REQUI"`
	Requi      nulls.String        `db:"REQUI"`
	Datae      nulls.Time    `db:"DATAE"`
	Dtlan 	   nulls.Time    `db:"DTLAN"`
	Dtpag      nulls.Time    `db:"DTPAG"`
	Ano        int64         `db:"ANO"`
	Destino    string        `db:"DESTINO"`
	Entr       string        `db:"ENTR"`
	Said       string        `db:"SAID"`
	Comp       string        `db:"COMP"`
	Codif      nulls.String  `db:"CODIF"`
	Docum 	   nulls.String  `db:"DOCUM"`
	Recebe     nulls.String  `db:"RECEBE"`
	Codccusto  nulls.Int64   `db:"CODCCUSTO"`
}

type ModelIcadreq struct {
	IdRequi   int64         `db:"ID_REQUI"`
	Requi     string        `db:"REQUI"`
	Codccusto int64         `db:"CODCCUSTO"`
	Item 	int64         `db:"ITEM"`
	Destino   string        `db:"DESTINO"`
	Cadpro  string        `db:"CADPRO"`
	Quan1  float64      `db:"QUAN1"`
	Vaun1 float64      `db:"VAUN1"`
	Quan2 float64      `db:"QUAN2"`
	Vaun2 float64      `db:"VAUN2"`
	Motorista nulls.String `db:"MOTORISTA"`
	Placa    nulls.String `db:"PLACA"`
	Km nulls.String `db:"KM"`
}

type ModelPedido struct {
	Numped  string `db:"NUMPED"`
	Num    string `db:"NUM"`
	Ano    string `db:"ANO"`
	Codif  string `db:"CODIF"`
	Datped nulls.Time `db:"DATPED"`
	Ficha  nulls.Int64 `db:"FICHA"`
	Codccusto int64 `db:"CODCCUSTO"`
	Entrou string `db:"ENTROU"`
	IdProcesso string `db:"PROCESSO"`
	Obs     nulls.String `db:"OBS"`
	IdCadped int64 `db:"ID_CADPED"`
	IdContrato nulls.String `db:"ID_CONTRATO"`
	Cadpro string `db:"ID_MATERIAL"`
	Qtd   float64 `db:"QTD"`
	Prcunt float64 `db:"PRCUNT"`
	Pctot float64 `db:"PRCTOT"`
	Item int64 `db:"ITEM"`
	IdCadorc int64 `db:"ID_CADORC"`
}