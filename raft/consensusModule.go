package raft

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	Follower  string = "Follower"
	Candidate        = "Candidate"
	Leader           = "Leader"
)

type ConsensusModule struct {
	mu                 sync.Mutex
	currentTerm        int
	state              string
	electionResetEvent time.Time
}

func (cm *ConsensusModule) runElectionTimer() {
	timeoutDuration := cm.getElectionTimeout()
	cm.dlog("Election timeout set to " + string(timeoutDuration))
	// ...
	cm.mu.Lock()
	termStarted := cm.currentTerm
	cm.dlog(fmt.Sprintf("Election timer started (%v), term=%d", timeoutDuration, termStarted))
	cm.mu.Unlock()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		<-ticker.C
		cm.mu.Lock()
		// TODO: Why cannot we simply check cm.state == leader?
		if cm.state != Candidate && cm.state != Follower {
			cm.dlog(fmt.Sprintf("in election state=%s, bailing out", cm.state))
			cm.mu.Unlock()
			return
		}
		if cm.currentTerm != termStarted {
			cm.dlog(fmt.Sprintf("in elcetion timer term changed from %d to %d, bailing out", termStarted, cm.currentTerm))
			cm.mu.Unlock()
			return
		}

		// start an election if we haven't hearf from a leader or haven't voted for
		// someone for the duration of the timeout
		if elapsed := time.Since(cm.electionResetEvent); elapsed > timeoutDuration {
			cm.mu.Unlock()
			cm.startElection()
			return
		}
	}
}

func (cm *ConsensusModule) startElection() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.dlog(fmt.Sprintf("Starting election"))
}

func (cm *ConsensusModule) dlog(message string) {
	// ...
	log.Println(message)
}

func (cm *ConsensusModule) getElectionTimeout() time.Duration {
	// ...
	//generate a random timeout
	cm.electionResetEvent = time.Now()
	return time.Duration(rand.Intn(150)+150) * time.Millisecond
}
