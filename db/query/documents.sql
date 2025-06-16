-- name: CreateDocument :one
INSERT INTO documents (id, user_id, filename, file_type)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDocumentByID :one
SELECT * FROM documents
WHERE id = $1;

-- name: ListDocumentsByUser :many
SELECT * FROM documents
WHERE user_id = $1
ORDER BY uploaded_at DESC LIMIT $1 OFFSET $2;

-- name: UpdateDocumentFilename :exec
UPDATE documents
SET filename = $2
WHERE id = $1;

-- name: DeleteDocument :exec
DELETE FROM documents
WHERE id = $1;

