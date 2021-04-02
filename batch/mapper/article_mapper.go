package mapper

import (
	"fmt"

	"github.com/KotaroKimura/note_analyzer/batch/client/mysql"
	"github.com/KotaroKimura/note_analyzer/batch/model"
)

type ArticleMapper struct {
	conn *mysql.Conn
}

func NewArticleMapper(conn *mysql.Conn) ArticleMapper {
	return ArticleMapper{
		conn: conn,
	}
}

func (m *ArticleMapper) FindByTitle(title string) ([]model.Article, error) {
	article := &model.Article{}

	query := fmt.Sprintf(`
SELECT
  *
FROM %s
WHERE title = ?
`, article.TableName())

	rows, err := m.conn.M.Query(query, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.CreatedAt,
			&article.UpdatedAt,
		); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (m *ArticleMapper) Insert(a *model.Article) error {
	query := fmt.Sprintf(`
INSERT INTO %s (
  title, created_at, updated_at
) VALUE (
  ?, ?, ?
)
`, a.TableName())

	stmt, err := m.conn.M.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(a.Title, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
