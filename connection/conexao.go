package connection

import (
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/joho/godotenv"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    _ "github.com/nakagami/firebirdsql"
    "database/sql"
)

var dsnFdb string
var dsnPg string
var dsnFdbAnx string

func init() {
    envPath, err := os.Getwd()
    if err != nil {
        log.Fatalf("Erro ao obter diretório: %v", err)
    }

    if err = godotenv.Load(filepath.Join(envPath, ".env")); err != nil {
        log.Fatalf("Erro ao carregar .env: %v", err)
    }

    dsnFdb = fmt.Sprintf("%s:%s@%s:%s/%s?charset=win1252&auth_plugin_name=Legacy_Auth",
        os.Getenv("FDB_USER"),
        os.Getenv("FDB_PASS"),
        os.Getenv("FDB_HOST"),
        os.Getenv("FDB_PORT"),
        os.Getenv("FDB_PATH"))

    dsnFdbAnx = fmt.Sprintf("%s:%s@%s:%s/%s?charset=win1252&auth_plugin_name=Legacy_Auth",
        os.Getenv("FDB_USER_ANX"),
        os.Getenv("FDB_PASS_ANX"),
        os.Getenv("FDB_HOST_ANX"),
        os.Getenv("FDB_PORT_ANX"),
        os.Getenv("FDB_PATH_ANX"))

    dsnPg = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("PG_HOST"),
        os.Getenv("PG_PORT"),
        os.Getenv("PG_USER"),
        os.Getenv("PG_PASS"),
        os.Getenv("PG_DB"))
}

// GetConexoes retorna conexão Firebird via sql e PostgreSQL via sqlx
func GetConexoes() (*sqlx.DB, *sql.DB, error) {
    // Conexão com Firebird (segue usando sql.DB)
    ConexaoFdb, err := sql.Open("firebirdsql", dsnFdb)
    if err != nil {
        return nil, nil, fmt.Errorf("erro ao estabelecer conexão FDB: %v", err)
    }

    // Conexão com PostgreSQL usando sqlx
    ConexaoPg, err := sqlx.Connect("postgres", dsnPg)
    if err != nil {
        ConexaoFdb.Close()
        return nil, nil, fmt.Errorf("erro ao estabelecer conexão PostgreSQL: %v", err)
    }

    return ConexaoPg, ConexaoFdb, err
}

func GetConexoesAnexos() (*sql.DB, *sqlx.DB, error) {
    // Conexão com Firebird (segue usando sql.DB)
    ConexaoFdbAnx, err := sql.Open("firebirdsql", dsnFdbAnx)
    if err != nil {
        return nil, nil, fmt.Errorf("erro ao estabelecer conexão FDB: %v", err)
    }

    // Conexão com PostgreSQL usando sqlx
    ConexaoPg, err := sqlx.Connect("postgres", dsnPg)
    if err != nil {
        ConexaoFdbAnx.Close()
        return nil, nil, fmt.Errorf("erro ao estabelecer conexão PostgreSQL: %v", err)
    }

    return ConexaoFdbAnx, ConexaoPg, nil
}