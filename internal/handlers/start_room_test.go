package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/shuto.sawaki/elmo-project/internal/ai"
	"github.com/shuto.sawaki/elmo-project/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// テストファイル内では、テスト関数だけを記述します。
// RoomHandlerなどの構造体定義は room.go にあるので、ここでは書きません。

func TestStartRoomHandler_Integration(t *testing.T) {
	err := godotenv.Load("../../.env")
	require.NoError(t, err, ".envファイルの読み込みに失敗しました")

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	realAIGenerator, err := ai.NewGeminiAIGenerator(ctx)
	require.NoError(t, err, "本物のAIジェネレータの初期化に失敗しました")

	roomID := "r001"
	rows := sqlmock.NewRows([]string{"id", "title", "description", "status"}).
		AddRow(roomID, "Go言語のテスト", "テストコードの書き方について議論する部屋", "not started")
	mock.ExpectQuery(`SELECT id, title, description, status FROM rooms WHERE id = \$1`).WithArgs(roomID).WillReturnRows(rows)

	mock.ExpectExec(`UPDATE rooms SET status = \$1, initial_question = \$2 WHERE id = \$3`).
		WithArgs("inprogress", sqlmock.AnyArg(), roomID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	participantRows := sqlmock.NewRows([]string{"id", "user_name"}) // ★ 修正：DBのカラム名に合わせる
	mock.ExpectQuery(`SELECT u.id, u.user_name FROM participants p JOIN users u ON p.user_id = u.id WHERE p.room_id = \$1`).
		WithArgs(roomID).
		WillReturnRows(participantRows)

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID+"/start", nil)
	w := httptest.NewRecorder()

	handler := NewRoomHandler(db, realAIGenerator)
	handler.StartRoom(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode, "期待したステータスコードは200 OKでしたが、実際は%dでした", res.StatusCode)

	var response models.StartRoomResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.InitialQuestion)
	assert.Equal(t, roomID, response.RoomInfo.RoomID)
	assert.Equal(t, "inprogress", response.RoomInfo.Status)
	assert.Len(t, response.Participants, 0)

	assert.NoError(t, mock.ExpectationsWereMet())
}