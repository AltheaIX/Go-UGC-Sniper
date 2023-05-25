package main

type Database struct {
	Version Version
	Trial   Trial
}

type Trial struct {
	Status string
}

type Version struct {
	Version string
}
