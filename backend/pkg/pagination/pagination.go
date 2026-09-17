package pagination

import (
	"math"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

type Params struct {
	Page   int
	Limit  int
	Offset int
	Sort   []SortField
}

type SortField struct {
	Field string
	Desc  bool
}

type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
	Pages int   `json:"pages"`
}

func NewMeta(p Params, total int64) Meta {
	pages := 0
	if p.Limit > 0 {
		pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	}
	return Meta{
		Page:  p.Page,
		Limit: p.Limit,
		Total: total,
		Pages: pages,
	}
}

// Parse reads pagination params from query string.
// Allowed sort fields must be passed in to prevent SQL injection.
func Parse(c *gin.Context, allowedSorts map[string]string) Params {
	p := Params{
		Page:  DefaultPage,
		Limit: DefaultLimit,
	}

	if v := c.Query("page"); v != "" {
		if n, err := parseIntSafe(v); err == nil && n > 0 {
			p.Page = n
		}
	}
	if v := c.Query("limit"); v != "" {
		if n, err := parseIntSafe(v); err == nil && n > 0 && n <= MaxLimit {
			p.Limit = n
		}
	}
	p.Offset = (p.Page - 1) * p.Limit

	// sort=field:asc,field2:desc
	if v := c.Query("sort"); v != "" {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			field := part
			desc := false
			if strings.HasPrefix(part, "-") {
				field = part[1:]
				desc = true
			} else if i := strings.Index(part, ":"); i > 0 {
				field = part[:i]
				desc = strings.ToLower(part[i+1:]) == "desc"
			}
			if mapped, ok := allowedSorts[field]; ok {
				p.Sort = append(p.Sort, SortField{Field: mapped, Desc: desc})
			}
		}
	}
	return p
}

// OrderClause produces a SQL ORDER BY fragment, or "" if no sorting.
func (p Params) OrderClause() string {
	if len(p.Sort) == 0 {
		return ""
	}
	var parts []string
	for _, s := range p.Sort {
		if s.Desc {
			parts = append(parts, s.Field+" DESC")
		} else {
			parts = append(parts, s.Field+" ASC")
		}
	}
	return strings.Join(parts, ", ")
}

func parseIntSafe(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errNotNumber
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

var errNotNumber = &simpleErr{"not a number"}

type simpleErr struct{ s string }

func (e *simpleErr) Error() string { return e.s }
