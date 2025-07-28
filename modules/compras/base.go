package compras

import (
	"Eddydata-WEB/connection"
	"Eddydata-WEB/modules"
	"fmt"
	"github.com/vbauerster/mpb"
)

func Cadunimedida(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadunimedida"})

	cnxOrig, cnxDest, err := connection.GetConexoes()
    if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
    }
    defer cnxOrig.Close()
    defer cnxDest.Close()

	tx, err := cnxDest.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err.Error()))
	}
	defer tx.Commit()

	insert, err := tx.Prepare("insert into cadunimedida(sigla,descricao) values(?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert: %v", err.Error()))
	}
	defer insert.Close()

	query := `select substring(nome,1,4) sigla from "Y174" order by 1`
	
	rows, err := cnxOrig.Query(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta: %v", err.Error()))
	}
	defer rows.Close()

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar registros: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "CADUNIMEDIDA")

	for rows.Next() {
		var sigla string
		if err := rows.Scan(&sigla); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta: %v", err.Error()))
		}

		_, err = insert.Exec(sigla, sigla)
		if err != nil {
			panic(fmt.Sprintf("erro ao executar insert: %v", err.Error()))
		}
		bar.Increment()
	}
	bar.Completed()
}

func GrupoSubgrupo(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadsubgr","cadgrupo"})

	cnxOrig, cnxDest, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxOrig.Close()
	defer cnxDest.Close()

	tx, err := cnxDest.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err.Error()))
	}
	defer tx.Commit()

	insert, err := tx.Prepare("insert into cadgrupo(grupo,nome,balco_tce,balco_tce_saida,ocultar) values(?,?,?,?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadgrupo: %v", err.Error()))
	}

	query := `select TO_CHAR(ID, 'fm000') grupo, nome, null id_audesp, 'N' ocultar from "Y95"`

	rows, err := cnxOrig.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta cadgrupo: %v", err.Error()))
	}
	defer rows.Close()

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar registros cadgrupo: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "CADGRUPO")
	for rows.Next() {
		var registro ModelCadgrupo

		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta cadgrupo: %v", err.Error()))
		}

		_, err = insert.Exec(registro.Grupo, registro.Nome, registro.BalcoTce, registro.BalcoTce, registro.Ocultar)
		if err != nil {
			panic(fmt.Sprintf("erro ao executar insert cadgrupo: %v", err.Error()))
		}

		bar.Increment()
	}
	tx.Commit()

	query = `select to_char(grupo_id, 'fm000') grupo, to_char(id, 'fm000') subgrupo, nome, 'N' ocultar from "Y155"`

	tx, err = cnxDest.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err.Error()))
	}
	defer tx.Commit()

	insert, err = tx.Prepare("insert into cadsubgr(grupo,subgrupo,nome,ocultar) values(?,?,?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadsubgr: %v", err.Error()))
	}
	defer insert.Close()

	rows, err = cnxOrig.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta cadsubgr: %v", err.Error()))
	}
	defer rows.Close()

	totalRows, err = modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar registros cadsubgr: %v", err.Error()))
	}

	bar = modules.NewProgressBar(p, totalRows, "CADSUBGR")
	for rows.Next() {
		var registro ModelCadsubgr

		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta cadsubgr: %v", err.Error()))
		}

		_, err = insert.Exec(registro.Grupo, registro.Subgrupo, registro.Nome, registro.Ocultar)
		if err != nil {
			panic(fmt.Sprintf("erro ao executar insert cadsubgr: %v", err.Error()))
		}

		bar.Increment()
	}
}

