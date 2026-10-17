package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Record struct {
	ID             string
	TenantID       string
	VendaID        string
	ChaveAcesso    string
	NumeroNota     int
	Serie          int
	Status         string
	TpEmis         int
	XMLEnvio       []byte
	XMLRetorno     []byte
	MotivoRejeicao string
	Tentativas     int
}

type Repository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir sqlite: %w", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS fiscal_records (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		venda_id TEXT NOT NULL,
		chave_acesso TEXT UNIQUE,
		numero_nota INTEGER NOT NULL,
		serie INTEGER NOT NULL,
		status TEXT NOT NULL,
		tp_emis INTEGER NOT NULL,
		xml_envio BLOB NOT NULL,
		xml_retorno BLOB,
		motivo_rejeicao TEXT,
		tentativas INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_status_tp_emis ON fiscal_records(status, tp_emis);
	`
	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar tabelas: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Save(rec *Record) error {
	query := `INSERT INTO fiscal_records (id, tenant_id, venda_id, chave_acesso, numero_nota, serie, status, tp_emis, xml_envio, xml_retorno, motivo_rejeicao)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, rec.ID, rec.TenantID, rec.VendaID, rec.ChaveAcesso, rec.NumeroNota, rec.Serie, rec.Status, rec.TpEmis, rec.XMLEnvio, rec.XMLRetorno, rec.MotivoRejeicao)
	return err
}

func (r *Repository) GetPendingContingency() ([]*Record, error) {
	rows, err := r.db.Query("SELECT id, tenant_id, venda_id, chave_acesso, numero_nota, serie, status, tp_emis, xml_envio FROM fiscal_records WHERE status = 'CONTINGENCIA' AND tp_emis = 9 LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*Record
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.ID, &rec.TenantID, &rec.VendaID, &rec.ChaveAcesso, &rec.NumeroNota, &rec.Serie, &rec.Status, &rec.TpEmis, &rec.XMLEnvio); err != nil {
			return nil, err
		}
		records = append(records, &rec)
	}
	return records, nil
}

func (r *Repository) UpdateStatus(id string, status string, xmlRetorno []byte, motivo string) error {
	_, err := r.db.Exec("UPDATE fiscal_records SET status = ?, xml_retorno = ?, motivo_rejeicao = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", status, xmlRetorno, motivo, id)
	return err
}

func (r *Repository) Close() error {
	return r.db.Close()
}
