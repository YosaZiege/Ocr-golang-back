package db


import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/yosa/ocr-golang-back/util"
)

func TestCreateDocument(t *testing.T)         {

arg := CreateDocumentParams{
		ID:     util.RandomInit(),
	}
}
func TestDeleteDocument(t *testing.T)         {}
func TestGetDocumentByID(t *testing.T)        {}
func TestListDocumentsByUser(t *testing.T)    {}
func TestUpdateDocumentFilename(t *testing.T) {}
