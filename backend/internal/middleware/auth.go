package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/apperror"
	"github.com/campusx/api/pkg/jwt"
	"github.com/campusx/api/pkg/response"
)

const (
	ctxUserID    = "auth.user_id"
	ctxUserRole  = "auth.user_role"
	ctxCollegeID = "auth.college_id"
)

// RequireAuth validates the access token and injects user context.
func RequireAuth(issuer *jwt.Issuer) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			abort(c, apperror.Unauthorized("missing bearer token"))
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := issuer.ParseAccess(token)
		if err != nil {
			abort(c, apperror.Unauthorized("invalid or expired token"))
			return
		}
		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxUserRole, claims.Role)
		if claims.CollegeID != "" {
			c.Set(ctxCollegeID, claims.CollegeID)
		}
		c.Next()
	}
}

// RequireRole enforces that the user has one of the allowed roles.
func RequireRole(roles ...models.Role) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[string(r)] = struct{}{}
	}
	return func(c *gin.Context) {
		role, ok := c.Get(ctxUserRole)
		if !ok {
			abort(c, apperror.Unauthorized("missing auth context"))
			return
		}
		if _, ok := allowed[role.(string)]; !ok {
			abort(c, apperror.Forbidden("insufficient permissions"))
			return
		}
		c.Next()
	}
}

// RequireTenant enforces that the request targets the user's own college.
// Super admins bypass (explicitly).
// Reads target college from path (":college_id") or query ("college_id").
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ctxUserRole)
		if role == string(models.RoleSuperAdmin) {
			c.Next()
			return
		}

		collegeID, ok := c.Get(ctxCollegeID)
		if !ok || collegeID == "" {
			abort(c, apperror.Forbidden("no college context"))
			return
		}

		target := c.Param("college_id")
		if target == "" {
			target = c.Query("college_id")
		}
		if target != "" && target != collegeID {
			abort(c, apperror.Forbidden("cross-tenant access denied"))
			return
		}
		c.Next()
	}
}

// --- Accessors ---

func MustUserID(c *gin.Context) string {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func MustUserRole(c *gin.Context) models.Role {
	v, _ := c.Get(ctxUserRole)
	s, _ := v.(string)
	return models.Role(s)
}

func MustCollegeID(c *gin.Context) string {
	v, _ := c.Get(ctxCollegeID)
	s, _ := v.(string)
	return s
}

func abort(c *gin.Context, err error) {
	response.Fail(c, err)
	c.Abort()
}

// RequireTenantParam enforces that the :college_id path param matches the
// authenticated user's college. SUPER_ADMIN bypasses.
func RequireTenantParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ctxUserRole)
		if role == string(models.RoleSuperAdmin) {
			c.Next()
			return
		}

		userCollege, ok := c.Get(ctxCollegeID)
		if !ok || userCollege == "" {
			abort(c, apperror.Forbidden("no college context"))
			return
		}

		target := c.Param("id")
		if target == "" {
			abort(c, apperror.BadRequest("missing college_id"))
			return
		}
		if target != userCollege.(string) {
			abort(c, apperror.Forbidden("cross-tenant access denied"))
			return
		}
		c.Next()
	}
}
