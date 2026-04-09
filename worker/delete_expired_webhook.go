package worker

import (
	"context"
	"hookfy/models"
	"log"
	"time"

	"gorm.io/gorm"
)

const cleanupInterval = 10 * time.Minute

func StartDeleteExpiredWorker(db *gorm.DB, ctx context.Context) {
	ticker := time.NewTicker(cleanupInterval)

	go func() {
		defer ticker.Stop() // Cleanup when goroutine exits

		for {
			select {
			case <-ticker.C:
				runCleanup(db)
			case <-ctx.Done():
				log.Println("[worker] shutting down cleanup worker")
				return
			}
		}
	}()
}

func runCleanup(db *gorm.DB) {
	result := db.Unscoped().
		Where("expires_at < ?", time.Now()).
		Delete(&models.Webhook{})

	if result.Error != nil {
		log.Printf("[worker] erro ao apagar webhooks expirados: %v", result.Error)
		return
	}
	if result.RowsAffected > 0 {
		log.Printf("[worker] %d webhook(s) expirado(s) apagado(s)", result.RowsAffected)
	}
}
