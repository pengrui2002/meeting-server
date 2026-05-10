"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.PeerManager = void 0;
class PeerManager {
    constructor() {
        this.peers = new Map();
    }
    addPeer(peerId, name) {
        const peer = {
            id: peerId,
            name,
            producers: new Map(),
            consumers: new Map(),
        };
        this.peers.set(peerId, peer);
        return peer;
    }
    getPeer(peerId) {
        return this.peers.get(peerId);
    }
    removePeer(peerId) {
        this.peers.delete(peerId);
    }
    getPeers() {
        return Array.from(this.peers.values());
    }
    getPeerCount() {
        return this.peers.size;
    }
}
exports.PeerManager = PeerManager;
//# sourceMappingURL=peer.js.map