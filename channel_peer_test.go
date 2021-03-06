package graft

import (
	"github.com/benmills/quiz"
	"testing"
)

func TestChannelPeerCanBeAddedToAServersListOfPeers(t *testing.T) {
	server := New("id")
	peer := NewChannelPeer(server)
	server2 := New("id")
	server2.Peers = []Peer{peer}
}

func TestChannelPeerRespondsToVoteMessages(t *testing.T) {
	test := quiz.Test(t)

	server := New("id")
	peer := NewChannelPeer(server)
	peer.Start()
	defer peer.ShutDown()
	requestVote := RequestVoteMessage{
		Term:         1,
		CandidateId:  "foo",
		LastLogIndex: 0,
		LastLogTerm:  0,
	}
	response, err := peer.ReceiveRequestVote(requestVote)

	test.Expect(err).ToEqual(nil)
	test.Expect(response.VoteGranted).ToBeTrue()
	test.Expect(server.Term).ToEqual(1)
}

func TestChannelPeerRespondsToAppendEntriesMessages(t *testing.T) {
	test := quiz.Test(t)

	server := New("id")
	peer := NewChannelPeer(server)
	peer.Start()
	defer peer.ShutDown()
	message := AppendEntriesMessage{
		Term:         2,
		LeaderId:     "leader_id",
		PrevLogIndex: 2,
		Entries:      []LogEntry{},
		CommitIndex:  0,
	}

	peer.ReceiveAppendEntries(message)

	test.Expect(server.Term).ToEqual(2)
}

func TestPartitionedPeerRespondsWithError(t *testing.T) {
	test := quiz.Test(t)

	server := New("id")
	peer := NewChannelPeer(server)
	peer.Start()
	defer peer.ShutDown()
	requestVote := RequestVoteMessage{
		Term:         1,
		CandidateId:  "foo",
		LastLogIndex: 0,
		LastLogTerm:  0,
	}
	peer.Partition()
	_, err := peer.ReceiveRequestVote(requestVote)

	test.Expect(err != nil).ToBeTrue()
}

func TestPartitionCanBeHealed(t *testing.T) {
	test := quiz.Test(t)

	server := New("id")
	peer := NewChannelPeer(server)
	peer.Start()
	defer peer.ShutDown()
	requestVote := RequestVoteMessage{
		Term:         1,
		CandidateId:  "foo",
		LastLogIndex: 0,
		LastLogTerm:  0,
	}
	peer.Partition()
	_, err := peer.ReceiveRequestVote(requestVote)

	test.Expect(err != nil).ToBeTrue()

	peer.HealPartition()

	_, err2 := peer.ReceiveRequestVote(requestVote)

	test.Expect(err2).ToEqual(nil)

}
