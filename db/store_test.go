package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/yosa/ocr-golang-back/util"
)

func TestCreateDocumentWithTexts(t *testing.T) {
	user := createRandomUser(t)

	docArg := CreateDocumentParams{
		UserID:   user.Username,
		Filename: pgtype.Text{String: util.RandomFilename(), Valid: true},
		FileType: pgtype.Text{String: "pdf", Valid: true},
	}

	var texts []CreateExtractedTextParams
	for i := 0; i < 3; i++ {
		texts = append(texts, CreateExtractedTextParams{
			Content: pgtype.Text{String: util.RandomContent(), Valid: true},
		})
	}

	err := testStore.CreateDocumentWithTexts(context.Background(), docArg, texts)
	require.NoError(t, err)

	documents, err := testQueries.ListDocumentsByUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, documents)

	textsInDb, err := testQueries.ListExtractedTextsByDocument(context.Background(), documents[0].ID)
	require.NoError(t, err)
	require.Len(t, textsInDb, 3)
}

