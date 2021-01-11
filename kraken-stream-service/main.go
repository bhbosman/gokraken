package main

import (
	"flag"
	"github.com/bhbosman/gokraken/kraken-stream-service/internal"
	"github.com/kardianos/service"
	"log"
)

func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	program := &internal.ServiceInterfaceImpl{}

	svcConfig := &service.Config{
		Name:        "bhbosman.Kraken.Stream.Service",
		DisplayName: "bhbosman.Kraken.Stream.Service",
		Description: "bhbosman.Kraken.Stream.Service",
		Option:      make(service.KeyValue),
	}
	svcConfig.Option["UserService"] = true
	svcConfig.Option["KeepAlive"] = true
	svcConfig.Option["RunAtLoad"] = true
	svcConfig.Option["SessionCreate"] = true
	svcConfig.Option["LaunchdConfig"] = ""

	s, err := service.New(program, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		return
	}
}
