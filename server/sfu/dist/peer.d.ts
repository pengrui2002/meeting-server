export interface Peer {
    id: string;
    name: string;
    producers: Map<string, any>;
    consumers: Map<string, any>;
}
export declare class PeerManager {
    private peers;
    addPeer(peerId: string, name: string): Peer;
    getPeer(peerId: string): Peer | undefined;
    removePeer(peerId: string): void;
    getPeers(): Peer[];
    getPeerCount(): number;
}
//# sourceMappingURL=peer.d.ts.map