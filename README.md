# lightswarm

[![Build Status](https://travis-ci.org/thisissoon/lightswarm.svg?branch=master)](https://travis-ci.org/thisissoon/lightswarm)

A Go library for communicating with LightSwarm LED's.

## Usage

This example uses http://github.com/tarm/serial to open a serial connection. There
are many other serial libraries availible. All this package requires is that the
library implements the `io.Writer` interface.

``` go
package main

import (
	"log"
	"time"

	"github.com/tarm/serial"
	"github.com/thisissoon/lightswarm"
)

func main() {
	w, err := serial.OpenPort(&serial.Config{
		Name: "/dev/tty.usbserial-DA00YSEB",
		Baud: 38400,
	})
	if err != nil {
		log.Fatal(err)
	}
	led := lightswarm.New(690, w)
	if _, err := led.On(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5)
	if _, err := led.Off(); err != nil {
		log.Fatal(err)
	}
}
```
