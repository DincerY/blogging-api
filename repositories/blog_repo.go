package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/DincerY/blogging-api/models"
)

type BlogRepository struct {
	DB *sql.DB
}

func NewBlogRepository(db *sql.DB) *BlogRepository {
	return &BlogRepository{DB: db}
}

func (r *BlogRepository) Create(blog models.Blog) (int, error) {
	var id int
	createBlogQuery := `
	INSERT INTO blog (title,content,category,tags,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;
	`
	tags := strings.Join(blog.Tags, ",")

	err := r.DB.QueryRow(createBlogQuery, blog.Title, blog.Content, blog.Category, tags, time.Now(), time.Time{}).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *BlogRepository) Update(blog models.Blog, id int) error {
	tags := strings.Join(blog.Tags, ",")
	updateBlogSql := `UPDATE blog SET title = $1,content = $2,category = $3,tags = $4,updated_at=$5 WHERE id = $6`
	res, err := r.DB.Exec(updateBlogSql, blog.Title, blog.Content, blog.Category, tags, time.Now(), id)
	if err != nil {
		return err
	}
	affectedRow, _ := res.RowsAffected()
	if affectedRow < 1 {
		return err
	}
	return nil
}

func (r *BlogRepository) Delete(id int) error {

	deleteBlogSql := `DELETE FROM blog WHERE id = $1`

	res, err := r.DB.Exec(deleteBlogSql, id)
	if err != nil {
		return err
	}
	affectedRow, err := res.RowsAffected()
	if affectedRow < 1 {
		return fmt.Errorf("no blogs have this id : %w", err)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *BlogRepository) GetAll() ([]models.Blog, error) {
	var blogs []models.Blog
	getAll := `SELECT * FROM blog`

	rows, err := r.DB.Query(getAll)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var b models.Blog
		var tags string
		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Content,
			&b.Category,
			&tags,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		b.Tags = strings.Split(tags, ",")
		blogs = append(blogs, b)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return blogs, nil
}

func (r *BlogRepository) GetById(id int) (models.Blog, error) {

	getByIdSql := `SELECT id, title, content, category, tags, created_at, updated_at FROM blog WHERE id = $1`
	var b models.Blog
	var tags string
	err := r.DB.QueryRow(getByIdSql, id).Scan(
		&b.ID,
		&b.Title,
		&b.Content,
		&b.Category,
		&tags,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return models.Blog{}, fmt.Errorf("blog was not found")
	}
	if err != nil {
		return models.Blog{}, err
	}
	b.Tags = strings.Split(tags, ",")

	return b, nil
}

func (r *BlogRepository) GetByCategory(category string) ([]models.Blog, error) {

	getByCategorySql := `SELECT * FROM blog WHERE category = $1`
	var blogs []models.Blog
	rows, err := r.DB.Query(getByCategorySql, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b models.Blog
		var tags string
		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Content,
			&b.Category,
			&tags,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		b.Tags = strings.Split(tags, ",")
		blogs = append(blogs, b)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return blogs, nil
}
