package zabbix

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
)

type response struct {
	Response string `json:"response"`
	Info     string `json:"info"`
}

type Response struct {
	Processed uint64
	Failed    uint64
	Total     uint64
	Time      time.Duration
}

func parseResponse(in []byte) (ret Response, err error) {
	if !bytes.Equal(in[0:4], []byte(protocol)) {
		err = fmt.Errorf("Not a Zabbix response")
		return
	}

	if in[4] != 0x01 {
		err = fmt.Errorf("Unsupported Zabbix protocol")
		return
	}

	l := binary.LittleEndian.Uint64(in[5:13])
	if (l + 8 + 4 + 1) != uint64(len(in)) {
		err = fmt.Errorf("Response length mismatch")
		return
	}

	var r response
	err = json.Unmarshal(in[13:], &r)
	if err != nil {
		return
	}

	if r.Response != "success" {
		err = fmt.Errorf("Request was not successful")
		return
	}

	var n int
	var t float64
	n, err = fmt.Sscanf(r.Info, "processed: %d; failed: %d; total: %d; seconds spent: %f", &ret.Processed, &ret.Failed, &ret.Total, &t)
	if err != nil {
		return
	}
	if n != 4 {
		err = fmt.Errorf("Failed to parse the response")
		return
	}
	ret.Time, err = time.ParseDuration(fmt.Sprintf("%fs", t))

	return
}
