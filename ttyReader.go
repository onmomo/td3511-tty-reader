package main

import (
	"bufio"
	"flag"
	"net"
	"time"

	"github.com/ian-kent/go-log/appenders"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/log"
	"github.com/tarm/serial"
)

func main() {
	devicePtr := flag.String("device", "/dev/ttyUSB0", "the smartmeter device")
	hostPtr := flag.String("host", "10.0.1.2:9999", "the host that will receive the data")
	protocolPtr := flag.String("protocol", "udp", "the protocol for the host connection (tcp, udp and IP networks)")
	flag.Parse()
	initLogger()
	read(*devicePtr, *hostPtr, *protocolPtr)
}

func initLogger() {
	logger := log.Logger()
	logger.SetAppender(appenders.RollingFile("smartmeter.log", true))
	appender := logger.Appender()
	appender.SetLayout(layout.Pattern("%d %p - %m%n"))
}

func read(device string, host string, protocol string) {
	log.Info("Opening smartmeter device '%s' ...", device)
	readTimeout := time.Minute * 5
	c := &serial.Config{Name: device, Baud: 300, Size: 7, Parity: 'E', ReadTimeout: readTimeout}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Waiting a maximum of %d minutes for smartmeter data.", readTimeout)

	log.Info("Configure smartmeter Td3511 for read out ...")
	_, err = s.Write([]byte("1:0:9a7:0:3:1c:7f:15:4:5:1:0:11:13:1a:0:12:f:17:16:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.Write([]byte("\x2F\x3F\x21\x0D\x0A"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.Write([]byte("\x06\x30\x30\x30\x0D\x0A"))
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Waiting for data ...")
	reader := bufio.NewReader(s)
	readData, err := reader.ReadString('\x21')
	if err != nil {
		log.Error("Couldn't read any data.")
		log.Fatal(err)
	}

	matchedData := matchData(readData)
	matchedDataSize := len(matchedData)
	log.Info("Received %d data records.", matchedDataSize)
	if matchedDataSize > 0 {
		log.Info("Opening %s connection to %s ...", protocol, host)
		conn, err := net.Dial(protocol, host)
		if err != nil {
			log.Fatal(err)
		}

		for key, value := range matchedData {
			if value["omis"] == "1.7.0" {
				data := value["data"]
				conn.Write([]byte("1.7.0:" + data))
				log.Info("1.7.0/%d: Aktueller Verbrauch (%s): %s", key, value["unit"], data)
			} else if value["omis"] == "1.8.0" {
				data := value["data"]
				conn.Write([]byte("1.8.0:" + data))
				log.Info("1.8.0/%d: Verbrauch Gesamt (%s): %s", key, value["unit"], data)
			} else if value["omis"] == "1.8.1" {
				data := value["data"]
				conn.Write([]byte("1.8.1:" + data))
				log.Info("1.8.1/%d: Verbrauch Tarif 1 (%s): %s", key, value["unit"], data)
			} else if value["omis"] == "1.8.2" {
				data := value["data"]
				conn.Write([]byte("1.8.2:" + data))
				log.Info("1.8.2/%d: Verbrauch Tarif 2 (%s): %s", key, value["unit"], data)
			} else if value["omis"] == "2.7.0" {
				data := value["data"]
				conn.Write([]byte("2.7.0:" + data))
				log.Info("2.7.0/%d: Aktuelle Lieferung (%s): %s", key, value["unit"], data)
			} else if value["omis"] == "2.8.0" {
				data := value["data"]
				conn.Write([]byte("2.8.0:" + data))
				log.Info("2.8.0/%d: Lieferung Gesamt (%s): %s", key, value["unit"], data)
			} else if value["omis"] == "2.8.1" {
				data := value["data"]
				conn.Write([]byte("2.8.1:" + data))
				log.Info("2.8.1/%d: Lieferung Tarif 1 (%s): %s", key, value["unit"], data)
			} else if value["omis"] == "2.8.2" {
				data := value["data"]
				conn.Write([]byte("2.8.2:" + data))
				log.Info("2.8.2/%d: Lieferung Tarif 2 (%s): %s", key, value["unit"], data)
			}
		}
		log.Info("Sucessfully processed all %d data records, closing connection ...", matchedDataSize)
		conn.Close()
	} else {
		log.Warn("Couldn't match any data records.")
	}

	log.Info("Closing TD3511 smartmeter, bye bye.")
}
