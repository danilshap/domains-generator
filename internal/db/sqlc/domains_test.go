package db

import (
	"context"
	"testing"

	"github.com/danilshap/domains-generator/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateDomain(t *testing.T) {
	arg := CreateDomainParams{
		Name:     utils.RandomName(),
		Provider: utils.RandomProvider(),
		Status:   1,
	}

	domain, err := testQueries.CreateDomain(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, domain)

	require.Equal(t, arg.Name, domain.Name)
	require.Equal(t, arg.Provider, domain.Provider)
	require.Equal(t, arg.Status, domain.Status)
	require.False(t, domain.IsDeleted.Bool)
	require.NotZero(t, domain.ID)
	require.NotZero(t, domain.CreatedAt)
}

func TestGetDomainByName(t *testing.T) {
	// Сначала создаем домен
	createdDomain := createRandomDomain(t)

	// Получаем домен по имени
	domain, err := testQueries.GetDomainByName(context.Background(), createdDomain.Name)
	require.NoError(t, err)
	require.NotEmpty(t, domain)

	require.Equal(t, createdDomain.ID, domain.ID)
	require.Equal(t, createdDomain.Name, domain.Name)
	require.Equal(t, createdDomain.Provider, domain.Provider)
	require.Equal(t, createdDomain.Status, domain.Status)
	require.Equal(t, createdDomain.IsDeleted, domain.IsDeleted)
}

func TestGetAllDomains(t *testing.T) {
	// Создаем несколько доменов
	for i := 0; i < 5; i++ {
		createRandomDomain(t)
	}

	domains, err := testQueries.GetAllDomains(context.Background(), GetAllDomainsParams{
		Limit:  5,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, domains)
	require.GreaterOrEqual(t, len(domains), 5)

	for _, domain := range domains {
		require.NotEmpty(t, domain.ID)
		require.NotEmpty(t, domain.Name)
		require.NotEmpty(t, domain.Provider)
		require.NotZero(t, domain.Status)
		require.NotZero(t, domain.CreatedAt)
	}
}

func TestSetDomainStatus(t *testing.T) {
	domain := createRandomDomain(t)

	arg := SetDomainStatusParams{
		ID:     domain.ID,
		Status: 2,
	}

	err := testQueries.SetDomainStatus(context.Background(), arg)
	require.NoError(t, err)

	// Проверяем что статус обновился
	updatedDomain, err := testQueries.GetDomainByName(context.Background(), domain.Name)
	require.NoError(t, err)
	require.Equal(t, arg.Status, updatedDomain.Status)
}

func TestDeleteDomain(t *testing.T) {
	domain := createRandomDomain(t)

	err := testQueries.DeleteDomain(context.Background(), domain.ID)
	require.NoError(t, err)

	// Проверяем что домен удален
	_, err = testQueries.GetDomainByName(context.Background(), domain.Name)
	require.Error(t, err) // Должна быть ошибка, так как домен удален
}

// Вспомогательная функция для создания случайного домена
func createRandomDomain(t *testing.T) Domain {
	arg := CreateDomainParams{
		Name:     utils.RandomName(),
		Provider: utils.RandomProvider(),
		Status:   1,
	}

	domain, err := testQueries.CreateDomain(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, domain)

	return domain
}
