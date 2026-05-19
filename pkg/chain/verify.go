package chain

func IsValidChain(hash string, custodyCount int) bool {
	return hash != "" && custodyCount >= 1
}
