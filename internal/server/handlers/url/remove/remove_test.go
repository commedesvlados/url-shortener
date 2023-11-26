package remove_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/commedesvlados/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/commedesvlados/url-shortener/internal/server/handlers/url/remove"
	"github.com/commedesvlados/url-shortener/internal/server/handlers/url/remove/mocks"
)

func TestDeleteHandler(t *testing.T) {
	testCases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
		{
			name:      "Empty alias",
			alias:     "",
			respError: "field ALIAS is a required field",
		},
		{
			name:      "DeleteURL Error",
			alias:     "test_alias",
			respError: "failed to delete url by alias",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).
					Once()
			}

			// chi chi chi
			router := chi.NewRouter()
			router.Delete("/url/{alias}", remove.New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			req, err := http.NewRequest(http.MethodDelete, "/url/"+tc.alias, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if tc.alias == "" {
				require.NotEqual(t, rr.Code, http.StatusOK)
			} else {
				require.Equal(t, rr.Code, http.StatusOK)
			}
		})
	}
}
