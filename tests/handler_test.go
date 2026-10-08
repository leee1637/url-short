package tests

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joho/godotenv"

	"url-short/internal/config"
	"url-short/internal/domain"
	"url-short/internal/handler"
	"url-short/internal/service"
)

// креды берём из того же .env, что и сервис — без хардкода
var (
	authUser string
	authPass string
)

func TestMain(m *testing.M) {
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		panic("не удалось загрузить ../.env: " + err.Error())
	}

	authUser = os.Getenv("AUTH_USER")
	authPass = os.Getenv("AUTH_PASS")
	if authUser == "" || authPass == "" {
		panic("AUTH_USER/AUTH_PASS пустые в .env")
	}

	os.Exit(m.Run())
}

// fakeRepo — заглушка репозитория, с помощью которой мы подменяем ошибки БД
type fakeRepo struct {
	saveErr    error
	getErr     error
	getResult  *domain.URL
	deleteErr  error
	updateErr  error
	savedAlias string
	savedURL   string
}

func (f *fakeRepo) Save(_ context.Context, u *domain.URL) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.savedAlias = u.Alias
	f.savedURL = u.OriginalURL
	return nil
}

func (f *fakeRepo) GetByAlias(_ context.Context, _ string) (*domain.URL, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.getResult, nil
}

func (f *fakeRepo) DeleteByAlias(_ context.Context, _ string) error {
	return f.deleteErr
}

func (f *fakeRepo) UpdateURLByAlias(_ context.Context, _, _ string) error {
	return f.updateErr
}

func newRouter(t *testing.T, repo domain.URLRepository) http.Handler {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return handler.New(service.New(repo, logger)).Router(
		config.Auth{User: authUser, Pass: authPass}, logger)
}

func do(t *testing.T, r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.SetBasicAuth(authUser, authPass)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestSave_Success(t *testing.T) {
	repo := &fakeRepo{}
	w := do(t, newRouter(t, repo), http.MethodPost, "/urls",
		`{"url":"https://example.com","alias":"myalias"}`)

	if w.Code != http.StatusCreated {
		t.Fatalf("got %d, want 201, body=%s", w.Code, w.Body)
	}
	if repo.savedAlias != "myalias" || repo.savedURL != "https://example.com" {
		t.Fatalf("repo записан неверно: alias=%q url=%q", repo.savedAlias, repo.savedURL)
	}
}

func TestSave_GenerateAlias(t *testing.T) {
	repo := &fakeRepo{}
	// alias не указан — сервис генерирует его сам
	w := do(t, newRouter(t, repo), http.MethodPost, "/urls",
		`{"url":"https://example.com"}`)

	if w.Code != http.StatusCreated {
		t.Fatalf("got %d, want 201, body=%s", w.Code, w.Body)
	}
	if len(repo.savedAlias) != 8 {
		t.Fatalf("сгенерированный alias = %q, want 8 символов", repo.savedAlias)
	}
}

func TestSave_BadRequest(t *testing.T) {
	// реально битый JSON — не валидный JSON с пустым полем
	w := do(t, newRouter(t, &fakeRepo{}), http.MethodPost, "/urls",
		`{"url": `)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400, body=%s", w.Code, w.Body)
	}
}

func TestSave_InvalidURL(t *testing.T) {
	// валидный JSON с невалидным URL (не http/https)
	w := do(t, newRouter(t, &fakeRepo{}), http.MethodPost, "/urls",
		`{"url":"ftp://example.com","alias":"myalias"}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400, body=%s", w.Code, w.Body)
	}
}

func TestSave_Conflict(t *testing.T) {
	// alias занят — storage вернул ErrConflict, ждём 409, не 400
	repo := &fakeRepo{saveErr: domain.ErrConflict}
	w := do(t, newRouter(t, repo), http.MethodPost, "/urls",
		`{"url":"https://example.com","alias":"taken22"}`)

	if w.Code != http.StatusConflict {
		t.Fatalf("got %d, want 409, body=%s", w.Code, w.Body)
	}
}

func TestSave_DBError(t *testing.T) {
	// внутренняя ошибка БД — 500, причём детали наружу не утекают
	repo := &fakeRepo{saveErr: context.DeadlineExceeded}
	w := do(t, newRouter(t, repo), http.MethodPost, "/urls",
		`{"url":"https://example.com","alias":"myalias"}`)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("got %d, want 500, body=%s", w.Code, w.Body)
	}
	if strings.Contains(w.Body.String(), "deadline") {
		t.Fatalf("детали внутренней ошибки утекли наружу: %s", w.Body)
	}
}

func TestNoAuth(t *testing.T) {
	// запрос без кредов — 401
	r := newRouter(t, &fakeRepo{})
	req := httptest.NewRequest(http.MethodPost, "/urls", strings.NewReader(`{}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("got %d, want 401", w.Code)
	}
	if w.Header().Get("WWW-Authenticate") == "" {
		t.Fatal("отсутствует заголовок WWW-Authenticate")
	}
}

