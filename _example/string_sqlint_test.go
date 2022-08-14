package example

import (
	"database/sql"
	driver "database/sql/driver"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mediaTypeSelect = `SELECT format FROM media WHERE id = ?`
	mediaTypeUpdate = `Update media SET format = ? WHERE id = ?`
)

func TestExampleSQLMediaTypeIntOnly(t *testing.T) {

	tests := map[string]struct {
		setupMock func(t testing.TB, mocks *MockSQL)
		tester    func(t testing.TB, db *sql.DB)
	}{
		"NonNullable select nil": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{nil}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.Equal(t, MediaTypeInt(""), status)
			},
		},
		"Nullable Select Null": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{nil}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getNullMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.False(t, status.Valid)
			},
		},
		"Select a string": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{MediaTypeIntFlac.String()}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.Equal(t, MediaTypeIntFlac, status)
			},
		},
		"Nullable select an int": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{int(2)}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getNullMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeIntMp3, status.MediaTypeInt)
			},
		},
		"Nullable select an int64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{int64(2)}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getNullMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeIntMp3, status.MediaTypeInt)
			},
		},
		"Nullable select an float64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{float64(2)}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getNullMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeIntMp3, status.MediaTypeInt)
			},
		},
		"Nullable select an uint": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{uint(2)}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getNullMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeIntMp3, status.MediaTypeInt)
			},
		},
		"Nullable select an uint64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{uint64(2)}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getNullMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeIntMp3, status.MediaTypeInt)
			},
		},
		"Nullable select an *uint64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := uint64(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{&val}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getNullMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeIntMp3, status.MediaTypeInt)
			},
		},
		"Nullable select an *int64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := int64(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{&val}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getNullMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeIntMp3, status.MediaTypeInt)
			},
		},
		"Nullable select an *float64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := float64(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{&val}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getNullMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeIntMp3, status.MediaTypeInt)
			},
		},
		"Nullable select an *uint": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := uint(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(mediaTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{&val}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getNullMediaTypeInt(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeIntMp3, status.MediaTypeInt)
			},
		},
		"standard update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare(mediaTypeUpdate).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(3, hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setMediaTypeInt(conn, MediaTypeIntOgg)
				require.NoError(t, err, "failed updating project status")
			},
		},
		"nullable update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare(mediaTypeUpdate).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(3, hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setMediaTypeInt(conn, NewNullMediaTypeInt(MediaTypeIntOgg))
				require.NoError(t, err, "failed updating project status")
			},
		},
		"nullable invalid update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare(mediaTypeUpdate).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(nil, hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setMediaTypeInt(conn, NullMediaTypeInt{})
				require.NoError(t, err, "failed updating project status")
			},
		},
	}

	driverctrl := gomock.NewController(t)
	driver := NewMockDriver(driverctrl)
	defer func() {
		driverctrl.Finish()
	}()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			require.NotNil(t, tc.setupMock)
			require.NotNil(t, tc.tester)

			sql.Register(t.Name(), driver)
			mocks, finish := WithMockSQL(t)
			defer finish()

			driver.EXPECT().Open(dataSourceName).Return(mocks.Conn, nil)

			conn, err := sql.Open(t.Name(), dataSourceName)
			require.NoError(t, err, "failed opening mock db")

			tc.setupMock(t, mocks)

			tc.tester(t, conn)

		})
	}
}

func getMediaTypeInt(db *sql.DB) (state MediaTypeInt, err error) {
	err = db.QueryRow(mediaTypeSelect, hardcodedProjectID).Scan(&state)
	return
}

func getNullMediaTypeInt(db *sql.DB) (state NullMediaTypeInt, err error) {
	err = db.QueryRow(mediaTypeSelect, hardcodedProjectID).Scan(&state)
	return
}

func setMediaTypeInt(db *sql.DB, state interface{}) error {
	_, err := db.Exec(mediaTypeUpdate, state, hardcodedProjectID)
	return err
}

func TestSQLMediaTypeIntExtras(t *testing.T) {

	assert.Equal(t, "22", MediaTypeInt("22").String(), "String value is not correct")

	_, err := ParseMediaTypeInt(`NotAStatus`)
	assert.Error(t, err, "Should have had an error parsing a non status")

	var (
		intVal           int          = 3
		strVal           string       = "flac"
		enumVal          MediaTypeInt = MediaTypeIntFlac
		nullInt          *int
		nullInt64        *int64
		nullUint         *uint
		nullUint64       *uint64
		nullString       *string
		nullMediaTypeInt *MediaTypeInt
	)

	tests := map[string]struct {
		input  interface{}
		result NullMediaTypeInt
	}{
		"nil": {},
		"val": {
			input: MediaTypeIntOgg,
			result: NullMediaTypeInt{
				MediaTypeInt: MediaTypeIntOgg,
				Valid:        true,
			},
		},
		"ptr": {
			input: &enumVal,
			result: NullMediaTypeInt{
				MediaTypeInt: MediaTypeIntFlac,
				Valid:        true,
			},
		},
		"string": {
			input: strVal,
			result: NullMediaTypeInt{
				MediaTypeInt: MediaTypeIntFlac,
				Valid:        true,
			},
		},
		"*string": {
			input: &strVal,
			result: NullMediaTypeInt{
				MediaTypeInt: MediaTypeIntFlac,
				Valid:        true,
			},
		},
		"invalid string": {
			input: "random value",
		},
		"[]byte": {
			input: []byte(MediaTypeIntMp4.String()),
			result: NullMediaTypeInt{
				MediaTypeInt: MediaTypeIntMp4,
				Valid:        true,
			},
		},
		"int": {
			input: intVal,
			result: NullMediaTypeInt{
				MediaTypeInt: MediaTypeIntOgg,
				Valid:        true,
			},
		},
		"*int": {
			input: &intVal,
			result: NullMediaTypeInt{
				MediaTypeInt: MediaTypeIntOgg,
				Valid:        true,
			},
		},
		"nullInt": {
			input: nullInt,
		},
		"nullInt64": {
			input: nullInt64,
		},
		"nullUint": {
			input: nullUint,
		},
		"nullUint64": {
			input: nullUint64,
		},
		"nullString": {
			input: nullString,
		},
		"nullMediaTypeInt": {
			input: nullMediaTypeInt,
		},
		"int as []byte": {
			input: []byte("3"),
			result: NullMediaTypeInt{
				MediaTypeInt: MediaTypeIntOgg,
				Valid:        true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			status := NewNullMediaTypeInt(tc.input)
			assert.Equal(t, status, tc.result)

		})
	}

}
