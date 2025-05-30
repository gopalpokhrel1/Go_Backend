package storage

import "github.com/gopalpokhrel1/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetList() ([]types.Student, error)
}
