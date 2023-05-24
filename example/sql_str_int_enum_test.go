//go:build example
// +build example

package example

import (
	"database/sql"
	driver "database/sql/driver"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExampleSQLStrIntOnly(t *testing.T) {
	tests := map[string]struct {
		setupMock func(t testing.TB, mocks *MockSQL)
		tester    func(t testing.TB, db *sql.DB)
	}{
		"NonNullable select nil": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT god FROM job WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`god`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{nil}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				god, err := getGreekGod(conn)
				require.NoError(t, err, "failed getting project god")
				require.Equal(t, GreekGod(""), god)
			},
		},
		"Nullable Select Null": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT god FROM job WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`god`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{nil}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				god, err := getNullGreekGod(conn)
				require.NoError(t, err, "failed getting project god")
				require.False(t, god.Valid)
			},
		},
		"Select a string": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT god FROM job WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`god`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{GreekGodAthena.String()}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				god, err := getGreekGod(conn)
				require.NoError(t, err, "failed getting project god")
				require.Equal(t, GreekGodAthena, god)
			},
		},
		"Nullable select an int": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT god FROM job WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`god`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{int(2)}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				god, err := getNullGreekGod(conn)
				require.NoError(t, err, "failed getting project god")
				require.True(t, god.Valid)
				require.Equal(t, GreekGodApollo, god.GreekGod)
			},
		},
		"Nullable select an int64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT god FROM job WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`god`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{int64(2)}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				god, err := getNullGreekGod(conn)
				require.NoError(t, err, "failed getting project god")
				require.True(t, god.Valid)
				require.Equal(t, GreekGodApollo, god.GreekGod)
			},
		},
		"Nullable select an uint": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT god FROM job WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`god`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{uint(2)}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				god, err := getNullGreekGod(conn)
				require.NoError(t, err, "failed getting project god")
				require.True(t, god.Valid)
				require.Equal(t, GreekGodApollo, god.GreekGod)
			},
		},
		"Nullable select an uint64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT god FROM job WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`god`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{uint64(2)}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				god, err := getNullGreekGod(conn)
				require.NoError(t, err, "failed getting project god")
				require.True(t, god.Valid)
				require.Equal(t, GreekGodApollo, god.GreekGod)
			},
		},
		"Nullable select an *uint64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := uint64(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT god FROM job WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`god`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{&val}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				god, err := getNullGreekGod(conn)
				require.NoError(t, err, "failed getting project god")
				require.True(t, god.Valid)
				require.Equal(t, GreekGodApollo, god.GreekGod)
			},
		},
		"Nullable select an *int64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := int64(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT god FROM job WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`god`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{&val}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				god, err := getNullGreekGod(conn)
				require.NoError(t, err, "failed getting project god")
				require.True(t, god.Valid)
				require.Equal(t, GreekGodApollo, god.GreekGod)
			},
		},
		"Nullable select an *uint": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := uint(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT god FROM job WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`god`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{&val}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				god, err := getNullGreekGod(conn)
				require.NoError(t, err, "failed getting project god")
				require.True(t, god.Valid)
				require.Equal(t, GreekGodApollo, god.GreekGod)
			},
		},
		"standard update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare("Update job SET god = ? WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(sqlIntGreekGodValue[GreekGodAres], hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				err := setGreekGod(conn, GreekGodAres)
				require.NoError(t, err, "failed updating project god")
			},
		},
		"nullable string update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare("Update job SET god = ? WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(GreekGodAres.String(), hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				err := setGreekGod(conn, NewNullGreekGod(GreekGodAres))
				require.NoError(t, err, "failed updating project god")
			},
		},
		"nullable invalid string update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare("Update job SET god = ? WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(nil, hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork god
				err := setGreekGod(conn, NullGreekGod{})
				require.NoError(t, err, "failed updating project god")
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

func getGreekGod(db *sql.DB) (state GreekGod, err error) {
	err = db.QueryRow(`SELECT god FROM job WHERE id = ?`, hardcodedProjectID).Scan(&state)
	return
}

func getNullGreekGod(db *sql.DB) (state NullGreekGod, err error) {
	err = db.QueryRow(`SELECT god FROM job WHERE id = ?`, hardcodedProjectID).Scan(&state)
	return
}

func setGreekGod(db *sql.DB, state interface{}) error {
	_, err := db.Exec(`Update job SET god = ? WHERE id = ?`, state, hardcodedProjectID)
	return err
}

func TestSQLStrIntExtras(t *testing.T) {
	assert.Equal(t, sqlIntGreekGodMap[20], GreekGod("athena"), "String value is not correct")

	_, err := ParseGreekGod(`NotAStatus`)
	assert.Error(t, err, "Should have had an error parsing a non god")
	assert.True(t, errors.Is(err, ErrInvalidGreekGod))

	var (
		intVal       int      = 3
		strVal       string   = "completed"
		enumVal      JobState = JobStateCompleted
		nullInt      *int
		nullInt64    *int64
		nullUint     *uint
		nullUint64   *uint64
		nullString   *string
		nullJobState *JobState
	)

	tests := map[string]struct {
		input  interface{}
		result NullJobState
	}{
		"nil": {},
		"val": {
			input: JobStateFailed,
			result: NullJobState{
				JobState: JobStateFailed,
				Valid:    true,
			},
		},
		"ptr": {
			input: &enumVal,
			result: NullJobState{
				JobState: JobStateCompleted,
				Valid:    true,
			},
		},
		"string": {
			input: strVal,
			result: NullJobState{
				JobState: JobStateCompleted,
				Valid:    true,
			},
		},
		"*string": {
			input: &strVal,
			result: NullJobState{
				JobState: JobStateCompleted,
				Valid:    true,
			},
		},
		"invalid string": {
			input: "random value",
		},
		"[]byte": {
			input: []byte(JobStateProcessing.String()),
			result: NullJobState{
				JobState: JobStateProcessing,
				Valid:    true,
			},
		},
		"int": {
			input: intVal,
			result: NullJobState{
				JobState: JobStateFailed,
				Valid:    true,
			},
		},
		"*int": {
			input: &intVal,
			result: NullJobState{
				JobState: JobStateFailed,
				Valid:    true,
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
			input: nullJobState,
		},
		"int as []byte": { // must have --sqlnullint flag to get this feature.
			input: []byte("3"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			god := NewNullJobState(tc.input)
			assert.Equal(t, god, tc.result)
		})
	}
}
