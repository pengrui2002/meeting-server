"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createSignalingRouter = createSignalingRouter;
function createSignalingRouter(routerManager, transportManager, peerManager) {
    const express = require('express');
    const router = express.Router();
    const { v4: uuidv4 } = require('uuid');
    // 创建/获取 Room
    router.post('/api/rooms', async (req, res) => {
        try {
            const { meetingId } = req.body;
            const roomId = meetingId || uuidv4();
            const room = await routerManager.createRoom(roomId);
            res.json({
                room_id: room.id,
                rtp_capabilities: room.router.rtpCapabilities,
            });
        }
        catch (error) {
            console.error('Error creating room:', error);
            res.status(500).json({ error: error.message });
        }
    });
    // 创建 Transport
    router.post('/api/rooms/:roomId/transports', async (req, res) => {
        try {
            const { roomId } = req.params;
            const { type, peerId } = req.body;
            const room = routerManager.getRoom(roomId);
            if (!room) {
                return res.status(404).json({ error: 'Room not found' });
            }
            const transportId = uuidv4();
            const transport = await transportManager.createTransport(room.router, transportId, type, peerId);
            res.json({
                transport_id: transportId,
                ice_params: transport.iceParameters,
                ice_candidates: transport.iceCandidates,
                dtls_params: transport.dtlsParameters,
            });
        }
        catch (error) {
            console.error('Error creating transport:', error);
            res.status(500).json({ error: error.message });
        }
    });
    // 关闭 Transport
    router.delete('/api/transports/:transportId', (req, res) => {
        transportManager.closeTransport(req.params.transportId);
        res.json({ ok: true });
    });
    // 关闭 Room
    router.delete('/api/rooms/:roomId', async (req, res) => {
        await routerManager.closeRoom(req.params.roomId);
        res.json({ ok: true });
    });
    // 获取 Room 信息
    router.get('/api/rooms/:roomId', async (req, res) => {
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
//# sourceMappingURL=signaling.js.map