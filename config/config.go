package config

// TODO: make it configurable by using viper, or even replacing it entirely by using viper.
//
//nolint:godox,lll
const (
	AppName             = "counter"
	AppShortDescription = "A simple thread-safe REST API counter"
	ApplLongDescription = "A simple thread-safe REST API counter of http requests fired up to the `GET /counter` route which supports restoring the last state of the counter from disk and graceful shutdown."

	ServerHTTPPort                           = 3000
	ServerTimeWindowInSecondsToCountRequests = 60
	ServerFilenameToStoreCounterState        = "counter.json"
)
