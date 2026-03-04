package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/BroMikey/blog_backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type coinTransferRequest struct {
	FromCoinID int64  `json:"from_coin_id" binding:"required,min=1"`
	ToCoinID   int64  `json:"to_coin_id" binding:"required,min=1"`
	Amount     int64  `json:"amount" binding:"required,gt=0"`
	CoinType   string `json:"coin_type" binding:"required,oneof=penny nickel dime"`
}

func (server *Server) createCoinTransfer(ctx *gin.Context) {
	var req coinTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !server.validCoin(ctx, req.ToCoinID, req.CoinType) {
		return
	}

	arg := db.CoinTransferTxParams{
		FromCoinID: req.FromCoinID,
		ToCoinID:   req.ToCoinID,
		Amount:     req.Amount,
	}

	result, err := server.store.CoinTransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validCoin(ctx *gin.Context, coinid int64, cointype string) bool {
	coin, err := server.store.GetCoin(ctx, coinid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if coin.CoinType != cointype {
		err := fmt.Errorf("coin [%d] type missmatch %s vs %s", coinid, coin.CoinType, cointype)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
