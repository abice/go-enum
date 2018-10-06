// +build go1.11
//go:generate mockgen -destination sql_mock_test.go -package example database/sql/driver Conn,Driver,Stmt,Result,Rows

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
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	dataSourceName     = "root:pass@tcp(localhost:3306)/database"
	hardcodedProjectID = 1
)

func TestExampleSQL(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockDriver := NewMockDriver(controller)
	mockConn := NewMockConn(controller)
	mockStmt := NewMockStmt(controller)
	mockResult := NewMockResult(controller)
	mockRows := NewMockRows(controller)

	// Set up the database for this test case
	mockDriver.EXPECT().Open(dataSourceName).Return(mockConn, nil)
	gomock.InOrder(
		// Update To InWork
		mockConn.EXPECT().Prepare("UPDATE project SET status = ? WHERE id = ?").Return(mockStmt, nil),
		mockStmt.EXPECT().NumInput().Return(2),
		mockStmt.EXPECT().Exec(gomock.Any()).Return(mockResult, nil),
		mockStmt.EXPECT().Close().Return(nil),
		// Select
		mockConn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mockStmt, nil),
		mockStmt.EXPECT().NumInput().Return(1),
		mockStmt.EXPECT().Query(gomock.Any()).Return(mockRows, nil),
		mockRows.EXPECT().Columns().Return([]string{`status`}),
		mockRows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{ProjectStatusInWork.String()}).Return(nil),
		mockRows.EXPECT().Close().Return(nil),
		mockStmt.EXPECT().Close().Return(nil),
		// Update to Completed
		mockConn.EXPECT().Prepare("UPDATE project SET status = ? WHERE id = ?").Return(mockStmt, nil),
		mockStmt.EXPECT().NumInput().Return(2),
		mockStmt.EXPECT().Exec(gomock.Any()).Return(mockResult, nil),
		mockStmt.EXPECT().Close().Return(nil),
		// Select
		mockConn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mockStmt, nil),
		mockStmt.EXPECT().NumInput().Return(1),
		mockStmt.EXPECT().Query(gomock.Any()).Return(mockRows, nil),
		mockRows.EXPECT().Columns().Return([]string{`status`}),
		mockRows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{[]byte(ProjectStatusCompleted.String())}).Return(nil),
		mockRows.EXPECT().Close().Return(nil),
		mockStmt.EXPECT().Close().Return(nil),
		// Select Nil Value response
		mockConn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mockStmt, nil),
		mockStmt.EXPECT().NumInput().Return(1),
		mockStmt.EXPECT().Query(gomock.Any()).Return(mockRows, nil),
		mockRows.EXPECT().Columns().Return([]string{`status`}),
		mockRows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{nil}).Return(nil),
		mockRows.EXPECT().Close().Return(nil),
		mockStmt.EXPECT().Close().Return(nil),
		// Select Non Status Value response
		mockConn.EXPECT().Prepare("SELECT status FROM project WHERE id = ?").Return(mockStmt, nil),
		mockStmt.EXPECT().NumInput().Return(1),
		mockStmt.EXPECT().Query(gomock.Any()).Return(mockRows, nil),
		mockRows.EXPECT().Columns().Return([]string{`status`}),
		mockRows.EXPECT().Next(gomock.Any()).SetArg(0, []driver.Value{"NotAStatus"}).Return(nil),
		mockRows.EXPECT().Columns().Return([]string{`status`}),
		mockRows.EXPECT().Close().Return(nil),
		mockStmt.EXPECT().Close().Return(nil),
	)

	sql.Register("mock", mockDriver)
	var status *ProjectStatus

	conn, err := sql.Open("mock", dataSourceName)
	require.NoError(t, err, "failed opening mock db")

	// Set status InWork
	err = setProjectStatus(conn, ProjectStatusInWork)
	require.NoError(t, err, "failed setting project status")

	// Get inWork status
	status, err = getProjectStatus(conn)
	require.NoError(t, err, "failed getting project status")
	require.Equal(t, ProjectStatusInWork, *status)

	// Set status completed
	err = setProjectStatus(conn, ProjectStatusCompleted)
	require.NoError(t, err, "failed setting project status")

	// Get status completed
	status, err = getProjectStatus(conn)
	require.NoError(t, err, "failed getting project status")
	require.Equal(t, ProjectStatusCompleted, *status)

	// Get status nil
	status, err = getProjectStatus(conn)
	require.NoError(t, err, "Nil Values do not error")
	require.Equal(t, ProjectStatus(0), *status)

	// Get Non Status
	_, err = getProjectStatus(conn)
	require.EqualError(t, err, "sql: Scan error on column index 0, name \"status\": NotAStatus is not a valid ProjectStatus", "should have failed getting project status")

}

func getProjectStatus(db *sql.DB) (*ProjectStatus, error) {
	var status ProjectStatus
	err := db.QueryRow(`SELECT status FROM project WHERE id = ?`, hardcodedProjectID).Scan(&status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func setProjectStatus(db *sql.DB, status ProjectStatus) error {
	_, err := db.Exec(`UPDATE project SET status = ? WHERE id = ?`, status, hardcodedProjectID)
	return err
}
