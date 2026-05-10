import { types } from 'mediasoup';
type Router = types.Router;
type WebRtcTransport = types.WebRtcTransport;
interface Transport {
    id: string;
    type: 'producer' | 'consumer';
    peerId: string;
    transport: WebRtcTransport;
}
export declare class TransportManager {
    private transports;
    createTransport(router: Router, transportId: string, type: 'producer' | 'consumer', peerId: string): Promise<WebRtcTransport>;
    getTransport(transportId: string): Transport | undefined;
    closeTransport(transportId: string): void;
}
export {};
//# sourceMappingURL=transport.d.ts.map