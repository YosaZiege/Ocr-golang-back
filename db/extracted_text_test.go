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

func CreateExtractedText(t *testing.T)          {}
func DeleteExtractedText(t *testing.T)          {}
func GetExtractedTextByID(t *testing.T)         {}
func ListExtractedTextsByDocument(t *testing.T) {}
func UpdateExtractedTextContent(t *testing.T)   {}
