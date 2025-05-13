package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Mandarinka0707/newRepoGOODarhit/internal/entity"
	"github.com/Mandarinka0707/newRepoGOODarhit/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *entity.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepo) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

type MockSessionRepo struct {
	mock.Mock
}

func (m *MockSessionRepo) CreateSession(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepo) GetSessionByToken(ctx context.Context, token string) (*entity.Session, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func setupTest(t *testing.T) (*AuthUsecase, *MockUserRepo, *MockSessionRepo) {
	userRepo := new(MockUserRepo)
	sessionRepo := new(MockSessionRepo)
	cfg := &auth.Config{
		TokenSecret:     "test-secret",
		TokenExpiration: time.Hour,
	}

	logger := zaptest.NewLogger(t)

	return NewAuthUsecase(userRepo, sessionRepo, cfg, logger), userRepo, sessionRepo
}

func TestGetUserByID_Success(t *testing.T) {
	uc, userRepo, _ := setupTest(t)
	ctx := context.Background()

	expectedUser := &entity.User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
	}

	userRepo.On("GetUserByID", ctx, int64(1)).Return(expectedUser, nil)

	user, err := uc.GetUserByID(ctx, 1) // Изменено: передаем int64 вместо строки

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	userRepo.AssertExpectations(t)
}

func TestRegister_Success(t *testing.T) {
	uc, userRepo, _ := setupTest(t)
	ctx := context.Background()

	req := &RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}

	userRepo.On("CreateUser", ctx, mock.AnythingOfType("*entity.User")).Return(int64(1), nil)

	resp, err := uc.Register(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.UserID)
	userRepo.AssertExpectations(t)
}

func TestRegister_HashError(t *testing.T) {
	uc, _, _ := setupTest(t)
	ctx := context.Background()

	req := &RegisterRequest{
		Username: "testuser",
		Password: string(make([]byte, 100)),
	}

	resp, err := uc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestRegister_DBError(t *testing.T) {
	uc, userRepo, _ := setupTest(t)
	ctx := context.Background()

	req := &RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}

	userRepo.On("CreateUser", ctx, mock.AnythingOfType("*entity.User")).Return(int64(0), errors.New("db error"))

	resp, err := uc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	uc, userRepo, sessionRepo := setupTest(t)
	ctx := context.Background()

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Password: string(hashedPassword),
		Role:     "user",
	}

	userRepo.On("GetUserByUsername", ctx, "testuser").Return(user, nil)
	sessionRepo.On("CreateSession", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)

	req := &LoginRequest{
		Username: "testuser",
		Password: password,
	}

	resp, err := uc.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "testuser", resp.Username)
	userRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	uc, userRepo, _ := setupTest(t)
	ctx := context.Background()

	userRepo.On("GetUserByUsername", ctx, "testuser").Return(nil, sql.ErrNoRows)

	req := &LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := uc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	uc, userRepo, _ := setupTest(t)
	ctx := context.Background()

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Password: "$2a$10$invalidhash",
		Role:     "user",
	}

	userRepo.On("GetUserByUsername", ctx, "testuser").Return(user, nil)

	req := &LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	resp, err := uc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "invalid username or password", err.Error())
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
}

func TestLogin_SessionCreationError(t *testing.T) {
	uc, userRepo, sessionRepo := setupTest(t)
	ctx := context.Background()

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Password: string(hashedPassword),
		Role:     "user",
	}

	userRepo.On("GetUserByUsername", ctx, "testuser").Return(user, nil)
	sessionRepo.On("CreateSession", ctx, mock.AnythingOfType("*entity.Session")).Return(errors.New("db error"))

	req := &LoginRequest{
		Username: "testuser",
		Password: password,
	}

	resp, err := uc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "internal server error", err.Error())
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	uc, userRepo, _ := setupTest(t)
	ctx := context.Background()

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
	}

	userRepo.On("GetUserByID", ctx, int64(1)).Return(user, nil)

	req := &GetUserRequest{UserID: 1}
	resp, err := uc.GetUser(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, user, resp.User)
	userRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	uc, userRepo, _ := setupTest(t)
	ctx := context.Background()

	userRepo.On("GetUserByID", ctx, int64(1)).Return(nil, nil)

	req := &GetUserRequest{UserID: 1}
	resp, err := uc.GetUser(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
}

