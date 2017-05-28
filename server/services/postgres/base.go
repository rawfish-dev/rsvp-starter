package postgres

import (
	"time"

	"gopkg.in/gorp.v1"
)

type baseModel struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (b *baseModel) PreInsert(s gorp.SqlExecutor) error {
	b.CreatedAt = time.Now()
	b.UpdatedAt = b.CreatedAt
	return nil
}

func (b *baseModel) PreUpdate(s gorp.SqlExecutor) error {
	b.UpdatedAt = time.Now()
	return nil
}
