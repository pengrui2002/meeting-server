"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.RouterManager = void 0;
const mediasoup_1 = require("mediasoup");
const config_1 = require("./config");
class RouterManager {
    constructor() {
        this.workers = [];
        this.rooms = new Map();
        this.workerIndex = 0;
    }
    async initialize() {
        for (let i = 0; i < config_1.config.mediasoup.numWorkers; i++) {
            const worker = await (0, mediasoup_1.createWorker)({
                ...config_1.config.mediasoup.worker,
            });
            this.workers.push(worker);
            console.log(`Worker ${i + 1} started`);
        }
    }
    async createRoom(roomId) {
        if (this.rooms.has(roomId)) {
            return this.rooms.get(roomId);
        }
        const worker = this.workers[this.workerIndex % this.workers.length];
        this.workerIndex++;
        const router = await worker.createRouter({
            mediaCodecs: config_1.config.mediasoup.router.mediaCodecs,
        });
        const room = {
            id: roomId,
            router,
            peers: new Map(),
            createdAt: new Date(),
        };
        this.rooms.set(roomId, room);
        console.log(`Room ${roomId} created with router`);
        return room;
    }
    getRoom(roomId) {
        return this.rooms.get(roomId);
    }
    async closeRoom(roomId) {
        const room = this.rooms.get(roomId);
        if (room) {
            await room.router.close();
            this.rooms.delete(roomId);
            console.log(`Room ${roomId} closed`);
        }
    }
    getRoomCount() {
        return this.rooms.size;
    }
}
exports.RouterManager = RouterManager;
//# sourceMappingURL=router.js.map