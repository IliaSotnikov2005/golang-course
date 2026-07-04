package app

import (
	"context"
	"time"
)

func (a *App) runBackgroundUpdater(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.log.Debug("background updater: fetching subscriptions")

			subs, err := a.subscriberAdapter.GetSubscriptions(ctx)
			if err != nil {
				a.log.Error("failed to fetch subscriptions", "err", err)
				continue
			}

			for _, s := range subs {
				_ = a.taskDispatcher.Dispatch(ctx, s.Owner, s.Repo)
			}
			a.log.Info("background tasks sent", "count", len(subs))

		case <-ctx.Done():
			return
		}
	}
}
