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
	imageTypeSelect = `SELECT format FROM images WHERE id = ?`
	imageTypeUpdate = `Update images SET format = ? WHERE id = ?`
)

func TestExampleSQLIntOnly(t *testing.T) {

	tests := map[string]struct {
		setupMock func(t testing.TB, mocks *MockSQL)
		tester    func(t testing.TB, db *sql.DB)
	}{
		"NonNullable select nil": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(imageTypeSelect).Return(mocks.Stmt, nil),
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
				status, err := getImageType(conn)
				require.NoError(t, err, "failed getting project status")
				require.Equal(t, ImageTypeJpeg, status)
			},
		},
		"Nullable Select Null": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(imageTypeSelect).Return(mocks.Stmt, nil),
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
				status, err := getNullImageType(conn)
				require.NoError(t, err, "failed getting project status")
				require.False(t, status.Valid)
			},
		},
		"Select a string": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(imageTypeSelect).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{ImageTypePng.String()}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getImageType(conn)
				require.NoError(t, err, "failed getting project status")
				require.Equal(t, ImageTypePng, status)
			},
		},
		"Nullable select an int": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(imageTypeSelect).Return(mocks.Stmt, nil),
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
				status, err := getNullImageType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ImageTypePng, status.ImageType)
			},
		},
		"Nullable select an int64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(imageTypeSelect).Return(mocks.Stmt, nil),
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
				status, err := getNullImageType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ImageTypePng, status.ImageType)
			},
		},
		"Nullable select an uint": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(imageTypeSelect).Return(mocks.Stmt, nil),
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
				status, err := getNullImageType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ImageTypePng, status.ImageType)
			},
		},
		"Nullable select an uint64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(imageTypeSelect).Return(mocks.Stmt, nil),
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
				status, err := getNullImageType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ImageTypePng, status.ImageType)
			},
		},
		"Nullable select an *uint64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := uint64(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(imageTypeSelect).Return(mocks.Stmt, nil),
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
				status, err := getNullImageType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ImageTypePng, status.ImageType)
			},
		},
		"Nullable select an *int64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := int64(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(imageTypeSelect).Return(mocks.Stmt, nil),
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
				status, err := getNullImageType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ImageTypePng, status.ImageType)
			},
		},
		"Nullable select an *uint": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := uint(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare(imageTypeSelect).Return(mocks.Stmt, nil),
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
				status, err := getNullImageType(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ImageTypePng, status.ImageType)
			},
		},
		"standard update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare(imageTypeUpdate).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(int64(ImageTypeTiff), hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setImageType(conn, ImageTypeTiff)
				require.NoError(t, err, "failed updating project status")
			},
		},
		"nullable update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare(imageTypeUpdate).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(int64(ImageTypeTiff), hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setImageType(conn, NewNullImageType(ImageTypeTiff))
				require.NoError(t, err, "failed updating project status")
			},
		},
		"nullable invalid update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare(imageTypeUpdate).Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(nil, hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setImageType(conn, NullImageType{})
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

func getImageType(db *sql.DB) (state ImageType, err error) {
	err = db.QueryRow(imageTypeSelect, hardcodedProjectID).Scan(&state)
	return
}

func getNullImageType(db *sql.DB) (state NullImageType, err error) {
	err = db.QueryRow(imageTypeSelect, hardcodedProjectID).Scan(&state)
	return
}

func setImageType(db *sql.DB, state interface{}) error {
	_, err := db.Exec(imageTypeUpdate, state, hardcodedProjectID)
	return err
}

func TestSQLIntExtras(t *testing.T) {

	assert.Equal(t, "ImageType(22)", ImageType(22).String(), "String value is not correct")

	_, err := ParseImageType(`NotAStatus`)
	assert.Error(t, err, "Should have had an error parsing a non status")

	var (
		intVal        int       = 3
		strVal        string    = "png"
		enumVal       ImageType = ImageTypeGif
		nullInt       *int
		nullInt64     *int64
		nullUint      *uint
		nullUint64    *uint64
		nullString    *string
		nullImageType *ImageType
	)

	tests := map[string]struct {
		input  interface{}
		result NullImageType
	}{
		"nil": {},
		"val": {
			input: ImageTypeTiff,
			result: NullImageType{
				ImageType: ImageTypeTiff,
				Valid:     true,
			},
		},
		"ptr": {
			input: &enumVal,
			result: NullImageType{
				ImageType: ImageTypeGif,
				Valid:     true,
			},
		},
		"string": {
			input: strVal,
			result: NullImageType{
				ImageType: ImageTypePng,
				Valid:     true,
			},
		},
		"*string": {
			input: &strVal,
			result: NullImageType{
				ImageType: ImageTypePng,
				Valid:     true,
			},
		},
		"invalid string": {
			input: "random value",
		},
		"[]byte": {
			input: []byte(ImageTypeGif.String()),
			result: NullImageType{
				ImageType: ImageTypeGif,
				Valid:     true,
			},
		},
		"int": {
			input: intVal,
			result: NullImageType{
				ImageType: ImageTypeTiff,
				Valid:     true,
			},
		},
		"*int": {
			input: &intVal,
			result: NullImageType{
				ImageType: ImageTypeTiff,
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
		"nullImageType": {
			input: nullImageType,
		},
		"int as []byte": {
			input: []byte("3"),
			result: NullImageType{
				ImageType: ImageTypeTiff,
				Valid:     true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			status := NewNullImageType(tc.input)
			assert.Equal(t, status, tc.result)

		})
	}

}
