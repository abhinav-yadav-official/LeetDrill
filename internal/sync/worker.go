package sync

import (
	"context"
	"log"
	"time"

	"leetdrill/internal/leetcode"
	"leetdrill/internal/store"
)

type RecentWorker struct {
	Store    *store.Store
	Client   leetcodeClient
	Interval time.Duration
	Logger   *log.Logger
}

func (w *RecentWorker) Start(ctx context.Context) {
	interval := w.Interval
	if interval <= 0 {
		interval = 30 * time.Minute
	}
	logger := w.Logger
	if logger == nil {
		logger = log.Default()
	}
	go func() {
		w.runOnce(ctx, logger)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.runOnce(ctx, logger)
			}
		}
	}()
}

func (w *RecentWorker) runOnce(ctx context.Context, logger *log.Logger) {
	if w.Store == nil {
		return
	}
	client := w.Client
	if client == nil {
		client = leetcode.New()
	}
	users, err := store.ListUsersForRecentSync(ctx, w.Store.DB())
	if err != nil {
		logger.Printf("recent sync users: %v", err)
		return
	}
	for _, user := range users {
		importer := &ColdStartImporter{
			Store:              w.Store,
			Client:             client,
			MaxSubmissionPages: 0,
			SubmissionPageSize: 0,
		}
		result, err := importer.Run(ctx, user.ID, user.LeetcodeUsername)
		if err != nil {
			logger.Printf("recent sync user_id=%d username=%s: %v", user.ID, user.LeetcodeUsername, err)
			continue
		}
		if result.RecentImported > 0 {
			logger.Printf("recent sync user_id=%d username=%s imported=%d", user.ID, user.LeetcodeUsername, result.RecentImported)
		}
		time.Sleep(1100 * time.Millisecond)
	}
}
