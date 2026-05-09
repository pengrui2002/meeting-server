package media

import (
	"sync"
)

// Room represents a video conferencing room managed by the SFU
type Room struct {
	ID    string
	Peers map[string]*Peer
	mu    sync.RWMutex
}

// Peer represents a participant in a room
type Peer struct {
	ID          string
	DisplayName string
}

// NewRoom creates a new room with the given ID
func NewRoom(id string) *Room {
	return &Room{
		ID:    id,
		Peers: make(map[string]*Peer),
	}
}

// AddPeer adds a peer to the room
func (r *Room) AddPeer(peer *Peer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Peers[peer.ID] = peer
}

// RemovePeer removes a peer from the room
func (r *Room) RemovePeer(peerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Peers, peerID)
}

// GetPeers returns all peers in the room
func (r *Room) GetPeers() []*Peer {
	r.mu.RLock()
	defer r.mu.RUnlock()
	peers := make([]*Peer, 0, len(r.Peers))
	for _, p := range r.Peers {
		peers = append(peers, p)
	}
	return peers
}
