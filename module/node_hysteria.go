package module

type NodeHysteria struct {
	HysteriaProtocol *string `ddb:"hysteria_protocol"`
	HysteriaUpMbps   *int    `ddb:"hysteria_up_mbps"`
	HysteriaDownMbps *int    `ddb:"hysteria_down_mbps"`
}
