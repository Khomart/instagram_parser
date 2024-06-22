package server

import (
	"log/slog"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/khomart/instagram_recipe_parser/internal/interfaces"
)

type RequestHandler struct {
	downloader interfaces.Downloader
	parser     interfaces.Parser
}

func NewRequestHandler(d interfaces.Downloader, p interfaces.Parser) *RequestHandler {
	return &RequestHandler{
		downloader: d,
		parser:     p,
	}
}

type URLParsingRequest struct {
	URL string `form:"URL" binding:"required" example:"https://www.instagram.com/p/C7h7ksJOBvc/"`
}

func (rh *RequestHandler) Test(ctx *gin.Context) {
	var req URLParsingRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, "test")
}

func (rh *RequestHandler) ParseURL(ctx *gin.Context) {
	var req URLParsingRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	m := regexp.MustCompile(`\.?([^.]*.com)`)
	if urlDomain := m.FindStringSubmatch(req.URL); urlDomain == nil || urlDomain[1] != "instagram.com" {
		slog.Error("Url domain mismatch", "domain", urlDomain)
		invalidURL(ctx)
		return
	}

	downloadFolder := rh.downloader.FetchPost(req.URL)
	summary, err := rh.parser.SummarizeContent(downloadFolder)
	if err != nil {
		slog.Error("Error summarizing content", "error", err)
		errorResponse(ctx, err)
		return
	}

	slog.Info("Recipe Instruction:")
	slog.Info(summary)
	successResponse(ctx, summary)
}
