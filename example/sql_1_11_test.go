//go:generate ../bin/mockgen -destination sql_mock_test.go -package example database/sql/driver Conn,Driver,Stmt,Result,Rows

/*
This example shows the conversion of enumerations between GO and SQL database.
You can run this example with the command: `go test -tags example github.com/abice/go-enum/example -v -run ^ExampleSQL$`
Don't forget to change the constant "dataSourceName" if necessary and apply the sql query.

SQL query to create a database and fill the initial data:

	CREATE TABLE project
	(
		id INT PRIMARY KEY AUTO_INCREMENT,
		status ENUM('pending', 'inWork', 'completed', 'rejected')
	);
	INSERT INTO project (`id`, `status`) VALUES (1, 'pending')
*/

package example

import (
	"database/sql"
	driver "database/sql/driver"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	dataSourceName     = "root:pass@tcp(localhost:3306)/database"
	hardcodedProjectID = 1
)

// A Matcher is a representation of a class of values.
// It is used to represent the valid or expected arguments to a mocked method.
type DriverValueMatcher struct {
	values []driver.Value
}

func MatchesValues(vals ...driver.Value) *DriverValueMatcher {
	return &DriverValueMatcher{
		values: vals,
	}
}

// Matches returns whether x is a match.
func (d *DriverValueMatcher) Matches(x interface{}) bool {

	switch values := x.(type) {
	case []driver.Value:
		if len(values) != len(d.values) {
			return false
		}
		for i, value := range values {
			if !assert.ObjectsAreEqualValues(d.values[i], value) {
				fmt.Printf("%v != %v\n", value, d.values[i])
				return false
			}
		}
	default:
		return false
	}

	return true
}

// String describes what the matcher matches.
func (d *DriverValueMatcher) String() string {
	return fmt.Sprintf("%v", d.values)
}

type MockSQL struct {
	Driver *MockDriver
	Conn   *MockConn
	Stmt   *MockStmt
	Result *MockResult
	Rows   *MockRows
}

func WithMockSQL(t testing.TB) (*MockSQL, func()) {
	t.Helper()

	ctrl := gomock.NewController(t)

	mocks := &MockSQL{
		Conn:   NewMockConn(ctrl),
		Stmt:   NewMockStmt(ctrl),
		Result: NewMockResult(ctrl),
		Rows:   NewMockRows(ctrl),
	}

	return mocks, func() {
		ctrl.Finish()
	}
}

