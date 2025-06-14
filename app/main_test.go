package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// testServer crée un serveur de test à partir de notre handler principal
func testServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/signup", signupHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/welcome", welcomeHandler)
	mux.HandleFunc("/trylater", tryLaterHandler)
	return mux
}

func TestIndexHandler(t *testing.T) {
	ts := httptest.NewServer(testServer())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestSignupMissingFields(t *testing.T) {
	ts := httptest.NewServer(testServer())
	defer ts.Close()

	resp, err := http.PostForm(ts.URL+"/signup", url.Values{
		"email":    {""},
		"username": {"testuser"},
		"pwd":      {"password"},
		"pwd2":     {"password"},
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

	resp, err := http.PostForm(ts.URL+"/signup", url.Values{
		"email":    {"test@example.com"},
		"username": {"testuser"},
		"pwd":      {"password123"},
		"pwd2":     {"password123"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", resp.StatusCode)
	}
}

func TestLoginFailWrongPassword(t *testing.T) {
	ts := httptest.NewServer(testServer())
	defer ts.Close()

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

	resp, err := http.PostForm(ts.URL+"/login", url.Values{
		"username": {"testuser"},
		"pwd":      {"password123"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("expected 303 on successful login, got %d", resp.StatusCode)
	}
}

func TestLogout(t *testing.T) {
	ts := httptest.NewServer(testServer())
	defer ts.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL+"/logout", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", resp.StatusCode)
	}
}
