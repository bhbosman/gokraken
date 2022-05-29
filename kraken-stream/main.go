package main

import "github.com/bhbosman/gokraken/internal"

func main() {
	fxApp := internal.CreateFxApp()
	fxApp.RunTerminalApp()
}
