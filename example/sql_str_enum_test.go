package example

import (
	"database/sql"
	driver "database/sql/driver"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExampleSQLStrOnly(t *testing.T) {

	tests := map[string]struct {
		setupMock func(t testing.TB, mocks *MockSQL)
		tester    func(t testing.TB, db *sql.DB)
	}{
		"NonNullable select nil": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM job WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getJobState(conn)
				require.NoError(t, err, "failed getting project status")
				require.Equal(t, JobStatePending, status)
			},
		},
		"Nullable Select Null": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM job WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullJobState(conn)
				require.NoError(t, err, "failed getting project status")
				require.False(t, status.Valid)
			},
		},
		"Select a string": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM job WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{JobStateProcessing.String()}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getJobState(conn)
				require.NoError(t, err, "failed getting project status")
				require.Equal(t, JobStateProcessing, status)
			},
		},
		"Nullable select an int": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM job WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullJobState(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, JobStateCompleted, status.JobState)
			},
		},
		"Nullable select an int64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM job WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullJobState(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, JobStateCompleted, status.JobState)
			},
		},
		"Nullable select an uint": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM job WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullJobState(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, JobStateCompleted, status.JobState)
			},
		},
		"Nullable select an uint64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM job WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullJobState(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, JobStateCompleted, status.JobState)
			},
		},
		"Nullable select an *uint64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := uint64(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM job WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullJobState(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, JobStateCompleted, status.JobState)
			},
		},
		"Nullable select an *int64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := int64(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM job WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullJobState(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, JobStateCompleted, status.JobState)
			},
		},
		"Nullable select an *uint": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := uint(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM job WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullJobState(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, JobStateCompleted, status.JobState)
			},
		},
		"standard update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare("Update job SET status = ? WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(JobStateFailed.String(), hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setJobState(conn, JobStateFailed)
				require.NoError(t, err, "failed updating project status")
			},
		},
		"nullable string update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare("Update job SET status = ? WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(JobStateFailed.String(), hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setJobState(conn, NewNullJobState(JobStateFailed))
				require.NoError(t, err, "failed updating project status")
			},
		},
		"nullable invalid string update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare("Update job SET status = ? WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(nil, hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setJobState(conn, NullJobState{})
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

func getJobState(db *sql.DB) (state JobState, err error) {
	err = db.QueryRow(`SELECT status FROM job WHERE id = ?`, hardcodedProjectID).Scan(&state)
	return
}

func getNullJobState(db *sql.DB) (state NullJobState, err error) {
	err = db.QueryRow(`SELECT status FROM job WHERE id = ?`, hardcodedProjectID).Scan(&state)
	return
}

func setJobState(db *sql.DB, state interface{}) error {
	_, err := db.Exec(`Update job SET status = ? WHERE id = ?`, state, hardcodedProjectID)
	return err
}

func TestSQLStrExtras(t *testing.T) {

	assert.Equal(t, "JobState(22)", JobState(22).String(), "String value is not correct")

	_, err := ParseJobState(`NotAStatus`)
	assert.Error(t, err, "Should have had an error parsing a non status")

	var (
		intVal  int      = 3
		strVal  string   = "completed"
		enumVal JobState = JobStateCompleted
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			status := NewNullJobState(tc.input)
			assert.Equal(t, status, tc.result)

		})
	}

}
