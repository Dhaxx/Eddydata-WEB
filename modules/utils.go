package modules

import (
	"Eddydata-WEB/connection"
	"bytes"
	"database/sql"
	"fmt"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

var Cache struct {
	Empresa int
	Ano     int
	Cadpros map[string]string
}

type ProcessoLicitatorio struct {
	Nro    int
	Licit  string
	Modlic string
	Codmod int
}

var Modalidades = []ProcessoLicitatorio{
	{Nro: 1, Licit: "CONCURSO", Modlic: "CS01", Codmod: 7},
	{Nro: 2, Licit: "MAT / SERV - CONVITE", Modlic: "CC02", Codmod: 2},
	{Nro: 3, Licit: "MAT / SERV - TOMADA", Modlic: "TOM3", Codmod: 3},
	{Nro: 4, Licit: "MAT / SERV - CONCORRENCIA", Modlic: "CON4", Codmod: 4},
	{Nro: 5, Licit: "DISPENSA", Modlic: "DI01", Codmod: 1},
	{Nro: 6, Licit: "INEXIGIBILIDADE", Modlic: "IN01", Codmod: 5},
	{Nro: 7, Licit: "PREGÃO PRESENCIAL", Modlic: "PP01", Codmod: 8},
	{Nro: 8, Licit: "PREGÃO ELETRÔNICO", Modlic: "PE01", Codmod: 9},
	{Nro: 9, Licit: "DISPENSA", Modlic: "DI01", Codmod: 1},
	{Nro: 10, Licit: "PREGÃO ELETRÔNICO", Modlic: "PE01", Codmod: 9},
	{Nro: 11, Licit: "LEILÃO", Modlic: "LEIL", Codmod: 6},
	{Nro: 12, Licit: "DISPENSA", Modlic: "DI01", Codmod: 1},
	{Nro: 13, Licit: "DISPENSA", Modlic: "DI01", Codmod: 1},
	{Nro: 14, Licit: "CONCORRÊNCIA ELETRÔNICA", Modlic: "CE01", Codmod: 13},
	{Nro: 15, Licit: "DISPENSA ELETRÔNICA", Modlic: "DE01", Codmod: 11},
	{Nro: 16, Licit: "DISPENSA", Modlic: "DI01", Codmod: 1},
	{Nro: 17, Licit: "DISPENSA", Modlic: "DI01", Codmod: 1},
}

func init() {
	_, cnxFdb, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()

	cnxFdb.QueryRow("Select empresa from cadcli").Scan(&Cache.Empresa)
	cnxFdb.QueryRow("Select mexer from cadcli").Scan(&Cache.Ano)
	var cadestOk int
	_ = cnxFdb.QueryRow("Select count(*) from cadest").Scan(&cadestOk)
	if cadestOk == 0 {
		fmt.Print("Cadest vazia")
	} else {
		cadpros, err := cnxFdb.Query(`select cadpro, 
			codreduz material
			From cadest t join cadgrupo g on g.GRUPO = t.GRUPO`)
		if err != nil {
			panic("Falha ao executar consulta: " + err.Error())
		}
		defer cadpros.Close()

		Cache.Cadpros = make(map[string]string)
		for cadpros.Next() {
			var cadpro string
			var material string
			if err := cadpros.Scan(&cadpro, &material); err != nil {
				panic("Falha ao ler resultados da consulta: " + err.Error())
			}
			Cache.Cadpros[material] = cadpro
		}
	}
}

func LimpaTabela(tabelas []string) {
	_, cnxFdb, err := connection.GetConexoes()
	if err != nil {
		fmt.Printf("Falha ao conectar com o banco de destino: %v", err)
	}
	defer cnxFdb.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		fmt.Printf("erro ao iniciar transação: %v", err)
	}

	for _, tabela := range tabelas {
		if _, err = tx.Exec(fmt.Sprintf("DELETE FROM %v", tabela)); err != nil {
			fmt.Printf("erro ao limpar tabela: %v", err)
			tx.Rollback()
		}
	}
	tx.Commit()
}

func CountRows(q string, args ...any) (int64, error) {
	cnxFdb, _, err := connection.GetConexoes()
	if err != nil {
		fmt.Printf("Falha ao conectar com o banco de destino: %v", err)
	}
	defer cnxFdb.Close()

	var count int64
	query := fmt.Sprintf("SELECT count(*) FROM (%v) AS qr", q)

	if err := cnxFdb.QueryRow(query).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("nenhuma linha recuperada: %v", sql.ErrNoRows.Error())
		}
		return 0, fmt.Errorf("erro ao contar registros: %v", err)
	}
	return count, nil
}

