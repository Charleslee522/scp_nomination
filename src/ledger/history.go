package ledger

import "log"

import . "github.com/Charleslee522/scp_nomination/src/common"

func init() {
	log.SetPrefix("[Trace] ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

type History struct {
	nodeName          string
	leaderName        string
	quorumThreshold   int
	blockingThreshold int

	votes    *VotingBox
	accepted *VotingBox
	confirm  *VotingBox

	selfMessageState map[Value]FederatedVotingState
	confirmValues    []Value
}

func NewHistory(nName string, lName string, n_validator int, qTh int) History {
	p := History{nodeName: nName, leaderName: lName, quorumThreshold: qTh,
		blockingThreshold: n_validator - qTh + 1}
	p.votes = NewVotingBox()
	p.accepted = NewVotingBox()
	p.confirm = NewVotingBox()
	p.selfMessageState = make(map[Value]FederatedVotingState)
	p.confirmValues = []Value{}
	return p
}

func (h *History) AppendMessage(msg SCPNomination) {
	h.AppendVotes(msg.Votes, msg.NodeName)
	h.AppendAccepted(msg.Accepted, msg.NodeName)
}

func (h *History) AppendVotes(values []Value, nodeName string) {
	for _, value := range values {
		if h.accepted.HasValue(value) ||
			h.confirm.HasValue(value) {
			continue
		}

		h.votes.Add(value, nodeName)
		if h.selfMessageState[value] == NONE {
			h.selfMessageState[value] = VOTES
		}

		if h.votes.Count(value) >= h.quorumThreshold {
			log.Println(value, "in votes exceed quorum threshold",
				h.quorumThreshold, ", so it is moved to accepted")
			h.accepted.Add(value, h.nodeName)
			h.selfMessageState[value] = ACCEPTED
		}
	}
}

func (h *History) AppendAccepted(values []Value, nodeName string) {
	for _, value := range values {
		if h.confirm.HasValue(value) {
			continue
		}

		h.accepted.Add(value, nodeName)

		if h.accepted.Count(value) >= h.blockingThreshold {
			log.Println(value, "in accepted exceed blocking threshold",
				h.quorumThreshold, ", so it is moved to accept")
			h.accepted.Add(value, h.nodeName)
			h.selfMessageState[value] = ACCEPTED
		}

		if h.accepted.Count(value) >= h.quorumThreshold {
			log.Println(value, "in accepted exceed quorum threshold",
				h.quorumThreshold, ", so it is moved to confirm")
			h.confirm.Add(value, h.nodeName)
			h.selfMessageState[value] = CONFIRM
			h.confirmValues = append(h.confirmValues, value)
		}
	}
}
