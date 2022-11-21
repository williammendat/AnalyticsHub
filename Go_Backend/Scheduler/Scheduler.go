package scheduler

type taskFunc func() 

type (
	Scheduler struct {
		tasks []taskFunc
	}
)

func NewScheduler(tasks ...taskFunc) *Scheduler{
	scheduler := &Scheduler{
		tasks: make([]taskFunc, 0, len(tasks)),
	}

	scheduler.tasks = append(scheduler.tasks, tasks...)

	return scheduler
}

func (s *Scheduler) StartTasksAsync() {

	for _, task := range s.tasks {
		go task()
	}
}