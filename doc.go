/*
A Go library for communicating with LightSwarm LED's over an io.Writer interface.

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
*/
package lightswarm
