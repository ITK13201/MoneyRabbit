package persistence_test

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/itk13201/money-rabbit/ent"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
)

// setupTestDB starts a MySQL container and returns a ready ent.Client.
// The container and client are cleaned up automatically via t.Cleanup.
func setupTestDB(t *testing.T) *ent.Client {
	t.Helper()
	ctx := context.Background()

	container, err := tcmysql.Run(ctx, "mysql:8.4",
		tcmysql.WithDatabase("testdb"),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("password"),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "3306")
	require.NoError(t, err)

	dsn := fmt.Sprintf("root:password@tcp(%s:%s)/testdb?parseTime=true", host, port.Port())
	client, err := ent.Open("mysql", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { client.Close() })

	err = client.Schema.Create(ctx)
	require.NoError(t, err)

	return client
}
