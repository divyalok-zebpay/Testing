package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
)

const ExternalAPIMockRouter = "ExternalAPIMockRouter"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func New(ctx context.Context) (client HTTPClient) {
	if r := ctx.Value(ExternalAPIMockRouter); r != nil {
		if router, ok := r.(*gin.Engine); ok {
			client = &mockClient{
				router: router,
			}
		}
	}

	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 20,
		}
	}

	return client
}

type mockClient struct {
	router *gin.Engine
}

func (c *mockClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	c.router.ServeHTTP(w, req)
	return w.Result(), nil
}
