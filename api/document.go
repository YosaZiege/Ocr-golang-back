package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/otiai10/gosseract"

	"github.com/yosa/ocr-golang-back/db"
	"github.com/yosa/ocr-golang-back/token"
)

// Cleanup only .png files matching the doc ID (safe and specific)
func cleanupDocumentPNGs(docID string) {
	uploadDir := "uploads"
	pattern := filepath.Join(uploadDir, fmt.Sprintf("%s_page-*.png", docID))

	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("Error finding PNGs for cleanup: %v", err)
		return
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			log.Printf("Failed to delete %s: %v", file, err)
		} else {
			log.Printf("Deleted: %s", file)
		}
	}
}

func extractTextFromPDFWithOCR(ctx context.Context, pdfPath, docID string) (string, error) {
	if _, err := os.Stat(pdfPath); err != nil {
		return "", fmt.Errorf("PDF file not found: %w", err)
	}

	outputPrefix := filepath.Join("uploads", fmt.Sprintf("%s_page", docID))

	// Add timeout for PDF conversion
	convertCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(convertCtx, "pdftoppm", "-png", pdfPath, outputPrefix)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(
			"failed to convert pdf to images: %w (stderr: %s)",
			err,
			stderr.String(),
		)
	}

	// Find all generated images
	images, err := filepath.Glob(fmt.Sprintf("%s-*.png", outputPrefix))
	if err != nil {
		return "", fmt.Errorf("failed to find images: %w", err)
	}

	if len(images) == 0 {
		return "", fmt.Errorf("no images generated from PDF")
	}

	client := gosseract.NewClient()
	defer client.Close()

	// Configure Tesseract for better results
	client.SetLanguage("eng")
	client.SetPageSegMode(gosseract.PSM_AUTO) // Add this

	var allText bytes.Buffer

	for _, imgPath := range images {
		// Clean up immediately after processing
		defer os.Remove(imgPath)

		if err := client.SetImage(imgPath); err != nil {
			return "", fmt.Errorf("failed to set image %s: %w", imgPath, err)
		}

		text, err := client.Text()
		if err != nil {
			// Consider: should one page failure fail the whole document?
			// Or log and continue?
			return "", fmt.Errorf("ocr error on %s: %w", imgPath, err)
		}

		allText.WriteString(text)
		allText.WriteString("\n\n")
	}

	return strings.TrimSpace(allText.String()), nil
}

func (s *Server) UploadDocument(ctx *gin.Context) {
	// 1. Authenticated user
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// 2. Upload file
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer file.Close()

	docID := uuid.New().String()
	uploadDir := "uploads"
	uploadPath := filepath.Join(uploadDir, fmt.Sprintf("%s.pdf", docID))

	// 3. Ensure upload folder exists
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(fmt.Errorf("failed to create upload dir: %w", err)),
		)
		return
	}

	// 4. Save PDF
	outFile, err := os.Create(uploadPath)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(fmt.Errorf("failed to save file: %w", err)),
		)
		return
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(fmt.Errorf("failed to write file: %w", err)),
		)
		return
	}

	// 5. Create document record
	_, err = s.queries.CreateDocument(ctx, db.CreateDocumentParams{
		ID:       docID,
		UserID:   authPayload.Username,
		Filename: pgtype.Text{String: header.Filename, Valid: true},
		FileType: pgtype.Text{String: header.Header.Get("Content-Type"), Valid: true},
	})
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(fmt.Errorf("failed to store document: %w", err)),
		)
		return
	}

	// 6. Extract OCR text
	content, err := extractTextFromPDFWithOCR(ctx.Request.Context(), uploadPath, docID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("OCR failed: %w", err)))
		return
	}
	if content == "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "No extractable text found in PDF"})
		return
	}

	// 7. Store OCR result
	_, err = s.queries.CreateExtractedText(ctx, db.CreateExtractedTextParams{
		ID:         uuid.New().String(),
		DocumentID: docID,
		Content:    pgtype.Text{String: content, Valid: true},
	})
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(fmt.Errorf("failed to save extracted text: %w", err)),
		)
		return
	}

	// 8. Optional cleanup
	go cleanupDocumentPNGs(docID) // clean .png
	go os.Remove(uploadPath)      // clean the .pdf

	// 9. Return success
	ctx.JSON(http.StatusOK, gin.H{
		"document_id": docID,
		"content":     content,
		"message":     "Document uploaded and text extracted successfully",
	})
}

type fetchDocumentsRequest struct {
	Username string `json:"username" binding:"required"`
	Limit    int32  `json:"limit"    binding:"required"`
	Offset   int32  `json:"offset"   binding:"required"`
}

func (s *Server) FetchDocuments(ctx *gin.Context) {
	var req fetchDocumentsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	documents, err := s.queries.ListDocumentsByUser(ctx, db.ListDocumentsByUserParams{
		UserID: req.Username,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Fetch Documents"})
		return
	}
	ctx.JSON(http.StatusOK, documents)
}
