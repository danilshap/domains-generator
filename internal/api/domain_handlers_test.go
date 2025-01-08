package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/danilshap/domains-generator/internal/auth"
	mockdb "github.com/danilshap/domains-generator/internal/db/mock"
	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/danilshap/domains-generator/internal/middleware"
	"github.com/danilshap/domains-generator/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHandleListDomains(t *testing.T) {
	user := createRandomUser(t)
	n := 5
	domains := make([]db.GetAllDomainsRow, n)
	for i := 0; i < n; i++ {
		domains[i] = db.GetAllDomainsRow{
			ID:       int32(i + 1),
			Name:     utils.RandomString(10),
			Provider: "test",
			Status:   1,
			UserID:   user.ID,
		}
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetAllDomainsParams{
					Limit:  10,
					Offset: 0,
					UserID: user.ID,
				}
				store.EXPECT().
					GetAllDomains(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(domains, nil)

				store.EXPECT().
					GetDomainsCount(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(int64(n), nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchDomains(t, rec.Body.String(), domains)
			},
		},
		{
			name: "InternalError",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllDomains(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, fmt.Errorf("internal error"))
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/domains"
			request := httptest.NewRequest(http.MethodGet, url, nil)

			// Add authorized user to context
			authPayload := &auth.Payload{
				UserID: user.ID,
			}
			request = request.WithContext(context.WithValue(request.Context(), middleware.UserContextKey, authPayload))

			server.handleListDomains(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleUpdateDomainStatus(t *testing.T) {
	user := createRandomUser(t)
	domain := db.Domain{
		ID:     1,
		Status: 1,
		UserID: user.ID,
	}

	testCases := []struct {
		name          string
		domainID      int32
		formData      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			domainID: domain.ID,
			formData: "status=2",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateDomainAndMailboxesStatus(gomock.Any(), gomock.Eq(domain.ID), gomock.Eq(int32(2))).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusSeeOther, rec.Code)
				require.Equal(t, fmt.Sprintf("/domains/%d", domain.ID), rec.Header().Get("Location"))
			},
		},
		{
			name:     "InvalidStatus",
			domainID: domain.ID,
			formData: "status=invalid",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateDomainAndMailboxesStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/domains/%d/status", tc.domainID)
			request := httptest.NewRequest(http.MethodPut, url, strings.NewReader(tc.formData))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprint(tc.domainID))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))

			// Add authorized user to context
			authPayload := &auth.Payload{
				UserID: user.ID,
			}
			request = request.WithContext(context.WithValue(request.Context(), middleware.UserContextKey, authPayload))

			server.handleUpdateDomainStatus(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func createRandomUser(t *testing.T) db.User {
	return db.User{
		ID:       1,
		Username: utils.RandomString(6),
		Email:    utils.RandomEmail(),
	}
}

func requireBodyMatchDomains(t *testing.T, body string, domains []db.GetAllDomainsRow) {
	require.Contains(t, body, "Domains")
	for _, domain := range domains {
		require.Contains(t, body, domain.Name)
	}
}

func TestIsValidDomain(t *testing.T) {
	testCases := []struct {
		name     string
		domain   string
		expected bool
	}{
		{
			name:     "ValidDomain",
			domain:   "example.com",
			expected: true,
		},
		{
			name:     "ValidDomainWithSubdomain",
			domain:   "sub.example.com",
			expected: false,
		},
		{
			name:     "InvalidDomainStartsWithHyphen",
			domain:   "-example.com",
			expected: false,
		},
		{
			name:     "InvalidDomainEndsWithHyphen",
			domain:   "example-.com",
			expected: false,
		},
		{
			name:     "InvalidDomainTooShort",
			domain:   "e.c",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidDomain(tc.domain)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestHandleDeleteDomain(t *testing.T) {
	user := createRandomUser(t)
	domain := db.Domain{
		ID:     1,
		UserID: user.ID,
	}

	testCases := []struct {
		name          string
		domainID      int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			domainID: domain.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDomain(gomock.Any(), gomock.Eq(domain.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusSeeOther, rec.Code)
				require.Equal(t, "/domains", rec.Header().Get("Location"))
			},
		},
		{
			name:     "InvalidID",
			domainID: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDomain(gomock.Any(), gomock.Any()).
					Times(1).
					Return(fmt.Errorf("invalid id"))
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/domains/%d", tc.domainID)
			request := httptest.NewRequest(http.MethodDelete, url, nil)

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprint(tc.domainID))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))

			server.handleDeleteDomain(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleDomainDetails(t *testing.T) {
	domain := db.GetDomainByIDRow{
		ID:     1,
		Name:   "test.com",
		Status: 1,
	}
	mailboxes := []db.GetMailboxesByDomainIDRow{
		{
			ID:       1,
			Address:  "test1@test.com",
			DomainID: domain.ID,
		},
	}

	testCases := []struct {
		name          string
		domainID      int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			domainID: domain.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDomainByID(gomock.Any(), gomock.Eq(domain.ID)).
					Times(1).
					Return(domain, nil)

				store.EXPECT().
					GetMailboxesByDomainID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(mailboxes, nil)

				store.EXPECT().
					GetMailboxesCountByDomainID(gomock.Any(), gomock.Eq(domain.ID)).
					Times(1).
					Return(int64(len(mailboxes)), nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				require.Contains(t, rec.Body.String(), domain.Name)
			},
		},
		{
			name:     "DomainNotFound",
			domainID: 999,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDomainByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetDomainByIDRow{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/domains/%d", tc.domainID)
			request := httptest.NewRequest(http.MethodGet, url, nil)

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprint(tc.domainID))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))

			server.handleDomainDetails(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleBulkMailboxesForm(t *testing.T) {
	domain := db.GetDomainByIDRow{
		ID:     1,
		Name:   "test.com",
		Status: 1,
	}

	testCases := []struct {
		name          string
		domainID      int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			domainID: domain.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDomainByID(gomock.Any(), gomock.Eq(domain.ID)).
					Times(1).
					Return(domain, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				require.Contains(t, rec.Body.String(), "Create Multiple Mailboxes")
			},
		},
		{
			name:     "InvalidDomainID",
			domainID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDomainByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetDomainByIDRow{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/domains/%d/bulk-mailboxes", tc.domainID)
			request := httptest.NewRequest(http.MethodGet, url, nil)

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprint(tc.domainID))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))

			server.handleBulkMailboxesForm(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleCreateDomain(t *testing.T) {
	user := createRandomUser(t)
	domain := db.Domain{
		ID:       1,
		Name:     "example.com",
		Provider: "test",
		Status:   1,
		UserID:   user.ID,
	}

	testCases := []struct {
		name          string
		formData      map[string]string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			formData: map[string]string{
				"name":     domain.Name,
				"provider": domain.Provider,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateDomainParams{
					Name:     domain.Name,
					Provider: domain.Provider,
					Status:   1,
					UserID:   user.ID,
				}

				store.EXPECT().
					CreateDomain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(domain, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusSeeOther, rec.Code)
				require.Equal(t, fmt.Sprintf("/domains/%d", domain.ID), rec.Header().Get("Location"))
			},
		},
		{
			name: "InvalidDomainName",
			formData: map[string]string{
				"name":     "invalid..com",
				"provider": domain.Provider,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDomain(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "DBError",
			formData: map[string]string{
				"name":     domain.Name,
				"provider": domain.Provider,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDomain(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Domain{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			form := url.Values{}
			for k, v := range tc.formData {
				form.Add(k, v)
			}

			request := httptest.NewRequest(http.MethodPost, "/domains", strings.NewReader(form.Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Add authorized user to context
			authPayload := &auth.Payload{
				UserID: user.ID,
			}
			request = request.WithContext(context.WithValue(request.Context(), middleware.UserContextKey, authPayload))

			server.handleCreateDomain(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
