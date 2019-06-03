package config

type StatsdConfig struct {
	Prefix   *string
	Addr     *string
	StuckDur *string
	Mocking  *string
	Disabled *string
}
