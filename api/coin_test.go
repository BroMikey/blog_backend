package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/BroMikey/blog_backend/db/mock"
	db "github.com/BroMikey/blog_backend/db/sqlc"
	"github.com/BroMikey/blog_backend/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetCoin(t *testing.T) {
	user := db.Users{}

	coin := randomCoin(user)

	testCases := []struct {
		name          string
		coinID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			coinID: coin.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetCoin(gomock.Any(), gomock.Eq(coin.ID)).
					Times(1).
					Return(coin, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMathCoin(t, recorder.Body, coin)
			},
		},
		{
			name:   "NotFound",
			coinID: coin.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetCoin(gomock.Any(), gomock.Eq(coin.ID)).
					Times(1).
					Return(db.Coin{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			coinID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetCoin(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			coinID: coin.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetCoin(gomock.Any(), gomock.Eq(coin.ID)).
					Times(1).
					Return(db.Coin{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/coin/%d", tc.coinID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)

		})
	}

}

func randomCoin(user db.Users) db.Coin {
	return db.Coin{
		ID:       utils.RandomInt(1, 1000),
		Uid:      user.Uid,
		Balance:  utils.RandomAmount(1, 1000),
		CoinType: utils.RandomCoinType(),
	}
}

func requireBodyMathCoin(t *testing.T, body *bytes.Buffer, coin db.Coin) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCoin db.Coin

	err = json.Unmarshal(data, &gotCoin)
	require.NoError(t, err)
	require.Equal(t, coin, gotCoin)
}
