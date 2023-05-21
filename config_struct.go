package main

type Config struct {
	FreeWebhook string
	PaidWebhook string
	Cookie      string
	LastId      int
	OffsaleId   []int
	Options     ConfigOptions
}

type ConfigOptions struct {
	Threads          int
	AlwaysCheckProxy bool
}
