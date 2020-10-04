package main

import "github.com/bhbosman/gokraken/internal"

func main() {
	fxApp, _ := internal.CreateFxApp()
	if fxApp.Err() != nil {
		return
	}
	fxApp.Run()
}
