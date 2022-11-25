package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/PaulYakow/metrics-track/internal/entity"
	"github.com/PaulYakow/metrics-track/internal/usecase"
)

func TestServerUC_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockIServerRepo(ctrl)
	hasher := NewMockIHasher(ctrl)
	serverUC := usecase.NewServerUC(repo, hasher)

	metric := &entity.Metric{}

	hasher.EXPECT().Check(metric).Return(nil)
	hasher.EXPECT().ProcessPointer(metric)
	repo.EXPECT().Store(metric).Return(nil)
	err := serverUC.Save(metric)
	require.NoError(t, err)

	hasher.EXPECT().Check(metric).Return(errors.New("hash mismatch"))
	err = serverUC.Save(metric)
	require.Error(t, err)
}

func TestServerUC_SaveBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockIServerRepo(ctrl)
	hasher := NewMockIHasher(ctrl)
	serverUC := usecase.NewServerUC(repo, hasher)

	metricsIn := []entity.Metric{
		{ID: "first", Hash: "in"},
		{ID: "second", Hash: "in"}}
	metricsOut := []entity.Metric{
		{ID: "first", Hash: "out"},
		{ID: "second", Hash: "out"}}

	hasher.EXPECT().ProcessBatch(metricsIn).Return(metricsOut)
	repo.EXPECT().StoreBatch(metricsOut).Return(nil)
	err := serverUC.SaveBatch(metricsIn)
	require.NoError(t, err)

	hasher.EXPECT().ProcessBatch(metricsIn).Return(metricsOut)
	repo.EXPECT().StoreBatch(metricsOut).Return(errors.New("repo error"))
	err = serverUC.SaveBatch(metricsIn)
	require.Error(t, err)
}

func TestServerUC_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockIServerRepo(ctrl)
	hasher := NewMockIHasher(ctrl)
	serverUC := usecase.NewServerUC(repo, hasher)

	ctx := context.Background()
	metricIn := entity.Metric{ID: "metricIn", Hash: "zero"}
	metricOut := &entity.Metric{ID: "metricOut", Hash: "some_value"}

	repo.EXPECT().Read(ctx, metricIn).Return(metricOut, nil)
	hasher.EXPECT().ProcessPointer(metricOut)
	metric, err := serverUC.Get(ctx, metricIn)
	require.Equal(t, metricOut, metric)
	require.NoError(t, err)

	repo.EXPECT().Read(ctx, metricIn).Return(nil, errors.New("repo error"))
	metric, err = serverUC.Get(ctx, metricIn)
	require.Nil(t, metric)
	require.Error(t, err)
}

func TestServerUC_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockIServerRepo(ctrl)
	hasher := NewMockIHasher(ctrl)
	serverUC := usecase.NewServerUC(repo, hasher)

	ctx := context.Background()
	metricsOut := []entity.Metric{
		{ID: "first", Hash: "out"},
		{ID: "second", Hash: "out"}}

	repo.EXPECT().ReadAll(ctx).Return(metricsOut, nil)
	metrics, err := serverUC.GetAll(ctx)
	require.Equal(t, metricsOut, metrics)
	require.NoError(t, err)
}

func TestServerUC_CheckRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockIServerRepo(ctrl)
	hasher := NewMockIHasher(ctrl)
	serverUC := usecase.NewServerUC(repo, hasher)

	repo.EXPECT().CheckConnection().Return(errors.New("no connection"))
	err := serverUC.CheckRepo()
	require.Error(t, err)

	repo.EXPECT().CheckConnection().Return(nil)
	err = serverUC.CheckRepo()
	require.NoError(t, err)
}