func TestWrongCreds(t *testing.T) {
	// неверные креды — 401
	r := newRouter(t, &fakeRepo{})
	req := httptest.NewRequest(http.MethodPost, "/urls", strings.NewReader(`{}`))
	req.SetBasicAuth("admin", "wrongpass")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("got %d, want 401", w.Code)
	}
}

func TestRedirect_Success(t *testing.T) {
	// публичная ручка, без auth, редирект на исходный URL
	repo := &fakeRepo{getResult: &domain.URL{OriginalURL: "https://example.com"}}
	req := httptest.NewRequest(http.MethodGet, "/myalias", nil)
	w := httptest.NewRecorder()
	newRouter(t, repo).ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("got %d, want 302", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "https://example.com" {
		t.Fatalf("Location = %q, want https://example.com", loc)
	}
}

func TestRedirect_NotFound(t *testing.T) {
	// alias не найден — 404, а не internal error
	repo := &fakeRepo{getErr: domain.ErrNotFound}
	req := httptest.NewRequest(http.MethodGet, "/missing1", nil)
	w := httptest.NewRecorder()
	newRouter(t, repo).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("got %d, want 404, body=%s", w.Code, w.Body)
	}
}

func TestDelete_Success(t *testing.T) {
	w := do(t, newRouter(t, &fakeRepo{}), http.MethodDelete, "/urls/myalias", "")

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want 200, body=%s", w.Code, w.Body)
	}
}

func TestDelete_NotFound(t *testing.T) {
	repo := &fakeRepo{deleteErr: domain.ErrNotFound}
	w := do(t, newRouter(t, repo), http.MethodDelete, "/urls/missing1", "")

	if w.Code != http.StatusNotFound {
		t.Fatalf("got %d, want 404, body=%s", w.Code, w.Body)
	}
}

func TestUpdate_Success(t *testing.T) {
	repo := &fakeRepo{}
	w := do(t, newRouter(t, repo), http.MethodPatch, "/urls/myalias",
		`{"url":"https://new.example.com"}`)

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want 200, body=%s", w.Code, w.Body)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repo := &fakeRepo{updateErr: domain.ErrNotFound}
	w := do(t, newRouter(t, repo), http.MethodPatch, "/urls/missing1",
		`{"url":"https://new.example.com"}`)

	if w.Code != http.StatusNotFound {
		t.Fatalf("got %d, want 404, body=%s", w.Code, w.Body)
	}
}

func TestUpdate_BadRequest(t *testing.T) {
	// битый JSON в update тоже даёт 400
	w := do(t, newRouter(t, &fakeRepo{}), http.MethodPatch, "/urls/myalias",
		`{"url": `)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400, body=%s", w.Code, w.Body)
	}
}

func TestGetAlias_Success(t *testing.T) {
	repo := &fakeRepo{getResult: &domain.URL{
		ID:          1,
		OriginalURL: "https://example.com",
		Alias:       "myalias",
	}}
	w := do(t, newRouter(t, repo), http.MethodGet, "/urls/myalias", "")

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want 200, body=%s", w.Code, w.Body)
	}
	if !strings.Contains(w.Body.String(), "https://example.com") {
		t.Fatalf("тело не содержит original_url: %s", w.Body)
	}
}
