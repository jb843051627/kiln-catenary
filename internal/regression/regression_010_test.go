package regression

import (
    "net/http"
    "net/http/httptest"
    "path/filepath"
    "testing"
    "github.com/jb843051627/kiln-catenary/internal/handler"
    "github.com/jb843051627/kiln-catenary/internal/service"
    "github.com/jb843051627/kiln-catenary/internal/store"
)

func TestBug10_MissingEventReturnsNotFound(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "bug10.db")); if err != nil { t.Fatal(err) }
	app := service.NewApp(db); t.Cleanup(func() { app.Close(); db.Close() })
	recorder := httptest.NewRecorder(); request := httptest.NewRequest(http.MethodPost, "/api/events/missing-event", nil)
	handler.New(app).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotFound { t.Fatalf("missing event returned HTTP %d", recorder.Code) }
}
