package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/TipsByPullak/phonepe-coding-round/constants"
	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow/handler"
)

type Pool struct {
	NumWorkers int
	startOnce  sync.Once
	wg         sync.WaitGroup
	Handler    *handler.Handler
}

func (wp *Pool) Start(ctx context.Context, logger *slog.Logger) error {
	wp.startOnce.Do(func() {
		for i := 1; i <= wp.NumWorkers; i++ {
			wp.wg.Add(1)
			go wp.worker(ctx, logger, i)
		}
	})
	return nil
}

func (wp *Pool) Wait() {
	wp.wg.Wait()
}

func (wp *Pool) worker(ctx context.Context, logger *slog.Logger, id int) {
	defer wp.wg.Done()

	ticker := time.NewTicker(constants.TickerTime)
	defer ticker.Stop()

	for {
		if ctx.Err() != nil {
			logger.InfoContext(ctx, "Worker halted before claiming new work", "id", id)
			return
		}

		dag := wp.Handler.GetDagForProcessing(ctx)
		logger.InfoContext(ctx, "dag seen>", "dag", dag)

		if dag == nil {
			// No work available right now; wait before polling again
			select {
			case <-ctx.Done():
				logger.InfoContext(ctx, "Worker halted while idle", "id", id)
				return
			case <-ticker.C:
				logger.InfoContext(ctx, "Worker %d waited for tick")
				continue
			}
		}
		// Do the work (failures here are isolated to this job)
		logger.InfoContext(ctx, "Worker processing job", "id", id, "workflowID", dag.WorkflowID)
		//if err := doWork(job); err != nil {
		//	wp.store.UpdateStatus(job.ID, "failed")
		//} else {
		//	wp.store.UpdateStatus(job.ID, "done")
		//}
	}
}
