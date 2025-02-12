package zabbix

type Data struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Request struct {
	Request string `json:"request"`
	Data    []Data `json:"data"`
}

func NewRequest() *Request {
	var r Request
	r.Request = "sender data"
	return &r
}

func (r *Request) AppendData(h, k, v string) {
	r.Data = append(r.Data, Data{Host: h, Key: k, Value: v})
}
