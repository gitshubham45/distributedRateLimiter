package domain

type EndpointConfig struct {
	Key           string
	Method        string
	Path          string
	TargetService string
	StrategyName  string
	Limit         int64
	Interval      int64
}
