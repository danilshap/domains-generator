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

func TestHandleListMailboxes(t *testing.T) {
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
					GetMailboxesWithFilters(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.GetMailboxesWithFiltersRow{}, nil)

				store.EXPECT().
					GetMailboxesCount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(int64(10), nil)
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
					GetMailboxesWithFilters(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.GetMailboxesWithFiltersRow{}, nil)

				store.EXPECT().
					GetMailboxesCount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(int64(10), nil)
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

			server, err := NewServer(store, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := "/mailboxes"
			if tc.page != "" {
				url += "?page=" + tc.page
			}
			request := httptest.NewRequest(http.MethodGet, url, nil)
			server.handleListMailboxes(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleMailboxDetails(t *testing.T) {
	mailbox := db.Mailbox{
		ID:       1,
		Address:  "test@test.com",
		DomainID: 1,
		Status:   1,
	}

	testCases := []struct {
		name          string
		mailboxID     int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			mailboxID: mailbox.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMailboxByID(gomock.Any(), gomock.Eq(mailbox.ID)).
					Times(1).
					Return(mailbox, nil)

				store.EXPECT().
					GetDomainByID(gomock.Any(), gomock.Eq(mailbox.DomainID)).
					Times(1).
					Return(db.GetDomainByIDRow{
						ID:       mailbox.DomainID,
						Name:     "test.com",
						Provider: "test",
						Status:   1,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			mailboxID: 2,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMailboxByID(gomock.Any(), gomock.Eq(int32(2))).
					Times(1).
					Return(db.Mailbox{}, fmt.Errorf("not found"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			mailboxID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				// No need to mock methods as error will occur during ID parsing
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

			server, err := NewServer(store, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/mailboxes/%d", tc.mailboxID)
			request := httptest.NewRequest(http.MethodGet, url, nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprint(tc.mailboxID))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))
			server.handleMailboxDetails(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleEditMailboxForm(t *testing.T) {
	mailbox := db.Mailbox{
		ID:       1,
		Address:  "test@test.com",
		DomainID: 1,
		Status:   1,
	}

	testCases := []struct {
		name          string
		mailboxID     int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			mailboxID: mailbox.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMailboxByID(gomock.Any(), gomock.Eq(mailbox.ID)).
					Times(1).
					Return(mailbox, nil)

				store.EXPECT().
					GetAllDomains(gomock.Any(), gomock.Any()).
					Times(1)
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

			server, err := NewServer(store, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/mailboxes/%d/edit", tc.mailboxID)
			request := httptest.NewRequest(http.MethodGet, url, nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprint(tc.mailboxID))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))
			server.handleEditMailboxForm(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleUpdateMailbox(t *testing.T) {
	mailbox := db.Mailbox{
		ID:       1,
		Address:  "test@test.com",
		DomainID: 1,
		Status:   1,
	}

	testCases := []struct {
		name          string
		mailboxID     int32
		formData      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			mailboxID: mailbox.ID,
			formData:  "address=test@test.com&domain_id=1",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateMailboxParams{
					ID:       mailbox.ID,
					Address:  mailbox.Address,
					DomainID: mailbox.DomainID,
				}
				store.EXPECT().
					UpdateMailbox(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusSeeOther, recorder.Code)
				require.Equal(t, fmt.Sprintf("/mailboxes/%d", mailbox.ID), recorder.Header().Get("Location"))
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

			server, err := NewServer(store, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/mailboxes/%d", tc.mailboxID)
			request := httptest.NewRequest(http.MethodPut, url, strings.NewReader(tc.formData))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprint(tc.mailboxID))
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))
			server.handleUpdateMailbox(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleDeleteMailbox(t *testing.T) {
	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupRequest: func() (*http.Request, error) {
				url := "/mailboxes/1"
				request := httptest.NewRequest(http.MethodDelete, url, strings.NewReader(""))
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", "1")
				return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx)), nil
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMailbox(gomock.Any(), gomock.Eq(int32(1))).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			setupRequest: func() (*http.Request, error) {
				url := "/mailboxes/invalid"
				request := httptest.NewRequest(http.MethodDelete, url, strings.NewReader(""))
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", "invalid")
				return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx)), nil
			},
			buildStubs: func(store *mockdb.MockStore) {
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

			server, err := NewServer(store, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			request, err := tc.setupRequest()
			require.NoError(t, err)

			server.handleDeleteMailbox(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
