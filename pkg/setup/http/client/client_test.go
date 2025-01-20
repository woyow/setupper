package client

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	maxHeaderBytes               = 1 << 20
	readHeaderTimeout            = 1 * time.Second
	idleTimeout                  = 1 * time.Second
	readTimeout                  = 1 * time.Second
	writeTimeout                 = 1 * time.Second
	youHasBeenRedirectedResponse = "you has been redirected"
)

func prepareHandler(port string) http.Handler {
	handler := gin.New()

	handler.Use(gin.Recovery(), gin.Logger())

	redirectPath := "handle_redirect"

	handler.GET(redirectPath, func(c *gin.Context) {
		c.String(http.StatusOK, youHasBeenRedirectedResponse)
	})

	redirectLocation := fmt.Sprintf("http://0.0.0.0:%s/%s", port, redirectPath)

	handler.GET("/test_moved_permanently_redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, redirectLocation)
	})
	handler.GET("/test_found_redirect", func(c *gin.Context) {
		c.Redirect(http.StatusFound, redirectLocation)
	})
	handler.GET("/test_see_other_redirect", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, redirectLocation)
	})
	handler.GET("/test_temporary_redirect", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, redirectLocation)
	})
	handler.GET("/test_permanent_redirect", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, redirectLocation)
	})

	return handler
}

func prepareHttpServer(port string, handler http.Handler) {
	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		MaxHeaderBytes:    maxHeaderBytes,
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(50 * time.Millisecond) // sleep for run http server
}

func TestAllowRedirects(t *testing.T) {
	port := "8988"
	handler := prepareHandler(port)
	prepareHttpServer(port, handler)

	cases := []struct {
		name     string
		endPoint string
		cfg      *Config
	}{
		{
			name:     "allowed moved permanently redirect 301",
			endPoint: "test_moved_permanently_redirect",
			cfg: &Config{
				Timeout:                   1,
				MaxIdleConnections:        1,
				MaxConnectionsPerHost:     1,
				MaxIdleConnectionsPerHost: 1,
				AllowFollowRedirect:       true,
			},
		},
		{
			name:     "allowed found redirect 302",
			endPoint: "test_found_redirect",
			cfg: &Config{
				Timeout:                   1,
				MaxIdleConnections:        1,
				MaxConnectionsPerHost:     1,
				MaxIdleConnectionsPerHost: 1,
				AllowFollowRedirect:       true,
			},
		},
		{
			name:     "allowed see other redirect 303",
			endPoint: "test_see_other_redirect",
			cfg: &Config{
				Timeout:                   1,
				MaxIdleConnections:        1,
				MaxConnectionsPerHost:     1,
				MaxIdleConnectionsPerHost: 1,
				AllowFollowRedirect:       true,
			},
		},
		{
			name:     "allowed temporary redirect 307",
			endPoint: "test_temporary_redirect",
			cfg: &Config{
				Timeout:                   1,
				MaxIdleConnections:        1,
				MaxConnectionsPerHost:     1,
				MaxIdleConnectionsPerHost: 1,
				AllowFollowRedirect:       true,
			},
		},
		{
			name:     "allowed permanent redirect 308",
			endPoint: "test_permanent_redirect",
			cfg: &Config{
				Timeout:                   1,
				MaxIdleConnections:        1,
				MaxConnectionsPerHost:     1,
				MaxIdleConnectionsPerHost: 1,
				AllowFollowRedirect:       true,
			},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			httpClient := getHTTPClient(testCase.cfg)

			resp, err := httpClient.Get(fmt.Sprintf("http://0.0.0.0:%s/%s", port, testCase.endPoint))
			if err == nil {
				if resp.StatusCode == http.StatusOK {
					bodyBytes, err := io.ReadAll(resp.Body)
					if err != nil {
						t.Fatal(err)
					}

					bodyString := string(bodyBytes)

					if bodyString != youHasBeenRedirectedResponse {
						t.Fatalf("error - body got: %s, want: %s", bodyString, youHasBeenRedirectedResponse)
					}
				}
			} else {
				t.Fatal("httpClient.Get error: ", err.Error())
			}

			if err := resp.Body.Close(); err != nil {
				t.Fatal("close response body error: ", err.Error())
			}
		})
	}
}

func TestDisallowRedirects(t *testing.T) {
	port := "8989"
	handler := prepareHandler(port)
	prepareHttpServer(port, handler)

	cases := []struct {
		name         string
		endPoint     string
		cfg          *Config
		expErrString string
	}{
		{
			name:     "not allowed moved permanently redirect 301",
			endPoint: "test_moved_permanently_redirect",
			cfg: &Config{
				Timeout:                   1,
				MaxIdleConnections:        1,
				MaxConnectionsPerHost:     1,
				MaxIdleConnectionsPerHost: 1,
				AllowFollowRedirect:       false,
			},
			expErrString: fmt.Sprintf("Get \"http://0.0.0.0:%s/handle_redirect\": redirect not allowed", port),
		},
		{
			name:     "not allowed found redirect 302",
			endPoint: "test_found_redirect",
			cfg: &Config{
				Timeout:                   1,
				MaxIdleConnections:        1,
				MaxConnectionsPerHost:     1,
				MaxIdleConnectionsPerHost: 1,
				AllowFollowRedirect:       false,
			},
			expErrString: fmt.Sprintf("Get \"http://0.0.0.0:%s/handle_redirect\": redirect not allowed", port),
		},
		{
			name:     "not allowed see other redirect 303",
			endPoint: "test_see_other_redirect",
			cfg: &Config{
				Timeout:                   1,
				MaxIdleConnections:        1,
				MaxConnectionsPerHost:     1,
				MaxIdleConnectionsPerHost: 1,
				AllowFollowRedirect:       false,
			},
			expErrString: fmt.Sprintf("Get \"http://0.0.0.0:%s/handle_redirect\": redirect not allowed", port),
		},
		{
			name:     "not allowed temporary redirect 307",
			endPoint: "test_temporary_redirect",
			cfg: &Config{
				Timeout:                   1,
				MaxIdleConnections:        1,
				MaxConnectionsPerHost:     1,
				MaxIdleConnectionsPerHost: 1,
				AllowFollowRedirect:       false,
			},
			expErrString: fmt.Sprintf("Get \"http://0.0.0.0:%s/handle_redirect\": redirect not allowed", port),
		},
		{
			name:     "not allowed permanent redirect 308",
			endPoint: "test_permanent_redirect",
			cfg: &Config{
				Timeout:                   1,
				MaxIdleConnections:        1,
				MaxConnectionsPerHost:     1,
				MaxIdleConnectionsPerHost: 1,
				AllowFollowRedirect:       false,
			},
			expErrString: fmt.Sprintf("Get \"http://0.0.0.0:%s/handle_redirect\": redirect not allowed", port),
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			httpClient := getHTTPClient(testCase.cfg)

			resp, err := httpClient.Get(fmt.Sprintf("http://0.0.0.0:%s/%s", port, testCase.endPoint))
			if err != nil {
				if err.Error() != testCase.expErrString {
					t.Fatalf("error - got: %s, want: %s", err.Error(), ErrRedirectNotAllowed)
				}
			} else {
				t.Fatalf("error must not be nil: %s, must be: %s", err, testCase.expErrString)
			}

			if err := resp.Body.Close(); err != nil {
				t.Fatal("close response body error: ", err.Error())
			}
		})
	}
}
