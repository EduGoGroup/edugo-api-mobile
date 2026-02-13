package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Tests para GetUserIDFromContext
func TestGetUserIDFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedID := uuid.New()
	c.Set(ContextKeyUserID, expectedID.String())

	userID, ok := GetUserIDFromContext(c)

	assert.True(t, ok)
	assert.Equal(t, expectedID, userID)
}

func TestGetUserIDFromContext_NotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	userID, ok := GetUserIDFromContext(c)

	assert.False(t, ok)
	assert.Equal(t, uuid.Nil, userID)
}

func TestGetUserIDFromContext_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyUserID, 12345) // tipo incorrecto

	userID, ok := GetUserIDFromContext(c)

	assert.False(t, ok)
	assert.Equal(t, uuid.Nil, userID)
}

func TestGetUserIDFromContext_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyUserID, "not-a-uuid")

	userID, ok := GetUserIDFromContext(c)

	assert.False(t, ok)
	assert.Equal(t, uuid.Nil, userID)
}

// Tests para MustGetUserIDFromContext
func TestMustGetUserIDFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedID := uuid.New()
	c.Set(ContextKeyUserID, expectedID.String())

	userID := MustGetUserIDFromContext(c)

	assert.Equal(t, expectedID, userID)
}

func TestMustGetUserIDFromContext_Panic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	assert.Panics(t, func() {
		MustGetUserIDFromContext(c)
	})
}

// Tests para GetSchoolIDFromContext
func TestGetSchoolIDFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedID := uuid.New()
	c.Set(ContextKeySchoolID, expectedID.String())

	schoolID, ok := GetSchoolIDFromContext(c)

	assert.True(t, ok)
	assert.Equal(t, expectedID, schoolID)
}

func TestGetSchoolIDFromContext_NotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	schoolID, ok := GetSchoolIDFromContext(c)

	assert.False(t, ok)
	assert.Equal(t, uuid.Nil, schoolID)
}

func TestGetSchoolIDFromContext_EmptyString(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeySchoolID, "")

	schoolID, ok := GetSchoolIDFromContext(c)

	assert.False(t, ok)
	assert.Equal(t, uuid.Nil, schoolID)
}

func TestGetSchoolIDFromContext_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeySchoolID, 12345)

	schoolID, ok := GetSchoolIDFromContext(c)

	assert.False(t, ok)
	assert.Equal(t, uuid.Nil, schoolID)
}

func TestGetSchoolIDFromContext_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeySchoolID, "not-a-uuid")

	schoolID, ok := GetSchoolIDFromContext(c)

	assert.False(t, ok)
	assert.Equal(t, uuid.Nil, schoolID)
}

// Tests para MustGetSchoolIDFromContext
func TestMustGetSchoolIDFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedID := uuid.New()
	c.Set(ContextKeySchoolID, expectedID.String())

	schoolID := MustGetSchoolIDFromContext(c)

	assert.Equal(t, expectedID, schoolID)
}

func TestMustGetSchoolIDFromContext_Panic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	assert.Panics(t, func() {
		MustGetSchoolIDFromContext(c)
	})
}

// Tests para GetRoleFromContext
func TestGetRoleFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyRole, "teacher")

	role, ok := GetRoleFromContext(c)

	assert.True(t, ok)
	assert.Equal(t, "teacher", role)
}

func TestGetRoleFromContext_NotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	role, ok := GetRoleFromContext(c)

	assert.False(t, ok)
	assert.Equal(t, "", role)
}

func TestGetRoleFromContext_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyRole, 12345)

	role, ok := GetRoleFromContext(c)

	assert.False(t, ok)
	assert.Equal(t, "", role)
}

// Tests para GetEmailFromContext
func TestGetEmailFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyEmail, "test@example.com")

	email, ok := GetEmailFromContext(c)

	assert.True(t, ok)
	assert.Equal(t, "test@example.com", email)
}

func TestGetEmailFromContext_NotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	email, ok := GetEmailFromContext(c)

	assert.False(t, ok)
	assert.Equal(t, "", email)
}

func TestGetEmailFromContext_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyEmail, 12345)

	email, ok := GetEmailFromContext(c)

	assert.False(t, ok)
	assert.Equal(t, "", email)
}

// Tests para IsAdminRole
func TestIsAdminRole_Admin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyRole, "admin")

	assert.True(t, IsAdminRole(c))
}

func TestIsAdminRole_SuperAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyRole, "super_admin")

	assert.True(t, IsAdminRole(c))
}

func TestIsAdminRole_NotAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyRole, "teacher")

	assert.False(t, IsAdminRole(c))
}

func TestIsAdminRole_NoRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	assert.False(t, IsAdminRole(c))
}

// Tests para HasRole
func TestHasRole_SingleRoleMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyRole, "teacher")

	assert.True(t, HasRole(c, "teacher"))
}

func TestHasRole_MultipleRolesMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyRole, "teacher")

	assert.True(t, HasRole(c, "admin", "teacher", "student"))
}

func TestHasRole_NoMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyRole, "student")

	assert.False(t, HasRole(c, "admin", "teacher"))
}

func TestHasRole_NoRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	assert.False(t, HasRole(c, "admin", "teacher"))
}

func TestHasRole_EmptyRolesList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	c.Set(ContextKeyRole, "teacher")

	assert.False(t, HasRole(c))
}
