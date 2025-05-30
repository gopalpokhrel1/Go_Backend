package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/gopalpokhrel1/students-api/internal/config"
	"github.com/gopalpokhrel1/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	age INTEGER,
	email TEXT
   )`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?,?,?)")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, nil
	}

	return id, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id=? LIMIT 1 ")

	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Age, &student.Email)

	if err != nil {

		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student dat is not availble with this id")
		}
		return types.Student{}, fmt.Errorf("query error")
	}

	return student, nil
}

func (s *Sqlite) GetList() ([]types.Student, error) {

	stmt, err := s.Db.Prepare("SELECT id, name,age, email FROM students")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	defer row.Close()

	var students []types.Student

	for row.Next() {
		var student types.Student

		err := row.Scan(&student.Id, &student.Name, &student.Age, &student.Email)

		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil

}
