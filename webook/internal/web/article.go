package web

import (
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/zmsocc/practice/webook/internal/domain"
	"github.com/zmsocc/practice/webook/internal/service"
	"github.com/zmsocc/practice/webook/internal/web/ijwt"
	"github.com/zmsocc/practice/webook/pkg/ginx"
	"github.com/zmsocc/practice/webook/pkg/logger"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	svc     service.ArticleService
	l       logger.Logger
	intrSvc service.InteractiveService
	biz     string
}

func NewArticleHandler(svc service.ArticleService, l logger.Logger,
	intrSvc service.InteractiveService) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		l:       l,
		intrSvc: intrSvc,
		biz:     "article",
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	ag := server.Group("/articles")
	ag.POST("/edit", h.Edit)
	ag.POST("/publish", h.Publish)
	ag.POST("/withdraw", h.Withdraw)
	ag.GET("/detail/:id", ginx.WrapBody(h.Detail))
	ag.POST("/list", ginx.WrapBody(h.List))

	pub := server.Group("/pub")
	pub.GET("/:id", h.PubDetail)
	pub.POST("/like", ginx.WrapBody(h.Like))
	pub.POST("/collect", ginx.WrapBody(h.Collect))
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	claims, ok := ctx.MustGet("users").(*ijwt.UserClaims)
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
	claims, ok := ctx.MustGet("users").(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	id, err := h.svc.Publish(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Data: id, Msg: "发表成功"})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	claims, ok := ctx.MustGet("users").(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		h.l.Error("未发现用户的 session 信息")
		return
	}
	err := h.svc.Withdraw(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		h.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "成功设置为仅自己可见"})
}

func (h *ArticleHandler) Detail(ctx *gin.Context) (Result, error) {
	var uc ijwt.UserClaims
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return Result{Code: 4, Msg: "参数错误"}, nil
	}
	data, err := h.svc.GetById(ctx, id)
	if err != nil {
		return Result{Code: 5, Msg: "系统错误"}, nil
	}
	if data.Author.Id != uc.Uid {
		return Result{Code: 4, Msg: "输入有误"},
			fmt.Errorf("非法访问文章，创作者 ID 不匹配 %d", uc.Uid)
	}
	return Result{
		Data: ArticleVO{
			Id:      data.Id,
			Title:   data.Title,
			Content: data.Content,
			Status:  data.Status.ToUint8(),
			Ctime:   data.Ctime.Format(time.DateTime),
			Utime:   data.Utime.Format(time.DateTime),
		},
	}, nil
}

func (h *ArticleHandler) List(ctx *gin.Context) (Result, error) {
	var (
		req ListReq
		uc  ijwt.UserClaims
	)
	res, err := h.svc.List(ctx, uc.Uid, req.Offset, req.Limit)
	if err != nil {
		return Result{Code: 5, Msg: "系统错误"}, nil
	}
	return Result{
		Data: slice.Map[domain.Article, ArticleVO](res, func(idx int, src domain.Article) ArticleVO {
			return ArticleVO{
				Id:       src.Id,
				Title:    src.Title,
				Abstract: src.Abstract(),
				// 这个列表请求，不需要返回内容
				//Content: src.Content,
				// 这个是创作者看自己的文章列表，也不需要这个字段
				//Author: src.Author,
				Status: src.Status.ToUint8(),
				Ctime:  src.Ctime.Format(time.DateTime),
				Utime:  src.Utime.Format(time.DateTime),
			}
		}),
	}, nil
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "参数错误"})
		h.l.Error("前端输入的 ID 不对", logger.Error(err))
		return
	}
	var eg errgroup.Group
	var art domain.Article
	uc := ctx.MustGet("users").(*ijwt.UserClaims)
	eg.Go(func() error {
		art, err = h.svc.GetPubById(ctx, id, uc.Uid)
		return err
	})
	err = eg.Wait()
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	go func() {
		er := h.intrSvc.IncrReadCnt(ctx, h.biz, art.Id)
		if er != nil {
			h.l.Error("增加阅读计数失败")
		}
	}()

	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVO{
			Id:      art.Id,
			Title:   art.Title,
			Content: art.Content,
			Status:  art.Status.ToUint8(),
			Author:  art.Author.Name,
			Ctime:   art.Ctime.Format(time.DateTime),
			Utime:   art.Utime.Format(time.DateTime),
		},
	})
}

func (h *ArticleHandler) Like(ctx *gin.Context) (Result, error) {
	var err error
	var req LikeReq
	var uc ijwt.UserClaims
	if req.Like {
		err = h.intrSvc.Like(ctx, h.biz, req.Id, uc.Uid)
	} else {
		err = h.intrSvc.CancelLike(ctx, h.biz, req.Id, uc.Uid)
	}
	if err != nil {
		return Result{Code: 5, Msg: "系统错误"}, nil
	}
	return Result{Msg: "ok"}, nil
}

func (h *ArticleHandler) Collect(ctx *gin.Context) (Result, error) {
	var err error
	var req CollectReq
	var uc ijwt.UserClaims
	if req.Collect {
		err = h.intrSvc.Collect(ctx, h.biz, req.Id, uc.Uid)
	} else {
		err = h.intrSvc.CancelCollect(ctx, h.biz, req.Id, uc.Uid)
	}
	if err != nil {
		return Result{Code: 5, Msg: "系统错误"}, nil
	}
	return Result{Msg: "ok"}, nil
}
