package dfnet

func FuzzNetAddrJSON(data []byte) int {
	var addr NetAddr
	if err := addr.UnmarshalJSON(data); err != nil {
		return 0
	}
	return 1
}
