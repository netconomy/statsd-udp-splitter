package main

import (
  "net"

  "fmt"
  collectd "github.com/paulhammond/gocollectd"
)

func main() {
  addr, _ := net.ResolveUDPAddr("udp", ":2000")
  sock, _ := net.ListenUDP("udp", addr)

  i := 0
  for {
    i++
    buf := make([]byte, 1024)
    rlen, _, err := sock.ReadFromUDP(buf)
    if err != nil {
      fmt.Println(err)
    }
    fmt.Println(string(buf[0:rlen]))
    fmt.Println(i)
    //go handlePacket(buf, rlen)
  }
}
