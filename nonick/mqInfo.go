package nonick

type MQInfo struct {
	AzCode        string     `json:"az_code"`
	Endpoint      string     `json:"endpoint"`
	MqChannel     string     `json:"mq_channel"`
	MqEndpoint    string     `json:"mq_endpoint"`
	MqPassword    string     `json:"mq_password"`
	MqUserName    string     `json:"mq_user_name"`
	Name          string     `json:"name"`
	NickName      string     `json:"nick_name"`
	SDN           bool       `json:"sdn"`
}

type MqInfoResp struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Data   []MQInfo `json:"data"`
	ResultType  int         `json:"resultType"`
	ElapsedTime int         `json:"elapsedTime"`
	Timestamp   int         `json:"timestamp"`
	Exception   interface{} `json:"exception"`
	TraceId     interface{} `json:"traceId"`
	RequestId   interface{} `json:"requestId"`
}
