package web

import (
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/practice/webook/internal/service"
	"github.com/zmsocc/practice/webook/internal/web/ijwt"
	"net/http"
)

type ArticleHandler struct {
	svc service.ArticleService
}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	ag := server.Group("/articles")
	ag.POST("/edit", h.Edit)
	ag.POST("/publish", h.Publish)
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	claims, ok := ctx.MustGet("uc").(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	id, err := h.svc.Save(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
	}
	ctx.JSON(http.StatusOK, Result{Data: id})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
}