func TestExampleSQL(t *testing.T) {

	tests := map[string]struct {
		setupMock func(t testing.TB, mocks *MockSQL)
		tester    func(t testing.TB, db *sql.DB)
	}{
		"NonNullable select nil": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getProjectStatus(conn)
				require.NoError(t, err, "failed getting project status")
				require.Equal(t, ProjectStatusPending, *status)
			},
		},
		"NullableStr Select Null": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullStrProjectStatus(conn)
				require.NoError(t, err, "failed getting project status")
				require.False(t, status.Valid)
			},
		},
		"Nullable Select Null": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullProjectStatus(conn)
				require.NoError(t, err, "failed getting project status")
				require.False(t, status.Valid)
			},
		},
		"Select a string": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(1),
					mocks.Stmt.EXPECT().Query(gomock.Any()).Return(mocks.Rows, nil),
					mocks.Rows.EXPECT().Columns().Return([]string{`status`}),
					mocks.Rows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{ProjectStatusInWork.String()}).Return(nil),
					mocks.Rows.EXPECT().Close().Return(nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				status, err := getProjectStatus(conn)
				require.NoError(t, err, "failed getting project status")
				require.Equal(t, ProjectStatusInWork, *status)
			},
		},
		"Nullable select an int": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullProjectStatus(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ProjectStatusCompleted, status.ProjectStatus)
			},
		},
		"Nullable select an int64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullProjectStatus(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ProjectStatusCompleted, status.ProjectStatus)
			},
		},
		"Nullable select an uint": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullProjectStatus(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ProjectStatusCompleted, status.ProjectStatus)
			},
		},
		"Nullable select an uint64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullProjectStatus(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ProjectStatusCompleted, status.ProjectStatus)
			},
		},
		"Nullable select an *uint64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := uint64(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullProjectStatus(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ProjectStatusCompleted, status.ProjectStatus)
			},
		},
		"Nullable select an *int64": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := int64(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullProjectStatus(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ProjectStatusCompleted, status.ProjectStatus)
			},
		},
		"Nullable select an *uint": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				val := uint(2)
				gomock.InOrder(
					// Select In Work
					mocks.Conn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mocks.Stmt, nil),
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
				status, err := getNullProjectStatus(conn)
				require.NoError(t, err, "failed getting project status")
				require.True(t, status.Valid)
				require.Equal(t, ProjectStatusCompleted, status.ProjectStatus)
			},
		},
		"standard update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare("UPDATE project SET status = ? WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(ProjectStatusRejected.String(), hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setProjectStatus(conn, ProjectStatusRejected)
				require.NoError(t, err, "failed updating project status")
			},
		},
		"nullable update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare("UPDATE project SET status = ? WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(3, hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setProjectStatus(conn, NewNullProjectStatus(ProjectStatusRejected))
				require.NoError(t, err, "failed updating project status")
			},
		},
		"nullable invalid update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare("UPDATE project SET status = ? WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(nil, hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setProjectStatus(conn, NullProjectStatus{})
				require.NoError(t, err, "failed updating project status")
			},
		},
		"nullable string update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare("UPDATE project SET status = ? WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(ProjectStatusRejected.String(), hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setProjectStatus(conn, NewNullProjectStatusStr(ProjectStatusRejected))
				require.NoError(t, err, "failed updating project status")
			},
		},
		"nullable invalid string update": {
			setupMock: func(t testing.TB, mocks *MockSQL) {
				gomock.InOrder(
					// Update value
					mocks.Conn.EXPECT().Prepare("UPDATE project SET status = ? WHERE id = ?").Return(mocks.Stmt, nil),
					mocks.Stmt.EXPECT().NumInput().Return(2),
					mocks.Stmt.EXPECT().Exec(MatchesValues(nil, hardcodedProjectID)).Return(mocks.Result, nil),
					mocks.Stmt.EXPECT().Close().Return(nil),
				)
			},
			tester: func(t testing.TB, conn *sql.DB) {
				// Get inWork status
				err := setProjectStatus(conn, NullProjectStatusStr{})
				require.NoError(t, err, "failed updating project status")
			},
		},
	}

	driverctrl := gomock.NewController(t)
	driver := NewMockDriver(driverctrl)
	defer func() {
		driverctrl.Finish()
	}()

	sql.Register("mock", driver)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			require.NotNil(t, tc.setupMock)
			require.NotNil(t, tc.tester)

			mocks, validate := WithMockSQL(t)
			defer validate()

			driver.EXPECT().Open(dataSourceName).Return(mocks.Conn, nil)

			conn, err := sql.Open("mock", dataSourceName)
			require.NoError(t, err, "failed opening mock db")

			tc.setupMock(t, mocks)

			tc.tester(t, conn)

		})
	}
}

func getProjectStatus(db *sql.DB) (*ProjectStatus, error) {
	var status ProjectStatus
	err := db.QueryRow(`SELECT status FROM project WHERE id = ?`, hardcodedProjectID).Scan(&status)
	if err != nil {
		return nil, err
	}

	return status.Ptr(), nil
}

func getNullStrProjectStatus(db *sql.DB) (status NullProjectStatusStr, err error) {
	err = db.QueryRow(`SELECT status FROM project WHERE id = ?`, hardcodedProjectID).Scan(&status)
	return
}

func getNullProjectStatus(db *sql.DB) (status NullProjectStatus, err error) {
	err = db.QueryRow(`SELECT status FROM project WHERE id = ?`, hardcodedProjectID).Scan(&status)
	return
}

func setProjectStatus(db *sql.DB, status interface{}) error {
	_, err := db.Exec(`UPDATE project SET status = ? WHERE id = ?`, status, hardcodedProjectID)
	return err
}
