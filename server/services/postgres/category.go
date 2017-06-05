package postgres

import (
	"fmt"
	"strings"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
)

type category struct {
	baseModel
	Tag string `db:"tag"`
}

type categoryAggregate struct {
	category
	Total int `db:"total"`
}

var (
	categoryColumns = strings.Join([]string{
		"id",
		"tag",
		"created_at",
		"updated_at",
	}, ",")
)

func (s *service) InsertCategory(req *domain.CategoryCreateRequest) (*domain.Category, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	category := &category{
		Tag: req.Tag,
	}

	err := s.gorpDB.Insert(category)
	if err != nil {
		if isCategoryTagUniqueConstraintError(err) {
			ctxLogger.Warn("postgres service - unable to insert category with a duplicate tag")
			return nil, NewPostgresCategoryTagUniqueConstraintError()
		}

		ctxLogger.Errorf("postgres service - unable to insert category due to %v", err)
		return nil, NewPostgresOperationError()
	}

	newCategory := &domain.Category{
		ID:  category.ID,
		Tag: category.Tag,
	}

	return newCategory, nil
}

func (s *service) FindCategoryByID(categoryID int64) (*domain.Category, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	query := fmt.Sprintf(`
		SELECT %v, COUNT(invitations.id) as total
		FROM categories
		LEFT JOIN invitations
        ON categories.id=invitations.category_id 
		WHERE categories.id=$1
		GROUP BY categories.id
	`, prependColumnsForJoin())

	var category categoryAggregate

	err := s.gorpDB.SelectOne(&category, query, categoryID)
	if err != nil {
		if isNotFoundError(err) {
			ctxLogger.Warnf("postgres service - unable to find category with id %v", categoryID)
			return nil, NewPostgresRecordNotFoundError()
		}

		ctxLogger.Errorf("postgres service - unable to find category with id %v due to %v", categoryID, err)
		return nil, NewPostgresOperationError()
	}

	domainCategory := &domain.Category{
		ID:    category.ID,
		Tag:   category.Tag,
		Total: category.Total,
	}

	return domainCategory, nil
}

func (s *service) ListCategories() ([]domain.Category, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	query := fmt.Sprintf(`
		SELECT %v, COUNT(invitations.id) as total
		FROM categories
		LEFT JOIN invitations
        ON categories.id=invitations.category_id
		GROUP BY categories.id
        ORDER BY tag DESC
	`, prependColumnsForJoin())

	var categories []categoryAggregate

	_, err := s.gorpDB.Select(&categories, query)
	if err != nil {
		ctxLogger.Errorf("postgres service - unable to retrieve categories due to %v", err)
		return nil, NewPostgresOperationError()
	}

	domainCategories := make([]domain.Category, len(categories))
	for idx := range categories {
		domainCategories[idx] = domain.Category{
			ID:    categories[idx].ID,
			Tag:   categories[idx].Tag,
			Total: categories[idx].Total,
		}
	}

	return domainCategories, nil
}

func (s *service) UpdateCategory(domainCategory *domain.Category) (*domain.Category, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	category := &category{
		baseModel: baseModel{
			ID: domainCategory.ID,
		},
		Tag: domainCategory.Tag,
	}

	_, err := s.gorpDB.Update(category)
	if err != nil {
		ctxLogger.Errorf("postgres service - unable to update category %+v due to %v", category, err)
		return nil, NewPostgresOperationError()
	}

	return domainCategory, nil
}

func (s *service) DeleteCategory(domainCategory *domain.Category) error {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	_, err := s.gorpDB.Delete(domainCategory)
	if err != nil {
		ctxLogger.Errorf("postgres service - unable to delete category with id %v due to %v", domainCategory.ID, err)
		return NewPostgresOperationError()
	}

	return nil
}

func prependColumnsForJoin() string {
	columns := strings.Split(categoryColumns, ",")
	prependedColumns := make([]string, len(columns))
	for idx := range columns {
		prependedColumns[idx] = "categories." + columns[idx]
	}

	return strings.Join(prependedColumns, ",")
}
