package main

import (
  "net"
  "strconv"
  "fmt"
//  collectd "github.com/paulhammond/gocollectd"
  goopt "github.com/droundy/goopt"
)


var port = goopt.Int([]string{"-p", "--port"}, 8126, "UDP Port to use")

func main() {
  goopt.Description = func() string {
		return "Metric Wrapper for (at first) graphite & elasticsearch."
  }
  goopt.Version = "1.0"
  goopt.Summary = "gostats"
  goopt.Parse(nil)


  addr, _ := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(*port))
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
