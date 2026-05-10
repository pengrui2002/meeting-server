export interface Peer {
  id: string;
  name: string;
  producers: Map<string, any>;
  consumers: Map<string, any>;
}

export class PeerManager {
  private peers: Map<string, Peer> = new Map();

  addPeer(peerId: string, name: string): Peer {
    const peer: Peer = {
      id: peerId,
      name,
      producers: new Map(),
      consumers: new Map(),
    };
    this.peers.set(peerId, peer);
    return peer;
  }

  getPeer(peerId: string): Peer | undefined {
    return this.peers.get(peerId);
  }

  removePeer(peerId: string): void {
    this.peers.delete(peerId);
  }

  getPeers(): Peer[] {
    return Array.from(this.peers.values());
  }

  getPeerCount(): number {
    return this.peers.size;
  }
}