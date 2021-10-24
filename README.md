# free-Epic-games-Slack-bot
Bot that periodically sends information about the currently free games in Epic Games' store to a Slack channel via Webhook URL.

## Configure
You need to paste the Webhook URL for the desired Slack channel into the slackWebhookURL constant.

By default, the bot sends the information every 24 hours. You can customize that period through the periodSize and period Unit (time.Hour, time.Minute, etc.) constants.

```
const (
	// URLs
	rawJsonURL      = "https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions"
	epicStoreURL    = "https://www.epicgames.com/store/en-US/p/"
	slackWebhookURL = "" // Paste your Slack Webhook URL here
	periodSize      = 24
	periodUnit      = time.Hour
)
```
