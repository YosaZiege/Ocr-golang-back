package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yosa/ocr-golang-back/db"
	"github.com/yosa/ocr-golang-back/util"
)

type createUserParams struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Provider string `json:"provider"`
}
type userResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email" `
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {

	return userResponse{
		Username:  user.Username,
		Email:     user.Email,
		Provider:  user.Provider.String,
		CreatedAt: user.CreatedAt.Time,
	}
}

func (s *Server) CreateUserHandler(ctx *gin.Context) {

	var req createUserParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed := util.HashPassword(req.Password)

	arg := db.CreateUserParams{

		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: pgtype.Text{String: hashed, Valid: true},
		Provider:     pgtype.Text{String: req.Provider, Valid: true},
	}

	user, err := s.queries.CreateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (s *Server) LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.queries.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.PasswordHash.String)
	fmt.Println(req.Password)
	fmt.Println(user.PasswordHash.String)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	expiresAt := pgtype.Timestamp{
		Time:  refreshPayload.ExpiresAt.Time,
		Valid: true,
	}

	pgUUID := pgtype.UUID{
		Bytes: refreshPayload.ID, // uuid.UUID is a [16]byte, which matches
		Valid: true,
	}
	session, err := s.queries.CreateSession(ctx, db.CreateSessionParams{
		ID:           pgUUID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             uuid.UUID(session.ID.Bytes),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)

}
