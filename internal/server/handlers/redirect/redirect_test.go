package redirect_test

import (
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/commedesvlados/url-shortener/internal/lib/api"
	"github.com/commedesvlados/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/commedesvlados/url-shortener/internal/server/handlers/redirect"
	"github.com/commedesvlados/url-shortener/internal/server/handlers/redirect/mocks"
)

func TestRedirectHandler(t *testing.T) {
	testCases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com/",
		},
	}

	for _, tc := range testCases {
		//tc := tc

		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			// check url after redirect
			require.Equal(t, tc.url, redirectedURL)
		})
	}
}
