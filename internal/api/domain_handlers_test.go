package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockdb "github.com/danilshap/domains-generator/internal/db/mock"
	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHandleHome(t *testing.T) {
	server, err := NewServer(nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	server.handleHome(recorder, request)

	require.Equal(t, http.StatusSeeOther, recorder.Code)
	require.Equal(t, "/domains", recorder.Header().Get("Location"))
}

func TestHandleListDomains(t *testing.T) {
	testCases := []struct {
		name          string
		page          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			page: "1",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllDomains(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.GetAllDomainsRow{}, nil)

				store.EXPECT().
					GetDomainsCount(gomock.Any()).
					Times(1).
					Return(int64(10), nil)

				store.EXPECT().
					GetMailboxesStats(gomock.Any()).
					Times(1).
					Return(db.GetMailboxesStatsRow{
						ActiveCount:   5,
						InactiveCount: 3,
						TotalCount:    10,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidPage",
			page: "invalid",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllDomains(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.GetAllDomainsRow{}, nil)

				store.EXPECT().
					GetDomainsCount(gomock.Any()).
					Times(1).
					Return(int64(10), nil)

				store.EXPECT().
					GetMailboxesStats(gomock.Any()).
					Times(1).
					Return(db.GetMailboxesStatsRow{
						ActiveCount:   5,
						InactiveCount: 3,
						TotalCount:    10,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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

			server, err := NewServer(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := "/domains"
			if tc.page != "" {
				url += "?page=" + tc.page
			}
			request := httptest.NewRequest(http.MethodGet, url, nil)
			server.handleListDomains(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleNewDomainForm(t *testing.T) {
	server, err := NewServer(nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/domains/new", nil)
	server.handleNewDomainForm(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestHandleDeleteDomain(t *testing.T) {
	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupRequest: func() (*http.Request, error) {
				url := "/domains/1"
				request := httptest.NewRequest(http.MethodDelete, url, strings.NewReader(""))
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", "1")
				return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx)), nil
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDomain(gomock.Any(), gomock.Eq(int32(1))).
					Times(1).
					Return(nil)

				store.EXPECT().
					GetAllDomains(gomock.Any(), db.GetAllDomainsParams{
						Limit:  10,
						Offset: 0,
					}).
					Times(1).
					Return([]db.GetAllDomainsRow{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "<table")
			},
		},
		{
			name: "InvalidID",
			setupRequest: func() (*http.Request, error) {
				url := "/domains/invalid"
				request := httptest.NewRequest(http.MethodDelete, url, strings.NewReader(""))
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", "invalid")
				return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx)), nil
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No need to mock methods as error will occur during ID parsing
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "DeleteError",
			setupRequest: func() (*http.Request, error) {
				url := "/domains/1"
				request := httptest.NewRequest(http.MethodDelete, url, strings.NewReader(""))
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", "1")
				return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx)), nil
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteDomain(gomock.Any(), gomock.Eq(int32(1))).
					Times(1).
					Return(fmt.Errorf("delete error"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			server, err := NewServer(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			request, err := tc.setupRequest()
			require.NoError(t, err)

			server.handleDeleteDomain(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleDomainDetails(t *testing.T) {
	domain := db.GetDomainByIDRow{
		ID:       1,
		Name:     "test.com",
		Provider: "test",
		Status:   1,
	}

	testCases := []struct {
		name          string
		domainID      int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
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
					GetMailboxesCountByDomainID(gomock.Any(), gomock.Eq(domain.ID)).
					Times(1).
					Return(int64(0), nil)

				store.EXPECT().
					GetMailboxesByDomainID(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Mailbox{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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

			server, err := NewServer(store)
			require.NoError(t, err)
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

func TestHandleEditDomainForm(t *testing.T) {
	domain := db.GetDomainByIDRow{
		ID:       1,
		Name:     "test.com",
		Provider: "test",
		Status:   1,
	}

	testCases := []struct {
		name          string
		domainID      int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
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
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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

			server, err := NewServer(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/domains/%d/edit", tc.domainID)
			request := httptest.NewRequest(http.MethodGet, url, nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprint(tc.domainID))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))
			server.handleEditDomainForm(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleUpdateDomain(t *testing.T) {
	domain := db.Domain{
		ID:       1,
		Name:     "test.com",
		Provider: "test",
	}

	testCases := []struct {
		name          string
		domainID      int32
		formData      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			domainID: domain.ID,
			formData: "name=test.com&provider=test",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateDomainParams{
					ID:       domain.ID,
					Name:     domain.Name,
					Provider: domain.Provider,
				}
				store.EXPECT().
					UpdateDomain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)

				store.EXPECT().
					GetAllDomains(gomock.Any(), db.GetAllDomainsParams{
						Limit:  10,
						Offset: 0,
					}).
					Times(1).
					Return([]db.GetAllDomainsRow{}, nil)

				store.EXPECT().
					GetDomainsCount(gomock.Any()).
					Times(1).
					Return(int64(0), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "<table")
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

			server, err := NewServer(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/domains/%d", tc.domainID)
			request := httptest.NewRequest(http.MethodPut, url, strings.NewReader(tc.formData))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprint(tc.domainID))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))
			server.handleUpdateDomain(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleStatusForm(t *testing.T) {
	domain := db.GetDomainByIDRow{
		ID:     1,
		Status: 1,
	}

	testCases := []struct {
		name          string
		domainID      int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
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
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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

			server, err := NewServer(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/domains/%d/status", tc.domainID)
			request := httptest.NewRequest(http.MethodGet, url, nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprint(tc.domainID))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))
			server.handleStatusForm(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleCreateDomains(t *testing.T) {
	domain := db.Domain{
		ID:       1,
		Name:     "test.com",
		Provider: "test",
		Status:   1,
	}

	testCases := []struct {
		name          string
		formData      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			formData: "name=test.com&provider=test",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateDomainParams{
					Name:     "test.com",
					Provider: "test",
					Status:   1,
				}
				store.EXPECT().
					CreateDomain(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(domain, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusSeeOther, recorder.Code)
				require.Equal(t, fmt.Sprintf("/domains/%d", domain.ID), recorder.Header().Get("Location"))
			},
		},
		{
			name:     "InvalidDomain",
			formData: "name=invalid&provider=test",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDomain(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			server, err := NewServer(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodPost, "/domains", strings.NewReader(tc.formData))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			server.handleCreateDomain(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleUpdateStatus(t *testing.T) {
	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupRequest: func() (*http.Request, error) {
				url := "/domains/1/status"
				formData := "status=2"
				request := httptest.NewRequest(http.MethodPut, url, strings.NewReader(formData))
				request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", "1")
				return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx)), nil
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateDomainAndMailboxesStatus(gomock.Any(), gomock.Eq(int32(1)), gomock.Eq(int32(2))).
					Times(1).
					Return(nil)

				store.EXPECT().
					GetAllDomains(gomock.Any(), db.GetAllDomainsParams{
						Limit:  10,
						Offset: 0,
					}).
					Times(1).
					Return([]db.GetAllDomainsRow{}, nil)

				store.EXPECT().
					GetDomainsCount(gomock.Any()).
					Times(1).
					Return(int64(0), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Contains(t, recorder.Body.String(), "<table")
			},
		},
		{
			name: "InvalidID",
			setupRequest: func() (*http.Request, error) {
				url := "/domains/invalid/status"
				formData := "status=2"
				request := httptest.NewRequest(http.MethodPut, url, strings.NewReader(formData))
				request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", "invalid")
				return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx)), nil
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No need to mock methods as error will occur during ID parsing
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidStatus",
			setupRequest: func() (*http.Request, error) {
				url := "/domains/1/status"
				formData := "status=invalid"
				request := httptest.NewRequest(http.MethodPut, url, strings.NewReader(formData))
				request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", "1")
				return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx)), nil
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No need to mock methods as error will occur during status parsing
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			server, err := NewServer(store)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			request, err := tc.setupRequest()
			require.NoError(t, err)

			server.handleUpdateStatus(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
