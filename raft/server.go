package raft

type Server struct {
	Addr string
	Cm   *ConsensusModule
}
