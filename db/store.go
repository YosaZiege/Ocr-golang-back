package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries // this is called Composition Prefered over heritage
	db       *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error %v, rb err %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *Store) CreateDocumentWithTexts(
	ctx context.Context,
	docArg CreateDocumentParams,
	texts []CreateExtractedTextParams,
) error {
	return store.execTx(ctx, func(q *Queries) error {
		doc, err := q.CreateDocument(ctx, docArg)
		if err != nil {
			return err
		}

		for _, text := range texts {
			text.DocumentID = doc.ID
			if _, err := q.CreateExtractedText(ctx, text); err != nil {
				return err
			}
		}

		return nil
	})
}

