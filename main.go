package main

import (
	"bufio"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"os"
)
func main() {
	options := serial.OpenOptions{
		PortName:        "/dev/ttyUSB0",
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}
	rx := make(chan []byte, 100)
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	go RXrun(port, rx)
	comd:=make(chan string,10)
	f,_:=os.Create("log.log")
	defer f.Close()
	go cmd(comd)
	for {
		select {
		case rx:=<-rx:
			f.Write(rx)
		case c:=<-comd:
			port.Write([]byte(c))
		}
	}
}

func RXrun(port io.ReadWriteCloser, bytes chan []byte) {
	dr:=make([]byte,1024)
	for{
		n,_:=port.Read(dr)
		ds:=dr[:n]
		bytes<-ds
	}
}
func cmd(s chan string)  {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		s<-text
	}
}