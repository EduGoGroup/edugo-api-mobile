package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// Tests para GetActiveContext
func TestGetActiveContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedCtx := &auth.UserContext{
		RoleID:      "role-123",
		RoleName:    "teacher",
		SchoolID:    "school-123",
		Permissions: []string{"read", "write"},
	}
	c.Set(ContextKeyActiveContext, expectedCtx)

	ctx := GetActiveContext(c)

	assert.NotNil(t, ctx)
	assert.Equal(t, expectedCtx, ctx)
}

func TestGetActiveContext_NotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	ctx := GetActiveContext(c)

	assert.Nil(t, ctx)
}

func TestGetActiveContext_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyActiveContext, "invalid-type")

	ctx := GetActiveContext(c)

	assert.Nil(t, ctx)
}

// Tests para RequirePermission
func TestRequirePermission_HasPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyActiveContext, &auth.UserContext{
		Permissions: []string{"materials:read", "materials:write"},
	})

	middleware := RequirePermission(enum.PermissionMaterialsRead)
	middleware(c)

	assert.False(t, c.IsAborted())
}

func TestRequirePermission_NoPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyActiveContext, &auth.UserContext{
		Permissions: []string{"materials:read"},
	})

	middleware := RequirePermission(enum.PermissionMaterialsUpdate)
	middleware(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "No tiene el permiso requerido")
}

func TestRequirePermission_NoActiveContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	middleware := RequirePermission(enum.PermissionMaterialsRead)
	middleware(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Se requiere contexto RBAC activo")
}

// Tests para RequireAnyPermission
func TestRequireAnyPermission_HasOnePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyActiveContext, &auth.UserContext{
		Permissions: []string{"materials:read"},
	})

	middleware := RequireAnyPermission(enum.PermissionMaterialsRead, enum.PermissionMaterialsUpdate)
	middleware(c)

	assert.False(t, c.IsAborted())
}

func TestRequireAnyPermission_HasMultiplePermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyActiveContext, &auth.UserContext{
		Permissions: []string{"materials:read", "materials:write", "materials:delete"},
	})

	middleware := RequireAnyPermission(enum.PermissionMaterialsRead, enum.PermissionMaterialsUpdate)
	middleware(c)

	assert.False(t, c.IsAborted())
}

func TestRequireAnyPermission_NoPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyActiveContext, &auth.UserContext{
		Permissions: []string{"assessments:read"},
	})

	middleware := RequireAnyPermission(enum.PermissionMaterialsRead, enum.PermissionMaterialsUpdate)
	middleware(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "No tiene ninguno de los permisos requeridos")
}

func TestRequireAnyPermission_NoActiveContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	middleware := RequireAnyPermission(enum.PermissionMaterialsRead, enum.PermissionMaterialsUpdate)
	middleware(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Se requiere contexto RBAC activo")
}

// Tests para RequireAdmin (Deprecated)
func TestRequireAdmin_AsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyRole, "admin")

	middleware := RequireAdmin()
	middleware(c)

	assert.False(t, c.IsAborted())
}

func TestRequireAdmin_AsSuperAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyRole, "super_admin")

	middleware := RequireAdmin()
	middleware(c)

	assert.False(t, c.IsAborted())
}

func TestRequireAdmin_AsTeacher(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyRole, "teacher")

	middleware := RequireAdmin()
	middleware(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// Tests para RequireSuperAdmin (Deprecated)
func TestRequireSuperAdmin_AsSuperAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyRole, "super_admin")

	middleware := RequireSuperAdmin()
	middleware(c)

	assert.False(t, c.IsAborted())
}

func TestRequireSuperAdmin_AsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyRole, "admin")

	middleware := RequireSuperAdmin()
	middleware(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// Tests para RequireTeacher (Deprecated)
func TestRequireTeacher_AsTeacher(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyRole, "teacher")

	middleware := RequireTeacher()
	middleware(c)

	assert.False(t, c.IsAborted())
}

func TestRequireTeacher_AsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyRole, "admin")

	middleware := RequireTeacher()
	middleware(c)

	assert.False(t, c.IsAborted())
}

func TestRequireTeacher_AsStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyRole, "student")

	middleware := RequireTeacher()
	middleware(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// Tests para RequireStudentOrAbove (Deprecated)
func TestRequireStudentOrAbove_AsStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyRole, "student")

	middleware := RequireStudentOrAbove()
	middleware(c)

	assert.False(t, c.IsAborted())
}

func TestRequireStudentOrAbove_AsTeacher(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyRole, "teacher")

	middleware := RequireStudentOrAbove()
	middleware(c)

	assert.False(t, c.IsAborted())
}

func TestRequireStudentOrAbove_AsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set(ContextKeyRole, "admin")

	middleware := RequireStudentOrAbove()
	middleware(c)

	assert.False(t, c.IsAborted())
}

func TestRequireStudentOrAbove_NoRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	middleware := RequireStudentOrAbove()
	middleware(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code) // Sin role = Unauthorized
}