func Cadest(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadest"})

	cnxOrig, cnxDest, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxOrig.Close()
	defer cnxDest.Close()

	tx, err := cnxDest.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err.Error()))
	}
	defer tx.Commit()

	insert, err := tx.Prepare("insert into cadest(grupo,subgrupo,codigo,cadpro,codreduz,disc1,ocultar,unid1,tipopro,usopro) values(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadest: %v", err.Error()))
	}
	defer insert.Close()

	query := `select
		to_char(b.id, 'fm000') grupo,
		b.sub_grupo_id,
		to_char(a.id, 'fm000') codigo,
		concat(to_char(b.id, 'fm000'),'.',b.sub_grupo_id,'.',to_char(a.id, 'fm000'))  cadpro,
		a.id codreduz,
		a.nome disc1,
		'N' inativo,
		substring(d.nome,1,4) unidade,
		'P' tipopro,
		'C' usopro
	from
		"Y126" a
	join "Y95" b on a.material_id = b.id 
	join "Y127" c on c.produto_id = a.id
	join "Y174" d on c.unidade_id = d.id
	order by
		1,2,3`

	rows, err := cnxOrig.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta cadest: %v", err.Error()))
	}
	defer rows.Close()

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar registros cadest: %v", err.Error()))
	}

	bar := modules.NewProgressBar(p, totalRows, "CADEST")
	for rows.Next() {
		var registro ModelCadest

		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta cadest: %v", err.Error()))
		}

		_, err = insert.Exec(registro.Grupo, registro.Subgrupo, registro.Codigo, registro.Cadpro, registro.Codreduz, registro.Disc1, registro.Ocultar, registro.Unid1, registro.Tipopro, registro.Usopro)
		if err != nil {
			panic(fmt.Sprintf("erro ao executar insert cadest: %v", err.Error()))
		}

		bar.Increment()
	}
	tx.Commit()

	if _, err := cnxDest.Exec("insert into cadsubgr(grupo, subgrupo, nome, ocultar) SELECT DISTINCT a.grupo, a.subgrupo, b.nome, b.ocultar FROM cadest a JOIN cadsubgr b USING (grupo, subgrupo) WHERE NOT EXISTS (SELECT 1 FROM cadsubgr x WHERE x.grupo = a.GRUPO AND x.SUBGRUPO = a.SUBGRUPO)"); err != nil {
		panic(fmt.Sprintf("erro ao executar insert cadsubgr: %v", err.Error()))
	}

	cadpros, err := cnxDest.Query(`select cadpro, 
			codreduz material
			From cadest t join cadgrupo g on g.GRUPO = t.GRUPO`)
		if err != nil {
			panic("Falha ao executar consulta: " + err.Error())
		}
		defer cadpros.Close()

	modules.Cache.Cadpros = make(map[string]string)
	for cadpros.Next() {
		var cadpro string
		var material string
		if err := cadpros.Scan(&cadpro, &material); err != nil {
			panic("Falha ao ler resultados da consulta: " + err.Error())
		}
		modules.Cache.Cadpros[material] = cadpro
	}
}

func Destino(p *mpb.Progress) {
	modules.LimpaTabela([]string{"destino"})

	cnxOrig, cnxDest, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxOrig.Close()
	defer cnxDest.Close()

	tx, err := cnxDest.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err.Error()))
	}
	defer tx.Commit()

	insert, err := tx.Prepare("insert into destino(cod,desti,empresa) values(?,?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert destino: %v", err.Error()))
	}

	query := fmt.Sprintf(`select to_char(id,'fm000000000') cod, nome from "Y58" where orgao_id = %v`, modules.Cache.Empresa)

	rows, err := cnxOrig.Query(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta destino: %v", err.Error()))
	}
	defer rows.Close()

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar registros destino: %v", err.Error()))
	}

	bar := modules.NewProgressBar(p, totalRows, "DESTINO")
	for rows.Next() {
		var cod, desti string
		if err := rows.Scan(&cod, &desti); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta destino: %v", err.Error()))
		}

		_, err := insert.Exec(cod, desti, modules.Cache.Empresa)
		if err != nil {
			panic(fmt.Sprintf("erro ao executar insert destino: %v", err.Error()))
		}

		bar.Increment()
	}
}

func CentroCusto(p *mpb.Progress) {
	modules.LimpaTabela([]string{"centrocusto"})

	cnxOrig, cnxDest, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxDest.Close()
	defer cnxOrig.Close()

	tx, err := cnxDest.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err.Error()))
	}
	defer tx.Commit()

	if _, err := cnxDest.Exec("INSERT INTO CENTROCUSTO(PODER, ORGAO, DESTINO, DESCR, CODCCUSTO, CCUSTO, EMPRESA) SELECT '02','01',(SELECT FIRST 1 COD FROM DESTINO), 'CONVERSAO', 0, '1', empresa from destino"); err != nil {
		panic(fmt.Sprintf("erro ao inserir codccusto: %v", err.Error()))
	}

	insert, err := cnxDest.Prepare("INSERT INTO CENTROCUSTO(PODER, ORGAO, DESTINO, DESCR, CODCCUSTO, CCUSTO, EMPRESA) values (?,?,?,?,?,?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert codccusto: %v", err.Error()))
	}

	query := `select
		'01' poder,
		'01' orgao,
		upper(nome) nome,
		codigo codccusto,
		1 ccusto
	from
		"Y153"`

	rows, err := cnxOrig.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta codccusto: %v", err.Error()))
	}
	defer rows.Close()

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar registros codccusto: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "CODCCUSTO")
	
	for rows.Next() {
		var registro modelCentroCusto

		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta codccusto: %v", err.Error()))
		}

		_, err := insert.Exec(registro.Poder, registro.Orgao, registro.Destino, registro.Descr, registro.Codccusto, registro.Ccusto, modules.Cache.Empresa)
		if err != nil {
			panic(fmt.Sprintf("erro ao executar insert codccusto: %v", err.Error()))
		}
		bar.Increment()
	}	
}