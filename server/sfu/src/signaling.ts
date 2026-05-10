import { types } from 'mediasoup';

type Router = types.Router;
import { RouterManager } from './router';
import { TransportManager } from './transport';
import { PeerManager } from './peer';

export function createSignalingRouter(
  routerManager: RouterManager,
  transportManager: TransportManager,
  peerManager: PeerManager
): any {
  const express = require('express');
  const router = express.Router();
  const { v4: uuidv4 } = require('uuid');

  // 创建/获取 Room
  router.post('/api/rooms', async (req: any, res: any) => {
    try {
      const { meetingId } = req.body;
      const roomId = meetingId || uuidv4();
      const room = await routerManager.createRoom(roomId);

      res.json({
        room_id: room.id,
        rtp_capabilities: room.router.rtpCapabilities,
      });
    } catch (error: any) {
      console.error('Error creating room:', error);
      res.status(500).json({ error: error.message });
    }
  });

  // 创建 Transport
  router.post('/api/rooms/:roomId/transports', async (req: any, res: any) => {
    try {
      const { roomId } = req.params;
      const { type, peerId } = req.body;

      const room = routerManager.getRoom(roomId);
      if (!room) {
        return res.status(404).json({ error: 'Room not found' });
      }

      const transportId = uuidv4();
      const transport = await transportManager.createTransport(
        room.router,
        transportId,
        type,
        peerId
      );

      res.json({
        transport_id: transportId,
        ice_params: transport.iceParameters,
        ice_candidates: transport.iceCandidates,
        dtls_params: transport.dtlsParameters,
      });
    } catch (error: any) {
      console.error('Error creating transport:', error);
      res.status(500).json({ error: error.message });
    }
  });

  // 关闭 Transport
  router.delete('/api/transports/:transportId', (req: any, res: any) => {
    transportManager.closeTransport(req.params.transportId);
    res.json({ ok: true });
  });

  // 关闭 Room
  router.delete('/api/rooms/:roomId', async (req: any, res: any) => {
    await routerManager.closeRoom(req.params.roomId);
    res.json({ ok: true });
  });

  // 获取 Room 信息
  router.get('/api/rooms/:roomId', async (req: any, res: any) => {
    const room = routerManager.getRoom(req.params.roomId);
    if (!room) {
      return res.status(404).json({ error: 'Room not found' });
    }

    const peers = peerManager.getPeers().filter(p => p.id.startsWith(req.params.roomId));

    res.json({
      room_id: room.id,
      peer_count: peers.length,
      created_at: room.createdAt,
    });
  });

  return router;
}