package render

import (
    "net/http/httptest"
    "reflect"
    "testing"
)

func TestContentTypeHeaderIsolationProbe(t *testing.T) {
    first := httptest.NewRecorder()
    JSON{}.WriteContentType(first)
    original := first.Header().Get("Content-Type")
    first.Header()["Content-Type"][0] = "application/custom"
    defer func() { first.Header()["Content-Type"][0] = original }()
    second := httptest.NewRecorder()
    JSON{}.WriteContentType(second)
    if got := second.Header().Get("Content-Type"); got != original {
        t.Fatalf("independent JSON response inherited another response header: got %q, want %q", got, original)
    }
}

func TestContentTypeExistingHeaderProbe(t *testing.T) {
    w := httptest.NewRecorder()
    expected := []string{"application/custom", "fallback"}
    w.Header()["Content-Type"] = expected
    JSON{}.WriteContentType(w)
    if !reflect.DeepEqual(w.Header()["Content-Type"], expected) {
        t.Fatalf("existing Content-Type was changed: %v", w.Header()["Content-Type"])
    }
}

func BenchmarkContentTypeHeaderIsolation(b *testing.B) {
    w := httptest.NewRecorder()
    header := w.Header()
    b.ReportAllocs()
    b.ResetTimer()
    for b.Loop() {
        delete(header, "Content-Type")
        JSON{}.WriteContentType(w)
    }
}
