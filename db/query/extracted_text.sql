-- name: CreateExtractedText :one
INSERT INTO extracted_texts (id, document_id, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetExtractedTextByID :one
SELECT * FROM extracted_texts
WHERE id = $1;

-- name: ListExtractedTextsByDocument :many
SELECT * FROM extracted_texts
WHERE document_id = $1
ORDER BY created_at DESC  LIMIT $1 OFFSET $2 ;

-- name: UpdateExtractedTextContent :exec
UPDATE extracted_texts
SET content = $2
WHERE id = $1;

-- name: DeleteExtractedText :exec
DELETE FROM extracted_texts
WHERE id = $1;

