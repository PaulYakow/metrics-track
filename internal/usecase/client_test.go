package usecase_test

//type test struct {
//	name string
//	mock func()
//	res  interface{}
//	err  error
//}

//func TestClientPoll(t *testing.T) {
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//
//	repo := NewMockIClientMemory(mockCtrl)
//	hasher := NewMockIHasher(mockCtrl)
//	gather := NewMockIClientGather(mockCtrl)
//	clientUC := NewClientUC(context.Background(), repo, hasher)
//
//	gather.EXPECT().Update().Return(map[string]*entity.Metric{})
//	repo.EXPECT().Store(gather.Update())
//	clientUC.Poll()
//	metrics := repo.ReadAll()
//	require.Equal(t, map[string]*entity.Metric{}, metrics)
//}
