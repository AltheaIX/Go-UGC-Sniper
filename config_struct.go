package main

type ConfigOptions struct {
	AlwaysCheckProxy bool
	Threads          int
}

type ProxyOptions struct {
	Enable   bool
	Username string
	Password string
}

type Config struct {
	FreeWebhook string
	PaidWebhook string
	Cookie      string
	OffsaleId   []int
	LastId      int
	Proxy       ProxyOptions
	Options     ConfigOptions
}
