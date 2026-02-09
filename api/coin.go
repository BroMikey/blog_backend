package api

import (
	"database/sql"
	"net/http"

	db "github.com/BroMikey/blog_backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createCoinRequest struct {
	Uid      int64  `json:"uid" biding:"required"`
	CoinType string `json:"coin_type" biding:"required, oneof=penny nickel dime"`
}

func (server *Server) createCoin(ctx *gin.Context) {
	var req createCoinRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCoinParams{
		Uid:      req.Uid,
		CoinType: req.CoinType,
		Balance:  0,
	}

	coin, err := server.store.CreateCoin(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, coin)
}

type getCoinRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCoin(ctx *gin.Context) {
	var req getCoinRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	coin, err := server.store.GetCoin(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, coin)
}

type listCoinRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCoin(ctx *gin.Context) {
	var req listCoinRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCoinsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	coins, err := server.store.ListCoins(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, coins)
}
