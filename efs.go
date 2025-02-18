package main

import "embed"

//go:embed assets/* data/*
var staticFiles embed.FS
