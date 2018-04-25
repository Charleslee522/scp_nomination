package ledger

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func GetPriority(round int, nodeName string) string {
	roundStr := strconv.Itoa(round)
	h := sha256.New()
	h.Write([]byte(roundStr + nodeName))
	return hex.EncodeToString(h.Sum(nil))
}
