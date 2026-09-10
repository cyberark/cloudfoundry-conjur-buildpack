package main

// Version is the buildpack version injected at build time via:
//
//	go build -ldflags "-X main.Version=<version>"
//
// If not set at build time, it defaults to "dev".
var Version = "dev"
