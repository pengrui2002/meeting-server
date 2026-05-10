"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TransportManager = void 0;
class TransportManager {
    constructor() {
        this.transports = new Map();
    }
    async createTransport(router, transportId, type, peerId) {
        const transport = await router.createWebRtcTransport({
            listenIps: [
                {
                    ip: '0.0.0.0',
                    announcedIp: process.env.MEDIASOUP_ANNOUNCED_IP || '127.0.0.1',
                },
            ],
            initialAvailableOutgoingBitrate: 1000000,
        });
        this.transports.set(transportId, { id: transportId, type, peerId, transport });
        return transport;
    }
    getTransport(transportId) {
        return this.transports.get(transportId);
    }
    closeTransport(transportId) {
        const transport = this.transports.get(transportId);
        if (transport) {
            transport.transport.close();
            this.transports.delete(transportId);
        }
    }
}
exports.TransportManager = TransportManager;
//# sourceMappingURL=transport.js.map