package td3511reader

import (
	"bufio"
	"github.com/tarm/serial"
	"log"
)

func main() {
	read("/dev/ttyUSB0")
}

func read(device string) {
	log.Printf("Opening device '%s'...", device)
	c := &serial.Config{Name: device, Baud: 300, Size: 7, Parity:'E'}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.Write([]byte("1:0:9a7:0:3:1c:7f:15:4:5:1:0:11:13:1a:0:12:f:17:16:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.Write([]byte("\x2F\x3F\x21\x0D\x0A"))
	if err != nil {
		log.Fatal(err)
	}

	_, err2 := s.Write([]byte("\x06\x30\x30\x30\x0D\x0A"))
	if err2 != nil {
		log.Fatal(err)
	}

	log.Print("Reading data...")
	reader := bufio.NewReader(s)
	reply, err := reader.ReadBytes('\x21')
	if err != nil {
		panic(err)
	}

	data := string(reply)
	//log.Print(data)

	matchedData := matchData(data)


	log.Print("start------------------------")
	for key, value := range matchedData {
		if value["omis"] == "1.8.0" {
			log.Printf("1.8.0/%d: Verbrauch Gesamt (%s): %s", key, value["unit"], value["data"])
		} else if value["omis"] == "1.8.1" {
			log.Printf("1.8.1/%d: Verbrauch Tarif 1 (%s): %s", key, value["unit"], value["data"])
		} else if value["omis"] == "1.8.2" {
			log.Printf("1.8.2/%d: Verbrauch Tarif 2 (%s): %s", key, value["unit"], value["data"])
		} else if value["omis"] == "2.8.0" {
			log.Printf("2.8.0/%d: Lieferung Gesamt (%s): %s", key, value["unit"], value["data"])
		} else if value["omis"] == "2.8.1" {
			log.Printf("2.8.1/%d: Lieferung Tarif 1 (%s): %s", key, value["unit"], value["data"])
		} else if value["omis"] == "2.8.2" {
			log.Printf("2.8.2/%d: Lieferung Tarif 2 (%s): %s", key, value["unit"], value["data"])
		}
	}
	log.Print("end------------------------")
}

