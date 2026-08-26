package jobs

import (
	"context"
	"flash-ticket/internal/bookings"
	"log"
	"time"
)

func StartSeatExpiryWorker(ctx context.Context, repo bookings.BookingRepository, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping seat expiry background worker...")
				return
			case <-ticker.C:
				releasedCount, err := repo.CancelExpiredHolds(ctx)
				if err != nil {
					log.Printf("error expiring seat holds: %v\n", err)
				} else if releasedCount > 0 {
					log.Printf("Cleaned up %d expired seat holds\n", releasedCount)
				}
			}
		}
	}()
}
