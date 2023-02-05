package ginmidctx

import (
	"github.com/gin-gonic/gin"
	"github.com/gmaschi/log-exp-eval/pkg/tools/pagination"
)

const (
	paginationPageHeader  = "x-pagination-page"
	paginationPagesHeader = "x-pagination-pages"
	paginationItemsHeader = "x-pagination-items"
)

func SetPaginationHeader(ctx *gin.Context, pgData pagination.Pagination) {
	ctx.Header(paginationPageHeader, pgData.Page)
	ctx.Header(paginationPagesHeader, pgData.Pages)
	ctx.Header(paginationItemsHeader, pgData.Items)
}