func NewProgressBar(p *mpb.Progress, total int64, label string) *mpb.Bar {
	return p.AddBar(total,
		mpb.BarWidth(60),
		mpb.BarStyle("[██████░░░░░░]"),
		mpb.PrependDecorators(
			decor.Name(label+": "),
			decor.CountersNoUnit("%d / %d"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
			decor.EwmaETA(decor.ET_STYLE_GO, 60),
		),
	)
}

func NewCol(table string, colName string) {
	_, cnxFdb, err := connection.GetConexoes()
	if err != nil {
		fmt.Printf("Falha ao conectar com o banco de destino: %v", err)
	}
	defer cnxFdb.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		fmt.Printf("erro ao iniciar transação: %v", err)
	}

	_, err = tx.Exec(fmt.Sprintf("ALTER TABLE %v ADD %v varchar(50)", table, colName))
	if err != nil {
		tx.Rollback()
		fmt.Printf("erro ao criar coluna %v: %v", colName, err)
	}

	tx.Commit()
}

func DecodeWin1252FromBytes(b []byte) (string, error) {
	decoder := charmap.Windows1252.NewDecoder()
	return decoder.String(string(b))
}

func EncodeToWin1252(input string) (string, error) {
	// Define uma tabela de caracteres válidos no Windows-1252
	validChars := charmap.Windows1252

	// Remove ou substitui caracteres inválidos
	t := transform.Chain(
		runes.Remove(runes.Predicate(func(r rune) bool {
			// Remove caracteres que não são válidos no Windows-1252
			_, ok := validChars.EncodeRune(r)
			return !ok
		})),
		validChars.NewEncoder(),
	)

	// Transforma a string
	var buf bytes.Buffer
	writer := transform.NewWriter(&buf, t)

	_, err := writer.Write([]byte(input))
	if err != nil {
		return "", fmt.Errorf("erro ao codificar para Windows-1252: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("erro ao finalizar o writer: %w", err)
	}

	return buf.String(), nil
}

func LimpaLicitacoes() {
	_, cnxAux, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxAux.Close()

	_, err = cnxAux.Exec(`execute block as
		begin
		DELETE FROM regpreco;
		DELETE FROM regprecohis;
		DELETE FROM regprecodoc;
		DELETE FROM CADPROLIC_DETALHE_FIC;
		DELETE FROM CADPRO;
		DELETE FROM CADPRO_FINAL;
		DELETE FROM CADPRO_LANCE;
		DELETE FROM CADPRO_PROPOSTA;
		DELETE FROM PROLICS;
		DELETE FROM PROLIC;
		DELETE FROM CADPRO_STATUS;
		DELETE FROM CADLIC_SESSAO;
		DELETE FROM CADPROLIC_DETALHE;
		DELETE FROM CADPROLIC;
		DELETE FROM CADLIC;
		end;`)
	if err != nil {
		panic("Falha ao executar delete: " + err.Error())
	}
}

func LimpaCompras() {
	_, cnxAux, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxAux.Close()

	Trigger("TD_ICADREQ", false)

	_, err = cnxAux.Exec(`execute block as
		begin
		DELETE FROM ICADREQ;
		DELETE FROM REQUI;
		DELETE FROM ICADPED;
		DELETE FROM CADPED;
		DELETE FROM regpreco;
		DELETE FROM regprecohis;
		DELETE FROM regprecodoc;
		DELETE FROM CADPRO_SALDO_ANT;
		DELETE FROM CADPROLIC_DETALHE_FIC;
		DELETE FROM CADPRO;
		DELETE FROM CADPRO_FINAL;
		DELETE FROM CADPRO_LANCE;
		DELETE FROM CADPRO_PROPOSTA;
		DELETE FROM PROLICS;
		DELETE FROM PROLIC;
		DELETE FROM CADPRO_STATUS;
		DELETE FROM CADLIC_SESSAO;
		DELETE FROM CADPROLIC_DETALHE;
		DELETE FROM CADPROLIC;
		DELETE FROM CADLIC;
		DELETE FROM VCADORC;
		DELETE FROM FCADORC;
		DELETE FROM ICADORC;
		DELETE FROM CADORC;
		DELETE FROM CADEST;
		DELETE FROM CENTROCUSTO;
		DELETE FROM DESTINO;
		DELETE FROM DESFORCRC_PADRAO;
		end;`)
	if err != nil {
		panic("Falha ao executar delete: " + err.Error())
	}
}

func LimpaPatrimonio() {
	_, cnxAux, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxAux.Close()

	_, err = cnxAux.Exec(`execute block as
		begin
		DELETE FROM PT_CADPAT_EMPEN;
		DELETE FROM PT_MOVBEM;
		DELETE FROM PT_CADPAT;
		DELETE FROM PT_CADPATS;
		DELETE FROM PT_CADPATD;
		DELETE FROM PT_CADPATG;
		DELETE FROM PT_CADTIP;
		DELETE FROM PT_CADSIT;
		DELETE FROM PT_CADBAI;
		DELETE FROM PT_CADAJUSTE;
		DELETE FROM PT_TIPOMOV;
		DELETE FROM PT_CADRESPONSAVEL;
		end;`)
	if err != nil {
		panic("Falha ao executar delete: " + err.Error())
	}
}

func Trigger(trigger string, status bool) {
	_, cnxFdb, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic("erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	var statusStr string
	if status {
		statusStr = "ACTIVE"
	} else {
		statusStr = "INACTIVE"
	}

	tx.Exec(fmt.Sprintf("ALTER TRIGGER %s %s", trigger, statusStr))
}

func ExtourouSubgrupo(codant string) string {
	_, cnxFdb, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic("erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	if _, err = tx.Exec("INSERT INTO CADSUBGR(grupo, subgrupo, nome, ocultar, key_subgrupo, base)  select grupo, lpad(max(cast((SELECT max(subgrupo) FROM cadsubgr) as integer) + 1), 3, '0'), nome, 'N', key_subgrupo, 'N' from cadsubgr where key_subgrupo = ? GROUP BY 1, 3, 5", codant); err != nil {
		panic("Falha ao inserir novo subgrupo: " + err.Error())
	}
	tx.Commit()

	var novoSubgrupo string
	err = cnxFdb.QueryRow("SELECT max(subgrupo) FROM cadsubgr where key_subgrupo = ?", codant).Scan(&novoSubgrupo)
	if err != nil {
		panic("Falha ao recuperar novo subgrupo: " + err.Error())
	}

	return novoSubgrupo
}

func DesforAnt(codant string) string {
	_, cnxFdb, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()

	var codif string

	err = cnxFdb.QueryRow("SELECT codif FROM desfor WHERE codant = ?", codant).Scan(&codif)
	if err != nil {
		codif = "0"
	}

	return codif
}

func ContratosAdit(p *mpb.Progress) {
	_, cnxFdb, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()

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

	query := `SELECT
		id_parente contrato,
		DT_INICIO,
		dt_termino, 
		valor,
		objeto descricao,
		extract(year from dt_inicio) as ano
	FROM
		contabil_contrato
	WHERE
		id_parente IS NOT NULL
	ORDER BY 1,2`

	rows, err := cnxOrig.Query(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta: %v", err.Error()))
	}
	defer rows.Close()

	insert, err := tx.Prepare(`INSERT INTO contratosaditamento (CONTRATO,DTLAN,DATAENCERRAMENTO,VALOR,DESCRICAO,DTPUBLICACAO,DATAINSC,TIPOHIST,TIPOALT,TIPO_TCE,DATAABERTURA,TERMO,CODIGO) values (?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert: %v", err.Error()))
	}
	defer insert.Close()

	totalRows, err := CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar registros: %v", err.Error()))
	}

	bar := NewProgressBar(p, totalRows, "ADITAMENTOS")
	var (
		codigo      int64
		contratoAnt string
		termo       int64
		ano         int
	)
	tipoHist := "Prorrogação"
	tipoTce := 12
	tipoAlt := "Bilateral"

	for rows.Next() {
		var (
			contratoIdAnt, descr, contrato string
			valor                          float64
			dtIni, dtEnc                   sql.NullTime
		)
		codigo++

		if err := rows.Scan(&contratoIdAnt, &dtIni, &dtEnc, &valor, &descr, &ano); err != nil {
			panic(fmt.Sprintf("erro ao ler linha: %v", err.Error()))
		}

		if contratoIdAnt == contratoAnt {
			termo++
		} else {
			termo = 1
		}
		contratoAnt = contratoIdAnt

		termoFormatado := fmt.Sprintf("%05d/%v", termo, ano%2000)

		if err = cnxDest.QueryRow("SELECT codigo FROM contratos WHERE id_id_contrato = ?", contratoIdAnt).Scan(&contrato); err != nil {
			panic(fmt.Sprintf("erro ao consultar codigo: %v", err.Error()))
		}

		if _, err := insert.Exec(contrato, dtIni, dtEnc, valor, descr, dtIni, dtIni, tipoHist, tipoAlt, tipoTce, dtIni, termoFormatado, codigo); err != nil {
			panic(fmt.Sprintf("erro ao inserir linha: %v", err.Error()))
		}

		bar.Increment()
	}
}

func DecodeWin1252(b []byte) string {
    decoder := charmap.Windows1252.NewDecoder()
    result, err := decoder.Bytes(b)
    if err != nil {
        return fmt.Sprintf("[erro]: %v", err)
    }
    return string(result)
}