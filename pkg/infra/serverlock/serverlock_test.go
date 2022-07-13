package serverlock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/sqlstore"
)

func createTestableServerLock(t *testing.T) *ServerLockService {
	t.Helper()

	sqlstore := sqlstore.InitTestDB(t)

	return &ServerLockService{
		SQLStore: sqlstore,
		log:      log.New("test-logger"),
	}
}

func TestServerLock(t *testing.T) {
	sl := createTestableServerLock(t)
	operationUID := "test-operation"

	first, err := sl.getOrCreate(context.Background(), operationUID)
	require.NoError(t, err)

	t.Run("trying to create three new row locks", func(t *testing.T) {
		expectedLastExecution := first.LastExecution
		var latest *serverLock

		for i := 0; i < 3; i++ {
			latest, err = sl.getOrCreate(context.Background(), operationUID)
			require.NoError(t, err)
			assert.Equal(t, operationUID, first.OperationUID)
			assert.Equal(t, int64(1), first.Id)
		}

		assert.Equal(t,
			expectedLastExecution,
			latest.LastExecution,
			"latest execution should not have changed")
	})

	t.Run("create lock on first row", func(t *testing.T) {
		gotLock, _, err := sl.acquireLock(context.Background(), first)
		require.NoError(t, err)
		assert.True(t, gotLock)

		gotLock, _, err = sl.acquireLock(context.Background(), first)
		require.NoError(t, err)
		assert.False(t, gotLock)
	})

	t.Run("create lock and then release it", func(t *testing.T) {
		lock, err := sl.getOrCreate(context.Background(), operationUID)

		gotLock, newVersion, err := sl.acquireLock(context.Background(), lock)
		require.NoError(t, err)
		assert.True(t, gotLock)

		err = sl.releaseLock(context.Background(), newVersion, lock)
		require.NoError(t, err)

		// and now we can acquire it again
		gotLock2, _, err2 := sl.acquireLock(context.Background(), lock)
		require.NoError(t, err2)
		assert.True(t, gotLock2)
	})
}
