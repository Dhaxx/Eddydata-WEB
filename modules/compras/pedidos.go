package compras

import (
	"Eddydata-WEB/connection"
	"Eddydata-WEB/modules"
	"fmt"

	"github.com/vbauerster/mpb"
)

func Cadped(p *mpb.Progress) {
	modules.LimpaTabela([]string{"ICADPED","CADPED"})
	modules.NewCol("CADPED", "ID_CONTRATO_ANT")

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

	insertCadped, err := tx.Prepare(`INSERT
		INTO
		cadped(numped,
		num,
		ano,
		codif,
		datped,
		ficha,
		codccusto,
		entrou,
		id_cadorc,
		obs,
		id_cadped,
		empresa)
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar statement: %v", err.Error()))
	}
	defer insertCadped.Close()

	insertIcadped, err := tx.Prepare(`INSERT
		INTO
		icadped(numped,
		id_cadped,
		item,
		cadpro,
		codccusto,
		qtd,
		prcunt,
		prctot,
		ficha)
	VALUES(?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar statement: %v", err.Error()))
	}
	defer insertIcadped.Close()

	query := `select
		concat(to_char(a.numero, 'fm00000/'), exercicio_id) numped,
		to_char(a.numero, 'fm00000') num,
		to_char(exercicio_id, 'fm2000') ano,
		favorecido_id codif,
		data_compra datped,
		ficha_id,
		b.codigo codccusto,
		rcms_id id_cadorc,
		concat('COMPRA ', concat(to_char(a.numero, 'fm00000/'), exercicio_id)) obs,
		a.id id_cadped,
		c.ordem item,
		c.produto_unidade_id codreduz,
		quantidade, 
		c.valor_unitario prcunt,
		quantidade * c.valor_unitario prctot
	from
		"Y29" a
	join "Y153" b on
		a.setor_id = b.id
	join "Y30" c on
		c.compra_id = a.id
	order by
		a.ID`

	rows, err := cnxOrig.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query: %v", err.Error()))
	}
	defer rows.Close()

	var inseridos = make(map[int64]bool)
	var item int64

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "PEDIDOS")

	for rows.Next() {
		var registro ModelPedido

		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro: %v", err.Error()))
		}

		if !inseridos[registro.IdCadped] {
			var numlic string
			inseridos[registro.IdCadped] = true
			item = 0

			if err := cnxDest.QueryRow(fmt.Sprintf(`select numlic from cadlic where processo = '%s'`, registro.IdProcesso)).Scan(&numlic); err != nil {
				if err.Error() != "sql: no rows in result set" {
					panic(fmt.Sprintf("erro ao buscar numlic: %v", err.Error()))
				}
			}

			if _, err := insertCadped.Exec(
				registro.Numped,
				registro.Num,
				registro.Ano,
				registro.Codif,
				registro.Datped,
				registro.Ficha,
				registro.Codccusto,
				"N",
				registro.IdCadorc,
				registro.Obs,
				registro.IdCadped,
				modules.Cache.Empresa); err != nil {
				panic(fmt.Sprintf("erro ao inserir cadped: %v", err.Error()))
			}
		}

		item++
		if _, err := insertIcadped.Exec(
			registro.Numped,
			registro.IdCadped,
			registro.Item,
			modules.Cache.Cadpros[registro.Cadpro],
			registro.Codccusto,
			registro.Qtd,
			registro.Prcunt,
			registro.Pctot,
			registro.Ficha); err != nil {
			panic(fmt.Sprintf("erro ao inserir icadped: %v", err.Error()))
		}
		bar.Increment()
	}
	tx.Commit()
}
