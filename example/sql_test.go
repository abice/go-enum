// +build example

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
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dataSourceName     = "root:pass@tcp(localhost:3306)/database"
	hardcodedProjectID = 1
)

func ExampleSQL() {
	conn, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	var status *ProjectStatus

	// Set status InWork
	err = setProjectStatus(conn, ProjectStatusInWork)
	if err != nil {
		panic(err)
	}

	// Get inWork status
	status, err = getProjectStatus(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println(status)

	// Set status completed
	err = setProjectStatus(conn, ProjectStatusCompleted)
	if err != nil {
		panic(err)
	}

	// Get status completed
	status, err = getProjectStatus(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println(status)

	// Output:
	// inWork
	// completed
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
