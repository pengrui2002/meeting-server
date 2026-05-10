import { types } from 'mediasoup';

type Router = types.Router;
type WebRtcTransport = types.WebRtcTransport;

interface Transport {
  id: string;
  type: 'producer' | 'consumer';
  peerId: string;
  transport: WebRtcTransport;
}

export class TransportManager {
  private transports: Map<string, Transport> = new Map();

  async createTransport(
    router: Router,
    transportId: string,
    type: 'producer' | 'consumer',
    peerId: string
  ): Promise<WebRtcTransport> {
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

  getTransport(transportId: string): Transport | undefined {
    return this.transports.get(transportId);
  }

  closeTransport(transportId: string): void {
    const transport = this.transports.get(transportId);
    if (transport) {
      transport.transport.close();
      this.transports.delete(transportId);
    }
  }
}