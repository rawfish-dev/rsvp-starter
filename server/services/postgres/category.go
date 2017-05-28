package postgres

import (
	"fmt"
	"strings"

	"bitbucket.org/rawfish-dev/wedding-rsvp/server/domain"
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
	category := &category{
		Tag: req.Tag,
	}

	err := s.gorpDB.Insert(category)
	if err != nil {
		if isCategoryTagUniqueConstraintError(err) {
			s.baseService.Warn("postgres service - unable to insert category with a duplicate tag")
			return nil, NewPostgresCategoryTagUniqueConstraintError()
		}

		s.baseService.Errorf("postgres service - unable to insert category due to %v", err)
		return nil, NewPostgresOperationError()
	}

	newCategory := &domain.Category{
		ID:  category.ID,
		Tag: category.Tag,
	}

	return newCategory, nil
}

func (s *service) FindCategoryByID(categoryID int64) (*domain.Category, error) {
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
			s.baseService.Warnf("postgres service - unable to find category with id %v", categoryID)
			return nil, NewPostgresRecordNotFoundError()
		}

		s.baseService.Errorf("postgres service - unable to find category with id %v due to %v", categoryID, err)
		return nil, NewPostgresOperationError()
	}

	domainCategory := &domain.Category{
		ID:    category.ID,
		Tag:   category.Tag,
		Total: category.Total,
	}

	return domainCategory, nil
}

func (s *service) FindAllCategories() ([]domain.Category, error) {
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
		s.baseService.Errorf("postgres service - unable to retrieve categories due to %v", err)
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
	category := &category{
		baseModel: baseModel{
			ID: domainCategory.ID,
		},
		Tag: domainCategory.Tag,
	}

	_, err := s.gorpDB.Update(category)
	if err != nil {
		s.baseService.Errorf("postgres service - unable to update category %+v due to %v", category, err)
		return nil, NewPostgresOperationError()
	}

	return domainCategory, nil
}

func (s *service) DeleteCategoryByID(categoryID int64) error {
	retrievedCategory, err := s.FindCategoryByID(categoryID)
	if err != nil {
		return err
	}

	category := &category{
		baseModel: baseModel{
			ID: retrievedCategory.ID,
		},
	}

	_, err = s.gorpDB.Delete(category)
	if err != nil {
		s.baseService.Errorf("postgres service - unable to delete category with id %v due to %v", categoryID, err)
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
