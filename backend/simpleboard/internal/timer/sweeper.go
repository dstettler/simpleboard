package timer

import (
	"log"
	"time"

	"simpleboard/internal/chess"
	"simpleboard/internal/repository"

	"gorm.io/gorm"
)

// StartSweeper launches a background goroutine that periodically scans
// in-progress games and ends any whose active side has flag-fallen.
// This guarantees a game can't sit forever just because no one is polling.
func StartSweeper(db *gorm.DB, interval time.Duration) {
	if interval <= 0 {
		interval = 30 * time.Second
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		log.Printf("Timer sweeper running every %s", interval)

		for range ticker.C {
			sweepOnce(db, time.Now())
		}
	}()
}

// sweepOnce does a single pass over in-progress games and ends timed-out ones.
// Pulled out so it's easy to invoke directly from a test.
func sweepOnce(db *gorm.DB, now time.Time) {
	var games []repository.Game
	if err := db.Where("status = ?", chess.InProgress.String()).Find(&games).Error; err != nil {
		log.Printf("timer sweeper: query failed: %v", err)
		return
	}

	for i := range games {
		g := &games[i]
		_, _, timedOut, loser := LiveRemaining(g, now)
		if !timedOut {
			continue
		}

		g.Status = FlagFallStatus(loser).String()
		if loser == "w" {
			g.WhiteRemainingMs = 0
		} else {
			g.BlackRemainingMs = 0
		}

		if err := db.Save(g).Error; err != nil {
			log.Printf("timer sweeper: save game %d failed: %v", g.ID, err)
			continue
		}
		log.Printf("timer sweeper: game %d ended on time (%s flagged)", g.ID, loser)
	}
}
