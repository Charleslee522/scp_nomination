package ledger

import "log"

import . "github.com/Charleslee522/scp_nomination/src/common"

func init() {
	log.SetPrefix("[Trace] ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

type History struct {
	nodeName          string
	quorumThreshold   int
	blockingThreshold int

	votes    *VotingBox
	accepted *VotingBox
	confirm  *VotingBox

	selfMessageState map[Value]FederatedVotingState
}

func NewHistory(nodeName string, n_validator int, quorumThreshold int) History {
	p := History{nodeName: nodeName, quorumThreshold: quorumThreshold,
		blockingThreshold: n_validator - quorumThreshold + 1}
	p.votes = NewVotingBox()
	p.accepted = NewVotingBox()
	p.confirm = NewVotingBox()
	p.selfMessageState = make(map[Value]FederatedVotingState)
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

		h.votes.AddVotedNode(value, nodeName)
		if h.selfMessageState[value] == NONE {
			h.selfMessageState[value] = VOTES
		}

		if h.votes.Count(value) >= h.quorumThreshold {
			log.Println(value, "in votes exceed quorum threshold",
				h.quorumThreshold, ", so it is moved to accepted")
			h.accepted.AddVotedNode(value, h.nodeName)
			h.selfMessageState[value] = ACCEPTED
		}
	}
}

func (h *History) AppendAccepted(values []Value, nodeName string) {
	for _, value := range values {
		if h.confirm.HasValue(value) {
			continue
		}

		h.accepted.AddVotedNode(value, nodeName)

		if h.accepted.Count(value) >= h.quorumThreshold {
			log.Println(value, "in accepted exceed quorum threshold",
				h.quorumThreshold, ", so it is moved to confirm")
			h.confirm.AddVotedNode(value, h.nodeName)
			h.selfMessageState[value] = CONFIRM
		}
	}
}
