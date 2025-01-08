package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	mockdb "github.com/danilshap/domains-generator/internal/db/mock"
	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/danilshap/domains-generator/pkg/config"
	"github.com/danilshap/domains-generator/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	cfg := &config.Config{
		TokenSynnetricKey: utils.RandomString(32),
	}

	server, err := NewServer(store, cfg)
	require.NoError(t, err)

	return server
}

func TestHandleRegister(t *testing.T) {
	// Generate a random user
	user := db.User{
		ID:             1,
		Username:       utils.RandomString(6), // handleRegister requires a minimum of 4 characters
		HashedPassword: utils.RandomString(12),
		Email:          utils.RandomEmail(),
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
				"username": user.Username,
				"email":    user.Email,
				"password": "secret",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Expect exactly 1 call to CreateUser which returns the user
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				// Expect a redirect to "/"
				require.Equal(t, http.StatusSeeOther, rec.Code)
				require.Equal(t, "/", rec.Header().Get("Location"))
			},
		},
		{
			name: "ShortUsername",
			formData: map[string]string{
				"username": "ab", // less than 4 characters
				"email":    user.Email,
				"password": "secret",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// In this scenario we expect no calls to CreateUser,
				// because the handler will return 400 before calling the store.
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "CreateUserError",
			formData: map[string]string{
				"username": user.Username,
				"email":    user.Email,
				"password": "secret",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Simulate an error while creating the user
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, fmt.Errorf("some db error"))
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

			mockStore := mockdb.NewMockStore(ctrl)
			tc.buildStubs(mockStore)

			server := newTestServer(t, mockStore)
			recorder := httptest.NewRecorder()

			form := url.Values{}
			for k, v := range tc.formData {
				form.Add(k, v)
			}

			request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			server.handleRegister(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleLogin(t *testing.T) {
	// Create a test user with the hashed password "secret"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := db.User{
		ID:             1,
		Username:       utils.RandomString(6),
		HashedPassword: string(hashedPassword),
		Email:          utils.RandomEmail(),
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
				"email":    user.Email,
				"password": "secret",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusSeeOther, rec.Code)
				require.Equal(t, "/", rec.Header().Get("Location"))
			},
		},
		{
			name: "UserNotFound",
			formData: map[string]string{
				"email":    "doesnot@exist.com",
				"password": "secret",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "WrongPassword",
			formData: map[string]string{
				"email":    user.Email,
				"password": "wrong-password",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Return the real user, but the password won't match
				store.EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name: "DBError",
			formData: map[string]string{
				"email":    user.Email,
				"password": "secret",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Simulate a database error
				store.EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(db.User{}, fmt.Errorf("db error"))
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			tc.buildStubs(mockStore)

			server := newTestServer(t, mockStore)
			recorder := httptest.NewRecorder()

			form := url.Values{}
			for k, v := range tc.formData {
				form.Add(k, v)
			}

			request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			server.handleLogin(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleLoginPage(t *testing.T) {
	testCases := []struct {
		name          string
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				require.Contains(t, rec.Body.String(), "Sign in")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodGet, "/login", nil)
			server.handleLoginPage(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleRegisterPage(t *testing.T) {
	testCases := []struct {
		name          string
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				require.Contains(t, rec.Body.String(), "Create new account")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodGet, "/register", nil)
			server.handleRegisterPage(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHandleLogout(t *testing.T) {
	testCases := []struct {
		name          string
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusSeeOther, rec.Code)
				require.Equal(t, "/login", rec.Header().Get("Location"))

				cookies := rec.Result().Cookies()
				require.Len(t, cookies, 1)
				cookie := cookies[0]
				require.Equal(t, "token", cookie.Name)
				require.Empty(t, cookie.Value)
				require.Equal(t, -1, cookie.MaxAge)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodPost, "/logout", nil)
			server.handleLogout(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
