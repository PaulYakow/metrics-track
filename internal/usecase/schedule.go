package usecase

// Реализация планировщика

type Schedule struct {
	fileRepo   IServerFile
	memoryRepo IServerMemory
}

func NewScheduleUC(file IServerFile, memory IServerMemory) *Schedule {
	return &Schedule{
		fileRepo:   file,
		memoryRepo: memory,
	}
}

func (s *Schedule) RunStoring() {
	s.fileRepo.SaveMetrics(s.memoryRepo.ReadAll())
}

func (s *Schedule) InitMetrics() {
	s.memoryRepo.InitializeMetrics(s.fileRepo.ReadMetrics())
}
