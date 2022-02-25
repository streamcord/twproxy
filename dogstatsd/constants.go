package dogstatsd

// GlobalRate is the sample rate for DogStatsD.
// https://docs.datadoghq.com/metrics/dogstatsd_metrics_submission/?code-lang=go#metric-submission-options
const GlobalRate = 1.0

// DogStatsD metric names
const (
	MetricRequest                          = "web.request"
	MetricTwitchEventSubSubscriptionsUsed  = "twitch.eventsub.subscriptions.used"
	MetricTwitchEventSubSubscriptionsTotal = "twitch.eventsub.subscriptions.total"
	MetricTwitchRequest                    = "twitch.request"
	MetricTwitchRatelimitRemaining         = "twitch.ratelimit.remaining"
)

// Twitch API routes
const (
	RouteCreateEventSubSubscription = "POST /helix/eventsub/subscriptions"
	RouteDeleteEventSubSubscription = "DELETE /helix/eventsub/subscriptions"
	RouteGetEventSubSubscriptions   = "GET /helix/eventsub/subscriptions"
	RouteGetClips                   = "/helix/clips"
	RouteGetGames                   = "/helix/games"
	RouteGetStreams                 = "/helix/streams"
	RouteGetUsers                   = "/helix/users"
	RouteGetUsersFollows            = "/helix/users/follows"
	RouteGetOauthToken              = "/oauth/token"
)
