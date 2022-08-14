package example

import (
	"database/sql"
	driver "database/sql/driver"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExampleSQLMediaTypeOnly(t *testing.T) {

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
				status, err := getMediaType(conn)
				require.NoError(t, err, "failed getting project status")
				require.Equal(t, MediaType(""), status)
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
				status, err := getNullMediaType(conn)
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
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{MediaTypeFlac.String()}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getMediaType(conn)
				require.NoError(t, err, "failed getting project status")
				require.Equal(t, MediaTypeFlac, status)
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
				status, err := getNullMediaType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeMp3, status.MediaType)
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
				status, err := getNullMediaType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeMp3, status.MediaType)
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
				status, err := getNullMediaType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeMp3, status.MediaType)
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
				status, err := getNullMediaType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeMp3, status.MediaType)
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
				status, err := getNullMediaType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeMp3, status.MediaType)
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
				status, err := getNullMediaType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeMp3, status.MediaType)
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
				status, err := getNullMediaType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeMp3, status.MediaType)
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
				status, err := getNullMediaType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeMp3, status.MediaType)
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
				status, err := getNullMediaType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, MediaTypeMp3, status.MediaType)
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
				err := setMediaType(conn, MediaTypeOgg)
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
				err := setMediaType(conn, NewNullMediaType(MediaTypeOgg))
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
				err := setMediaType(conn, NullMediaType{})
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

func getMediaType(db *sql.DB) (state MediaType, err error) {
	err = db.QueryRow(mediaTypeSelect, hardcodedProjectID).Scan(&state)
	return
}

func getNullMediaType(db *sql.DB) (state NullMediaType, err error) {
	err = db.QueryRow(mediaTypeSelect, hardcodedProjectID).Scan(&state)
	return
}

func setMediaType(db *sql.DB, state interface{}) error {
	_, err := db.Exec(mediaTypeUpdate, state, hardcodedProjectID)
	return err
}

func TestSQLMediaTypeExtras(t *testing.T) {

	assert.Equal(t, "22", MediaType("22").String(), "String value is not correct")

	_, err := ParseMediaType(`NotAStatus`)
	assert.Error(t, err, "Should have had an error parsing a non status")

	var (
		intVal        int       = 3
		strVal        string    = "flac"
		enumVal       MediaType = MediaTypeFlac
		nullInt       *int
		nullInt64     *int64
		nullUint      *uint
		nullUint64    *uint64
		nullString    *string
		nullMediaType *MediaType
	)

	tests := map[string]struct {
		input  interface{}
		result NullMediaType
	}{
		"nil": {},
		"val": {
			input: MediaTypeOgg,
			result: NullMediaType{
				MediaType: MediaTypeOgg,
				Valid:     true,
			},
		},
		"ptr": {
			input: &enumVal,
			result: NullMediaType{
				MediaType: MediaTypeFlac,
				Valid:     true,
			},
		},
		"string": {
			input: strVal,
			result: NullMediaType{
				MediaType: MediaTypeFlac,
				Valid:     true,
			},
		},
		"*string": {
			input: &strVal,
			result: NullMediaType{
				MediaType: MediaTypeFlac,
				Valid:     true,
			},
		},
		"invalid string": {
			input: "random value",
		},
		"[]byte": {
			input: []byte(MediaTypeMp4.String()),
			result: NullMediaType{
				MediaType: MediaTypeMp4,
				Valid:     true,
			},
		},
		"int": {
			input: intVal,
			result: NullMediaType{
				MediaType: MediaTypeOgg,
				Valid:     true,
			},
		},
		"*int": {
			input: &intVal,
			result: NullMediaType{
				MediaType: MediaTypeOgg,
				Valid:     true,
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
		"nullMediaType": {
			input: nullMediaType,
		},
		"int as []byte": {
			input: []byte("3"),
			result: NullMediaType{
				MediaType: MediaTypeOgg,
				Valid:     true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			status := NewNullMediaType(tc.input)
			assert.Equal(t, status, tc.result)

		})
	}

}
