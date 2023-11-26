package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"

	"github.com/commedesvlados/url-shortener/internal/lib/api"
	"github.com/commedesvlados/url-shortener/internal/server/handlers/url/save"
)

const (
	scheme = "http"
	host   = "localhost:8080"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: scheme,
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: gofakeit.Word(),
		}).
		WithBasicAuth("user", "user").
		Expect().
		Status(200).JSON().Object().ContainsKey("alias")
}

func TestURLShortener_SaveRedirectDelete(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
		// TODO: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: scheme,
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// SAVE URL
			respSave := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("user", "user").
				Expect().Status(http.StatusOK).JSON().Object()

			if tc.error != "" {
				respSave.NotContainsKey("alias")
				respSave.Value("error").String().IsEqual(tc.error)
				return
			}

			alias := tc.alias
			if tc.alias != "" {
				respSave.Value("alias").String().IsEqual(tc.alias)
			} else {
				respSave.Value("alias").String().NotEmpty().Length().IsEqual(6)
				alias = respSave.Value("alias").String().Raw()
			}

			// REDIRECT TO URL

			testRedirect(t, alias, tc.url)

			// DELETE URL

			e.DELETE("/url/"+alias).
				WithBasicAuth("user", "user").
				Expect().Status(http.StatusOK)

			// TRY REDIRECT AGAIN

			testRedirectNotFound(t, alias)
		})
	}

}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   alias,
	}

	redirectedURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedURL)
}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   alias,
	}

	_, err := api.GetRedirect(u.String())
	require.ErrorIs(t, err, api.ErrInvalidStatusCode)
}
