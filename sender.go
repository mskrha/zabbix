package zabbix

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const (
	defaultHost = "localhost"
	defaultPort = 10051

	protocol = "ZBXD"
)

type Zabbix struct {
	url string
}

func NewSender(h string, p uint) (*Zabbix, error) {
	if len(h) == 0 {
		return nil, fmt.Errorf("No hostname given")
	}
	if p == 0 {
		return nil, fmt.Errorf("No port given")
	}
	var z Zabbix
	z.url = fmt.Sprintf("%s:%d", h, p)
	return &z, nil
}

func NewDefaultSender() (*Zabbix, error) {
	return NewSender(defaultHost, defaultPort)
}

func (z *Zabbix) Send(r *Request) (ret Response, err error) {
	j, err := json.Marshal(r)
	if err != nil {
		return
	}

	data := []byte(protocol)          // Protocol header
	data = append(data, 0x01)         // Protocol flags
	data = append(data, getLen(j)...) // Length of payload (little endian order)
	data = append(data, j...)         // Payload
	l := len(data)

	c, err := net.Dial("tcp", z.url)
	if err != nil {
		return
	}
	defer c.Close()
	c.SetDeadline(time.Now().Add(10 * time.Second))

	n, err := c.Write(data)
	if err != nil {
		return
	}

	if n != l {
		err = fmt.Errorf("Successfully sent %d bytes of %d requested", n, l)
		return
	}

	buf := make([]byte, 1024)
	n, err = c.Read(buf)
	if err != nil {
		return
	}

	if n == 0 {
		err = fmt.Errorf("No response from Zabbix")
		return
	}

	if n < 13 {
		err = fmt.Errorf("Incomplete response from Zabbix")
		return
	}

	if n == 1024 {
		err = fmt.Errorf("Buffer overflow reading the response from Zabbix")
		return
	}

	ret, err = parseResponse(buf[:n])
	return
}

func getLen(s []byte) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(len(s)))
	return buf
}
