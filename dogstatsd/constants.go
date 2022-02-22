package dogstatsd

// GlobalRate is the sample rate for DogStatsD.
// https://docs.datadoghq.com/metrics/dogstatsd_metrics_submission/?code-lang=go#metric-submission-options
const GlobalRate = 1.0

// DogStatsD metric names
const (
	MetricRequest                  = "web.request"
	MetricTwitchRequest            = "twitch.request"
	MetricTwitchRatelimitRemaining = "twitch.ratelimit.remaining"
)
