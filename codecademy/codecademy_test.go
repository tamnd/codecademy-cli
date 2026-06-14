package codecademy_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/codecademy-cli/codecademy"
)

func newTestClient(ts *httptest.Server) *codecademy.Client {
	cfg := codecademy.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return codecademy.NewClient(cfg)
}

func fakeCatalogHTML(buildID string) string {
	return fmt.Sprintf(`<html><head></head><body><script id="__NEXT_DATA__" type="application/json">{"buildId":%q,"props":{"pageProps":{"initialCatalogResults":{"totalResults":2,"totalPages":1,"pageItems":[{"slug":"learn-python-3","title":"Learn Python 3","type":"course","difficulty":"beginner","lessonCount":14,"pro":false,"grantsCertificate":true,"urlPath":"/learn/learn-python-3"},{"slug":"learn-go","title":"Learn Go","type":"course","difficulty":"beginner","lessonCount":10,"pro":true,"grantsCertificate":false,"urlPath":"/learn/learn-go"}]}}}}</script></body></html>`, buildID)
}

func TestListCourses(t *testing.T) {
	const buildID = "test-build-id-123"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/catalog":
			w.Write([]byte(fakeCatalogHTML(buildID)))
		case r.URL.Path == fmt.Sprintf("/_next/data/%s/catalog.json", buildID):
			data := map[string]any{
				"pageProps": map[string]any{
					"initialCatalogResults": map[string]any{
						"totalResults": 2,
						"totalPages":   1,
						"pageItems": []any{
							map[string]any{"slug": "learn-python-3", "title": "Learn Python 3", "type": "course", "lessonCount": 14, "grantsCertificate": true, "urlPath": "/learn/learn-python-3"},
							map[string]any{"slug": "learn-go", "title": "Learn Go", "type": "course", "lessonCount": 10, "grantsCertificate": false, "urlPath": "/learn/learn-go"},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(data)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	c := newTestClient(ts)
	courses, err := c.ListCourses(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(courses) != 2 {
		t.Fatalf("want 2 courses, got %d", len(courses))
	}
	if courses[0].Slug != "learn-python-3" {
		t.Errorf("first slug = %q, want learn-python-3", courses[0].Slug)
	}
	if !courses[0].Certificate {
		t.Error("first course should grant certificate")
	}
}

func TestSearch(t *testing.T) {
	const buildID = "test-build-id-456"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/catalog":
			w.Write([]byte(fakeCatalogHTML(buildID)))
		case r.URL.Path == fmt.Sprintf("/_next/data/%s/catalog.json", buildID):
			data := map[string]any{
				"pageProps": map[string]any{
					"initialCatalogResults": map[string]any{
						"totalResults": 2,
						"totalPages":   1,
						"pageItems": []any{
							map[string]any{"slug": "learn-python-3", "title": "Learn Python 3", "type": "course", "lessonCount": 14, "urlPath": "/learn/learn-python-3"},
							map[string]any{"slug": "learn-go", "title": "Learn Go", "type": "course", "lessonCount": 10, "urlPath": "/learn/learn-go"},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(data)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	c := newTestClient(ts)
	results, err := c.Search(context.Background(), "python", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("want 1 result for 'python', got %d", len(results))
	}
	if results[0].Slug != "learn-python-3" {
		t.Errorf("result slug = %q, want learn-python-3", results[0].Slug)
	}
}
