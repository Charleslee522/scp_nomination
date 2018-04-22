package ledger

import (
	"fmt"
	"testing"
)

type Node struct {
	name     string
	priority int
	// message  []string
}

func (n Node) GetName() string {
	return n.name
}

type Message struct {
	data string
}

type Ledger struct {
	node       Node
	validators []Node
}

func (l Ledger) receiveMessage(msg Message) {
	fmt.Println(l.node.GetName())
	fmt.Println(msg)
}

// 노드 세 개 정의(n1, n2, n3)
// 이 테스트의 주인공은 n1
// 서로를 모두 validator 로 지정
// 리더 선출(n1) - TestLedgerSelfLeader
// 	다른 노드의 메시지를 받으면 저장만 하고, echoing 하지 않음
// 리더 선출(n2) - TestLedgerOtherLeader
// 	n2의 메시지를 받으면 저장 후에 echoing 함
// 	n3의 메시지를 받으면 저장만 하고, echoing 하지 않음

func TestLedger(t *testing.T) {

	var node1 Node = Node{"n1", 1}
	var node2 Node = Node{"n2", 0}
	var node3 Node = Node{"n3", 0}

	nodes := []Node{node1, node2, node3}

	var ledger1 Ledger = Ledger{Node{"n1", 1}, nodes}
	var ledger2 Ledger = Ledger{Node{"n2", 0}, nodes}
	var ledger3 Ledger = Ledger{Node{"n3", 0}, nodes}

	fmt.Println(ledger1.validators)

	var msg1 Message = Message{"I'm Charles!"}

	ledger1.receiveMessage(msg1)
	ledger2.receiveMessage(msg1)
	ledger3.receiveMessage(msg1)

	// bc1 = blockchain_factory(
	//     node_name_1,
	//     'http://localhost:5001',
	//     100,
	//     [node_name_2, node_name_3]
	// )

	// bc2 = blockchain_factory(
	//     node_name_2,
	//     'http://localhost:5002',
	//     100,
	//     [node_name_1, node_name_3]
	// )

	// bc3 = blockchain_factory(
	//     node_name_3,
	//     'http://localhost:5003',
	//     100,
	//     [node_name_1, node_name_2]
	// )

	// bc1.consensus.add_to_validator_connected(bc2.node)
	// bc1.consensus.add_to_validator_connected(bc3.node)
	// bc1.consensus.init()

	// message = Message.new('message')
	// ballot_init_1 = Ballot.new(node_name_1, message, IsaacState.INIT, BallotVotingResult.agree)
	// ballot_id = ballot_init_1.ballot_id
	// ballot_init_2 = Ballot(ballot_id, node_name_2, message, IsaacState.INIT, BallotVotingResult.agree,
	//                        ballot_init_1.timestamp)
	// ballot_init_3 = Ballot(ballot_id, node_name_3, message, IsaacState.INIT, BallotVotingResult.agree,
	//                        ballot_init_1.timestamp)

	// bc1.receive_ballot(ballot_init_1)
	// bc1.receive_ballot(ballot_init_2)
	// bc1.receive_ballot(ballot_init_3)

	// assert bc1.consensus.slot.get_ballot_state(ballot_init_1) == IsaacState.SIGN
}
