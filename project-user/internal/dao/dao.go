package dao

import (
	"github.com/a754962942/project-user/internal/database"
	"github.com/a754962942/project-user/internal/database/gorms"
)

type TransactionImpl struct {
	conn database.DbConn
}

func (t *TransactionImpl) Action(f func(conn database.DbConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	if err != nil {
		t.conn.Rollback()
		return err
	}
	t.conn.Commit()
	return nil
}

func NerTransactionImpl() *TransactionImpl {
	return &TransactionImpl{conn: gorms.NewTran()}
}
