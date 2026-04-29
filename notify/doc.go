// Package notify delivers post-sync webhook notifications for vaultpipe.
//
// After a sync, rotate, or diff operation completes, callers can construct
// an Event and dispatch it via a Notifier to any HTTP webhook endpoint
// (e.g. Slack incoming webhooks, PagerDuty, or a custom alerting service).
//
// Basic usage:
//
//	n, err := notify.New("https://hooks.slack.com/services/XXX/YYY/ZZZ")
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = n.Send(notify.Event{
//		Operation: "sync",
//		Role:      "backend",
//		Keys:      []string{"DB_HOST", "API_KEY"},
//	})
//
// The Timestamp field is automatically populated with the current UTC time
// if left at its zero value.
package notify
