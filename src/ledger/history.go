package ledger

import "log"

import . "github.com/Charleslee522/scp_nomination/src/common"

func init() {
	log.SetPrefix("[Trace] ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

type History struct {
	quorumThreshold uint16
	votes           map[Value]uint16
	accepted        map[Value]uint16
	confirm         map[Value]uint16
}

func NewHistory(n_validator uint16) History {
	p := History{quorumThreshold: n_validator}
	p.votes = make(map[Value]uint16)
	p.accepted = make(map[Value]uint16)
	p.confirm = make(map[Value]uint16)
	return p
}

func (h *History) AppendMessage(msg SCPNomination) {
	h.AppendVotes(msg.Votes)
	h.AppendAccepted(msg.Accepted)
}

func (h *History) AppendVotes(values []Value) {
	for _, value := range values {
		if h.accepted[value] > 0 || h.confirm[value] > 0 {
			continue
		}
		if h.votes[value] == 0 {
			h.votes[value] = 1
		} else {
			h.votes[value]++
		}

		if h.votes[value] >= h.quorumThreshold {
			log.Println(value, "in votes exceed quorum threshold",
				h.quorumThreshold, ", so it is moved to accepted")
			delete(h.votes, value)
			h.accepted[value] = 1
		}
	}
}

func (h *History) AppendAccepted(values []Value) {
	for _, value := range values {
		if h.confirm[value] > 0 {
			continue
		}
		if h.accepted[value] == 0 {
			h.accepted[value] = 1
		} else {
			h.accepted[value]++
		}

		if h.accepted[value] >= h.quorumThreshold {
			log.Println(value, "in accepted exceed quorum threshold",
				h.quorumThreshold, ", so it is moved to confirm")
			delete(h.accepted, value)
			h.confirm[value] = 1
		}
	}
}
