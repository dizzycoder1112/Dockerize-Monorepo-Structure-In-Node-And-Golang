```go
import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"your-project/logger"
	"your-project/internal/config"
)


func main() {
	config.Load()

	// Create logger first (needed by other components)
	zapLogger := logger.NewZapMust(logger.ZapOptions{
		ServiceName: "notify-worker",
		IsPretty:    config.AppConfig.Environment != "production",
		Level:       logger.InfoLevel,
	})



	// Build logger strategies
	strategies := []logger.Logger{zapLogger}

  if config.AppConfig.SlackErrorWebhookURL != "" {
		strategies = append(strategies, logger.NewSlack(logger.SlackOptions{
			WebhookURL:  config.AppConfig.SlackErrorWebhookURL,
			ServiceName: "promotion-worker",
			Environment: config.AppConfig.Environment,
		}))
	}

	appLogger := logger.NewMulti(strategies...)

  // some logic


  appLogger.Info("Shutting down consumer...")

	// Wait for pending async operations (e.g., Slack notifications) to complete
	if err := appLogger.Flush(); err != nil {
		log.Printf("Failed to flush logger: %v", err)
	}
}
```