func TestGetUser_DBError(t *testing.T) {
	uc, userRepo, _ := setupTest(t)
	ctx := context.Background()

	userRepo.On("GetUserByID", ctx, int64(1)).Return(nil, errors.New("db error"))

	req := &GetUserRequest{UserID: 1}
	resp, err := uc.GetUser(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
}

func TestGetUser_Logging(t *testing.T) {
	t.Run("Success logs info message", func(t *testing.T) {
		uc, userRepo, _ := setupTest(t)
		ctx := context.Background()

		user := &entity.User{
			ID:       1,
			Username: "testuser",
			Role:     "user",
		}

		userRepo.On("GetUserByID", ctx, int64(1)).Return(user, nil)

		req := &GetUserRequest{UserID: 1}
		resp, err := uc.GetUser(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, user, resp.User)
		userRepo.AssertExpectations(t)

	})

	t.Run("User not found logs error", func(t *testing.T) {
		uc, userRepo, _ := setupTest(t)
		ctx := context.Background()

		userRepo.On("GetUserByID", ctx, int64(1)).Return(nil, nil)

		req := &GetUserRequest{UserID: 1}
		resp, err := uc.GetUser(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.Nil(t, resp)
		userRepo.AssertExpectations(t)

	})

	t.Run("DB error logs error", func(t *testing.T) {
		uc, userRepo, _ := setupTest(t)
		ctx := context.Background()

		dbError := errors.New("database error")
		userRepo.On("GetUserByID", ctx, int64(1)).Return(nil, dbError)

		req := &GetUserRequest{UserID: 1}
		resp, err := uc.GetUser(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.Nil(t, resp)
		userRepo.AssertExpectations(t)

	})
}
func setupTestWithLogObserver(t *testing.T) (*AuthUsecase, *MockUserRepo, *MockSessionRepo, *observer.ObservedLogs) {
	userRepo := new(MockUserRepo)
	sessionRepo := new(MockSessionRepo)
	cfg := &auth.Config{
		TokenSecret:     "test-secret",
		TokenExpiration: time.Hour,
	}

	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	return NewAuthUsecase(userRepo, sessionRepo, cfg, logger), userRepo, sessionRepo, recorded
}

func TestGetUser_LoggingWithObserver(t *testing.T) {
	t.Run("Success logs info message", func(t *testing.T) {
		uc, userRepo, _, logs := setupTestWithLogObserver(t)
		ctx := context.Background()

		user := &entity.User{
			ID:       1,
			Username: "testuser",
			Role:     "user",
		}

		userRepo.On("GetUserByID", ctx, int64(1)).Return(user, nil)

		req := &GetUserRequest{UserID: 1}
		resp, err := uc.GetUser(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, user, resp.User)
		userRepo.AssertExpectations(t)

		assert.Equal(t, 1, logs.Len())
		assert.Contains(t, logs.All()[0].Message, "Get user request")
		assert.Equal(t, zap.InfoLevel, logs.All()[0].Level)
	})

	t.Run("User not found logs error", func(t *testing.T) {
		uc, userRepo, _, logs := setupTestWithLogObserver(t)
		ctx := context.Background()

		userRepo.On("GetUserByID", ctx, int64(1)).Return(nil, nil)

		req := &GetUserRequest{UserID: 1}
		resp, err := uc.GetUser(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.Nil(t, resp)
		userRepo.AssertExpectations(t)

		assert.Equal(t, 2, logs.Len())
		assert.Contains(t, logs.All()[1].Message, "User not found")
		assert.Equal(t, zap.ErrorLevel, logs.All()[1].Level)
	})

	t.Run("DB error logs error", func(t *testing.T) {
		uc, userRepo, _, logs := setupTestWithLogObserver(t)
		ctx := context.Background()

		dbError := errors.New("database error")
		userRepo.On("GetUserByID", ctx, int64(1)).Return(nil, dbError)

		req := &GetUserRequest{UserID: 1}
		resp, err := uc.GetUser(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		assert.Nil(t, resp)
		userRepo.AssertExpectations(t)

		assert.Equal(t, 2, logs.Len())
		assert.Contains(t, logs.All()[1].Message, "User not found")
		assert.Equal(t, zap.ErrorLevel, logs.All()[1].Level)
		assert.Equal(t, dbError, logs.All()[1].Context[0].Interface.(error))
	})
}

func TestGetUserByID_NotFound(t *testing.T) {
	uc, userRepo, _ := setupTest(t)
	ctx := context.Background()

	userRepo.On("GetUserByID", ctx, int64(1)).Return(nil, sql.ErrNoRows)

	user, err := uc.GetUserByID(ctx, 1) // Изменено: передаем int64 вместо строки

	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
	assert.Nil(t, user)
	userRepo.AssertExpectations(t)
}
