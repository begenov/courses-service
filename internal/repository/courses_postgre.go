package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/begenov/courses-service/internal/domain"

	"github.com/lib/pq"
)

type CoursesRepo struct {
	db *sql.DB
}

func NewCoursesRepo(db *sql.DB) *CoursesRepo {
	return &CoursesRepo{db: db}
}

func (r *CoursesRepo) Create(ctx context.Context, course domain.Courses) error {
	stmt := `INSERT INTO course(name, description, students, created_at) VALUES($1, $2, $3, $4)`
	if _, err := r.db.ExecContext(ctx, stmt, course.Name, course.Description, pq.Array(course.Students), course.CreatedAt); err != nil {
		return err
	}
	return nil
}

func (r *CoursesRepo) GetByID(ctx context.Context, id int) (domain.Courses, error) {
	var course domain.Courses
	stmt := `SELECT id, name, description, students, created_at FROM course WHERE id = $1`
	if err := r.db.QueryRowContext(ctx, stmt, id).Scan(&course.ID, &course.Name, &course.Description, pq.Array(&course.Students), &course.CreatedAt); err != nil {
		return course, err
	}

	return course, nil

}

func (r *CoursesRepo) Update(ctx context.Context, course domain.Courses) error {
	stmt := `UPDATE course SET name = $1, description = $2, students = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, stmt, course.Name, course.Description, pq.Array(course.Students), course.ID)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (r *CoursesRepo) Delele(ctx context.Context, id int) error {
	stmt := `DELETE FROM course WHERE id = $1`
	result, err := r.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return err
	}

	return nil
}

func (r *CoursesRepo) GetCoursesByIdStudent(ctx context.Context, studentId string) ([]domain.Courses, error) {
	var courses []domain.Courses
	stmt := `SELECT id, name, description, students, created_at FROM course WHERE $1 = ANY(students)`
	rows, err := r.db.QueryContext(ctx, stmt, studentId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var course domain.Courses
		err := rows.Scan(&course.ID, &course.Name, &course.Description, pq.Array(&course.Students), &course.CreatedAt)
		if err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}
	return courses, nil
}
