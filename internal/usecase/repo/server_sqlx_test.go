package repo

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/pkg/postgre/v2"
)

var testSqlx *ServerSqlxImpl

func TestMain(m *testing.M) {
	testDB, err := v2.New("postgresql://admin:root@localhost:5432/test_db", v2.ConnAttempts(1))
	if err != nil {
		log.Println(fmt.Errorf("skip repo tests because of: %w", err))
		os.Exit(0)
	}
	defer testDB.Close()

	testSqlx, err = NewSqlxImpl(testDB)
	if err != nil {
		log.Println(fmt.Errorf("repo tests - repo.New: %w", err))
	}

	os.Exit(m.Run())
}

func TestStore(t *testing.T) {
	var value = 77.7
	var delta int64 = 1

	mGauge := &entity.Metric{
		ID:    "testGauge",
		MType: "gauge",
		Value: &value,
		Hash:  "123",
	}

	mCounter := &entity.Metric{
		ID:    "testCounter",
		MType: "counter",
		Delta: &delta,
		Hash:  "987",
	}

	err := testSqlx.Store(mGauge)
	require.NoError(t, err)

	err = testSqlx.Store(mCounter)
	require.NoError(t, err)
}

func TestStoreBatch(t *testing.T) {
	var value = 77.7
	var delta int64 = 1

	metrics := []entity.Metric{
		{
			ID:    "testGaugeBatch",
			MType: "gauge",
			Value: &value,
		},
		{
			ID:    "testCounterBatch",
			MType: "counter",
			Delta: &delta,
		},
	}

	err := testSqlx.StoreBatch(metrics)
	require.NoError(t, err)
}

func TestRead(t *testing.T) {
	mGauge := entity.Metric{
		ID:    "testGauge",
		MType: "gauge",
	}

	mCounter := entity.Metric{
		ID:    "testCounter",
		MType: "counter",
	}

	gauge, err := testSqlx.Read(context.Background(), mGauge)
	require.NoError(t, err)
	require.NotEmpty(t, gauge)

	counter, err := testSqlx.Read(context.Background(), mCounter)
	require.NoError(t, err)
	require.NotEmpty(t, counter)
}

func TestReadAll(t *testing.T) {
	metrics, err := testSqlx.ReadAll(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, metrics)
}

func TestCheckConnection(t *testing.T) {
	err := testSqlx.CheckConnection()
	require.NoError(t, err)
}
