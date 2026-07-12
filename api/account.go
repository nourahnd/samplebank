package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/nourahnd/samplebank/db/sqlc"
)

type createAccountReq struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidID))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternal))
		return
	}

	ctx.JSON(http.StatusCreated, account)

}

type accountIDUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req accountIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidID))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(ErrAccountNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternal))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getListOfAccountsReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getListOfAccounts(ctx *gin.Context) {
	var req getListOfAccountsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidReq))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternal))
		return
	}
	ctx.JSON(http.StatusOK, accounts)

}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req accountIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidID))
		return
	}

	rows, err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				ctx.JSON(http.StatusConflict, errorResponse(ErrAccountHasTransfers))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternal))
		return
	}

	if rows == 0 {
		ctx.JSON(http.StatusNotFound, errorResponse(ErrAccountNotFound))
		return
	}

	ctx.Status(http.StatusNoContent)
}
