[![Go Report Card](https://goreportcard.com/badge/github.com/mskrha/zabbix)](https://goreportcard.com/report/github.com/mskrha/zabbix)

## zabbix

### Description
Go library for sending data to the Zabbix using the [Zabbix sender protocol](https://www.zabbix.com/documentation/6.0/en/manual/appendix/items/trapper).

### Installation
`go get github.com/mskrha/zabbix`

### Example usage
```go
package main

import (
	"fmt"

	"github.com/mskrha/zabbix"
)

func main() {
	s, err := zabbix.NewSender("zabbix-server.example.net", 10051)
	if err != nil {
		fmt.Println(err)
		return
	}

	r := zabbix.NewRequest()
	r.AppendData("server1.example.net", "custom.key.1", "value 1")
	r.AppendData("server1.example.net", "custom.key.2", "value 2")

	x, err := s.Send(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", x)
}
```
