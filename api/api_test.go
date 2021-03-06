package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/eiko-team/eiko/misc/data"
	"github.com/eiko-team/eiko/misc/misc"

	"github.com/julienschmidt/httprouter"
)

var (
	token, _ = misc.UserToToken(data.UserTest)
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func TestGet400(t *testing.T) {
	router := httprouter.New()
	ServeFiles(router)

	tests := []struct {
		name        string
		method      string
		URL         string
		code        int
		invalid     bool
		contentType string
	}{
		{"root", "GET", "/", http.StatusOK, false, "text/html"},
		{"index", "GET", "/index.html", http.StatusOK, false, "text/html"},
		{"login", "GET", "/login.html", http.StatusOK, false, "text/html"},
		{"app", "GET", "/eiko.html", http.StatusOK, false, "text/html"},
		{"service_worker", "GET", "/eiko-sw.js", http.StatusOK, false, "application/javascript; charset=utf-8"},
		{"favicon_Eiko", "GET", "/EIKO.ico", http.StatusOK, false, "image/vnd.microsoft.icon; charset=utf-8"},
		{"favicon", "GET", "/favicon.ico", http.StatusOK, false, "image/vnd.microsoft.icon; charset=utf-8"},
		{"manifest", "GET", "/manifest.json", http.StatusOK, false, "application/json; charset=utf-8"},
		{"js lib", "GET", "/js/lib.js", http.StatusOK, false, "application/javascript"},
		{"js color", "GET", "/js/color.js", http.StatusOK, false, "application/javascript"},
		{"json autocomplete_data", "GET", "/json/autocomplete_data.json", http.StatusOK, false, "application/json; charset=utf-8"},
		{"json manifest", "GET", "/json/manifest.json", http.StatusOK, false, "application/json; charset=utf-8"},
		{"not existing", "GET", "/blabla", 404, false, "text/plain; charset=utf-8"},
		{"no path", "GET", "/manifest.json", 500, true, "application/json"},
		{"no path", "GET", "/login.html", 500, false, "application/json"},
		{"no path", "GET", "/index.html", 500, false, "application/json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "no path" {
				Path = ""
			}
			t.Logf("%s", Path)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.URL, nil)
			router.ServeHTTP(w, req)
			if w.Code != tt.code {
				t.Errorf("Server error = %d, want %d for %s",
					w.Code, tt.code, tt.URL)
			}
			body := w.Body.String()
			cType := w.Header()["Content-Type"][0]
			if tt.contentType == "application/javascript" {
				if cType != "application/javascript" &&
					cType != "application/x-javascript" {
					t.Errorf("Content-Type = '%s' want '%s'",
						cType, tt.contentType)
				}
			} else if cType != tt.contentType {
				t.Logf("body: '%s'", body)
				t.Errorf("Content-Type = '%s' want '%s'",
					cType, tt.contentType)
			}
			if (body == "{\"error\":\"invalid_file\"}\n") != tt.invalid {
				t.Errorf("%+v", body)
			}
		})
	}
}

func TestWrapperFunction(t *testing.T) {
	router := InitAPI()
	fun := Functions[0]
	t.Run("TestWrapperFunction without body", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/api"+fun.Path, nil)
		router.ServeHTTP(w, r)
		if w.Code != 500 {
			t.Errorf("TestWrapperFunction = %d, want %d",
				w.Code, 500)
		}
	})
	t.Run("TestWrapperFunction with body", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := "{\"message\":\"message\"}"
		r, _ := http.NewRequest("POST", "/api"+fun.Path,
			strings.NewReader(body))
		router.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("TestWrapperFunction = %d, want %d",
				w.Code, http.StatusOK)
		}
	})
}

func TestWrapperFunctionCookie(t *testing.T) {
	router := InitAPI()
	fun := FunctionsWithToken[0]

	tests := []struct {
		name  string
		token string
		code  int
		err   string
	}{
		{"good token", token, 200, ""},
		{"no token", "", 500, "no_token_found"},
		{"invalid token", " ", 500, "token_invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/api"+fun.Path, nil)
			if tt.token != "" {
				r.Header.Set("Cookie", "token="+tt.token)
			}
			t.Log(misc.DumpRequest(r))
			router.ServeHTTP(w, r)
			if w.Code != tt.code {
				t.Errorf("WrapperFunctionCookie = %d, want %d",
					w.Code, tt.code)
			}

			body := w.Body.String()
			err := fmt.Sprintf("{\"error\":\"%s\"}\n", tt.err)
			if (body == err) != (tt.err != "") {
				t.Errorf("WrapperFunctionCookie = '%s', don't want '%s'",
					body, err)
			}

		})
	}
}

func TestExecuteAPI(t *testing.T) {
	os.Setenv("PROJECT_ID", "api_test.go")
	ExecuteAPI()
	os.Setenv("FILE_TYPE", "test")
	ExecuteAPI()
}
