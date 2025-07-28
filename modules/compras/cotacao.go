package compras

import (
	"Eddydata-WEB/connection"
	"Eddydata-WEB/modules"
	"fmt"

	"github.com/vbauerster/mpb"
)

func Cadorc(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadorc"})
	modules.NewCol("cadorc", "anexo_ant")

	cnxOrig, cnxDest, err := connection.GetConexoes()
    if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
    }
    defer cnxOrig.Close()
    defer cnxDest.Close()

	tx, err := cnxDest.Begin()
	if err != nil {
		panic("Erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	insert, err := tx.Prepare(`insert
		into
		cadorc (id_cadorc,
		num,
		ano,
		numorc,    
		dtorc,
		descr,  
		prioridade,
		status,
		liberado,
		codccusto,
		liberado_tela,
		empresa) values (?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		fmt.Printf("erro ao preparar insert: %v", err)
	}
	defer insert.Close()

	query := `select
		a.id id_cadorc,
		to_char(a.numero, 'fm00000') num,
		extract(year from data_rcms) ano,
		concat(to_char(a.numero, 'fm00000'), '/', exercicio_id) numorc,
		data_rcms dtorc,
		concat('COTACAO ', concat(to_char(a.numero, 'fm00000'), '/', exercicio_id)) descr,
		'NORMAL' prioridade,
		'EC' status,
		'S' liberado,
		b.codigo codccusto,
		'P' liberado_tela
	from
		"Y132" a
	join "Y153" b on a.setor_id = b.id
	order by 1`

	totalLinhas, _ := modules.CountRows(query)
	bar := modules.NewProgressBar(p, totalLinhas, "Cadorc")

	rows, err := cnxOrig.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelCadorc
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler registro: %v", err))
		}
		registro.Descr, err = modules.DecodeToWin1252(registro.Descr); if err != nil {
			panic(fmt.Sprintf("erro ao decodificar descricao: %v", err))
		}
		if _, err = insert.Exec(
			registro.IdCadorc,
			registro.Num,
			registro.Ano,
			registro.Numorc,
			registro.Dtorc,
			registro.Descr,
			registro.Prioridade,
			registro.Status,
			registro.Liberado,
			registro.Codccusto,
			registro.LiberadoTela,
			modules.Cache.Empresa,
		); err != nil {
			fmt.Sprintf("erro ao inserir cadorc: %v", err)
		}
		bar.Increment()
	}
}

func Icadorc(p *mpb.Progress) {
	modules.LimpaTabela([]string{"icadorc"})
	modules.NewCol("icadorc", "lote_ant")
	modules.NewCol("icadorc", "seq_ant")
	cnxOrig, cnxDest, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxOrig.Close()
	defer cnxDest.Close()

	tx, err := cnxDest.Begin()
	if err != nil {
		panic("Erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	insertIcadorc, err := tx.Prepare(`insert into icadorc (numorc, item, cadpro, qtd, valor, itemorc, codccusto, itemorc_ag, id_cadorc) values (?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		fmt.Printf("erro ao preparar insert: %v", err)
	}
	defer insertIcadorc.Close()

	query := `select 
		rcms_id id_cadorc,
		row_number() over (partition by rcms_id order by ordem, produto_unidade_id) item,
		--ordem item,
		produto_unidade_id codreduz,
		quantidade qtd,
		valor_unitario valor
	from "Y135"`

	totalLinhasItens, _ := modules.CountRows(query)
	bar := modules.NewProgressBar(p, totalLinhasItens, "Icadorc")

	rows, err := cnxOrig.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelIcadorc
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler registro: %v", err))
		}

		registro.Cadpro = modules.Cache.Cadpros[registro.Codreduz]
		if registro.Cadpro == "" {
			panic(fmt.Sprintf("erro ao encontrar produto: %v", registro.Codreduz))
		}

		if err := cnxDest.QueryRow("select numorc, codccusto from cadorc where id_cadorc = ?", registro.IdCadorc).Scan(&registro.Numorc, &registro.Codccusto); err != nil {
			panic(fmt.Sprintf("erro ao buscar numorc e codccusto: %v", err))
		}

		registro.Itemorc = registro.Item

		if _, err = insertIcadorc.Exec(registro.Numorc, registro.Item, registro.Cadpro, registro.Qtd, registro.Valor, registro.Itemorc, registro.Codccusto, registro.Itemorc, registro.IdCadorc); err != nil {
			panic(fmt.Sprintf("erro ao inserir icadorc: %v", err))
		}

		bar.Increment()
	}
}

func Vcadorc(p *mpb.Progress) {
	modules.LimpaTabela([]string{"vcadorc"})

	cnxOrig, cnxDest, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxOrig.Close()
	defer cnxDest.Close()

	tx, err := cnxDest.Begin()
	if err != nil {
		panic("Erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	query := `select
		a.favorecido_id codif,
		--d.nome,
		row_number() over (partition by a.rcms_id order by ordem, produto_unidade_id) item,
		b.valor_unitario,
		c.quantidade * b.valor_unitario total,
		a.rcms_id id_cadorc
	from
		"Y134" a
	join "Y133" b on
		a.id = b.rcms_favorecido_id
	join "Y135" c on
		b.rcmsitem_id = c.id
	join "Y65" d on
		a.favorecido_id = d.id
	order by
		a.rcms_id,
		c.ordem`

	insert, err := tx.Prepare("insert into vcadorc(numorc, item, codif, vlruni, vlrtot, ganhou, vlrganhou, classe, id_cadorc) values (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Printf("erro ao preparar insert: %v", err)
	}
	defer insert.Close()

	totalLinhas, _ := modules.CountRows(query)
	bar := modules.NewProgressBar(p, totalLinhas, "Vcadorc")

	rows, err := cnxOrig.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelVcadorc
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler registro: %v", err))
		}

		if err := cnxDest.QueryRow("select numorc from cadorc where id_cadorc = ?", registro.IdCadorc).Scan(&registro.Numorc); err != nil {
			panic(fmt.Sprintf("erro ao buscar numorc: %v", err))
		}

		registro.Classe = "UN"

		if _, err = insert.Exec(registro.Numorc, registro.Item, registro.Codif, registro.Vlruni, registro.Vlrtot, registro.Codif, registro.Vlruni, registro.Classe, registro.IdCadorc); err != nil {
			panic(fmt.Sprintf("erro ao inserir vcadorc: %v", err))
		}
		bar.Increment()
	}
	tx.Commit()

	cnxDest.Exec("insert into fcadorc (numorc, codif, nome, valorc, id_cadorc) select numorc, codif, (select nome from desfor x where x.codif = a.codif), sum(vlrtot), id_cadorc from vcadorc a group by numorc, codif, id_cadorc")
}