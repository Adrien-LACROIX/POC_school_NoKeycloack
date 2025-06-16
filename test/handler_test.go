package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	h "pocnokc/internal/handler"
)

// testServer crée un serveur de test à partir de notre handler principal
func testServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.IndexHandler)
	mux.HandleFunc("/signup", h.SignupHandler)
	mux.HandleFunc("/login", h.LoginHandler)
	mux.HandleFunc("/logout", h.LogoutHandler)
	mux.HandleFunc("/welcome", h.WelcomeHandler)
	mux.HandleFunc("/trylater", h.TryLaterHandler)
	return mux
}

func TestIndexHandler(t *testing.T) {
	ts := httptest.NewServer(testServer())
	defer ts.Close()
	ts.URL = "http://localhost:8080"

	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestSignupFailed(t *testing.T) {
	ts := httptest.NewServer(testServer())
	defer ts.Close()
	ts.URL = "http://localhost:8080"

	resp, err := http.PostForm(ts.URL+"/signup", url.Values{
		"email":    {"test@example.com"},
		"username": {"testuser"},
		"pwd":      {"password"},
		"pwd2":     {"pwd"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for form error, got %d", resp.StatusCode)
	}
}

func TestSignupSuccess(t *testing.T) {
	ts := httptest.NewServer(testServer())
	defer ts.Close()
	ts.URL = "http://localhost:8080"

	resp, err := singUp(ts)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 redirect, got %d", resp.StatusCode)
	}
}

func TestLoginFailWrongPassword(t *testing.T) {
	ts := httptest.NewServer(testServer())
	defer ts.Close()
	ts.URL = "http://localhost:8080"

	resp, err := http.PostForm(ts.URL+"/login", url.Values{
		"username": {"testuser"},
		"pwd":      {"wrongpass"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 on failed login, got %d", resp.StatusCode)
	}
}

func TestLoginSuccess(t *testing.T) {
	ts := httptest.NewServer(testServer())
	defer ts.Close()
	ts.URL = "http://localhost:8080"

	resp, err := login(ts)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 on successful login, got %d", resp.StatusCode)
	}
}

func TestLogout(t *testing.T) {
	ts := httptest.NewServer(testServer())
	defer ts.Close()
	ts.URL = "http://localhost:8080"

	client := &http.Client{}
	_, err := login(ts)
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("GET", ts.URL+"/logout", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 redirect, got %d", resp.StatusCode)
	}
}

func singUp(ts *httptest.Server) (*http.Response, error) {
	return http.PostForm(ts.URL+"/signup", url.Values{
		"email":    {"test@example.com"},
		"username": {"testuser"},
		"pwd":      {"password123"},
		"pwd2":     {"password123"},
	})
}

func login(ts *httptest.Server) (*http.Response, error) {
	resp, err := singUp(ts)
	if err != nil {
		return resp, err
	}

	return http.PostForm(ts.URL+"/login", url.Values{
		"username": {"testuser"},
		"pwd":      {"password123"},
	})
}
