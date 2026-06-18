package service_test

import (
	"context"
	"testing"
	"time"

	"cesizen/internal/domain"
	"cesizen/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockContentRepo struct {
	contents    []domain.Content
	findByIDErr error
	createErr   error
	updateErr   error
	deleteErr   error
}

func (m *mockContentRepo) FindAllPublished(_ context.Context) ([]domain.Content, error) {
	var out []domain.Content
	for _, c := range m.contents {
		if c.IsPublished {
			out = append(out, c)
		}
	}
	return out, nil
}

func (m *mockContentRepo) FindAll(_ context.Context) ([]domain.Content, error) {
	return m.contents, nil
}

func (m *mockContentRepo) FindByID(_ context.Context, id int) (*domain.Content, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	for _, c := range m.contents {
		if c.ID == id {
			cp := c
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockContentRepo) Create(_ context.Context, c *domain.Content) error {
	if m.createErr != nil {
		return m.createErr
	}
	c.ID = 42
	c.IsPublished = true
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (m *mockContentRepo) Update(_ context.Context, _ *domain.Content) error {
	return m.updateErr
}

func (m *mockContentRepo) Delete(_ context.Context, _ int) error {
	return m.deleteErr
}

var testContents = []domain.Content{
	{ID: 1, Title: "Article publié", Content: "Contenu 1", Author: "Alice", IsPublished: true},
	{ID: 2, Title: "Brouillon", Content: "Contenu 2", Author: "Bob", IsPublished: false},
}

func TestContentService_ListPublished(t *testing.T) {
	svc := service.NewContentService(&mockContentRepo{contents: testContents})
	list, err := svc.ListPublished(context.Background())
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "Article publié", list[0].Title)
}

func TestContentService_ListAll(t *testing.T) {
	svc := service.NewContentService(&mockContentRepo{contents: testContents})
	list, err := svc.ListAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestContentService_GetByID_Success(t *testing.T) {
	svc := service.NewContentService(&mockContentRepo{contents: testContents})
	c, err := svc.GetByID(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "Article publié", c.Title)
}

func TestContentService_GetByID_NotFound(t *testing.T) {
	svc := service.NewContentService(&mockContentRepo{contents: testContents})
	_, err := svc.GetByID(context.Background(), 99)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestContentService_Create(t *testing.T) {
	svc := service.NewContentService(&mockContentRepo{})
	c, err := svc.Create(context.Background(), domain.CreateContentInput{
		Title: "Nouveau", Content: "Corps", Author: "Charlie",
	})
	require.NoError(t, err)
	assert.Equal(t, 42, c.ID)
	assert.Equal(t, "Nouveau", c.Title)
	assert.True(t, c.IsPublished)
}

func TestContentService_Update_Success(t *testing.T) {
	svc := service.NewContentService(&mockContentRepo{contents: testContents})
	c, err := svc.Update(context.Background(), 1, domain.UpdateContentInput{
		Title: "Modifié", Content: "Nouveau corps", Author: "Alice", IsPublished: false,
	})
	require.NoError(t, err)
	assert.Equal(t, "Modifié", c.Title)
	assert.False(t, c.IsPublished)
}

func TestContentService_Update_NotFound(t *testing.T) {
	svc := service.NewContentService(&mockContentRepo{contents: testContents})
	_, err := svc.Update(context.Background(), 99, domain.UpdateContentInput{Title: "X", Content: "Y"})
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestContentService_Delete(t *testing.T) {
	svc := service.NewContentService(&mockContentRepo{})
	assert.NoError(t, svc.Delete(context.Background(), 1))
}
