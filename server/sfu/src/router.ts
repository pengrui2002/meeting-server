import { createWorker, types } from 'mediasoup';
import { config } from './config';

type Worker = types.Worker;
type Router = types.Router;

interface Room {
  id: string;
  router: Router;
  peers: Map<string, any>;
  createdAt: Date;
}

export class RouterManager {
  private workers: Worker[] = [];
  private rooms: Map<string, Room> = new Map();
  private workerIndex = 0;

  async initialize(): Promise<void> {
    for (let i = 0; i < config.mediasoup.numWorkers; i++) {
      const worker = await createWorker({
        ...config.mediasoup.worker,
      });
      this.workers.push(worker);
      console.log(`Worker ${i + 1} started`);
    }
  }

  async createRoom(roomId: string): Promise<Room> {
    if (this.rooms.has(roomId)) {
      return this.rooms.get(roomId)!;
    }

    const worker = this.workers[this.workerIndex % this.workers.length];
    this.workerIndex++;

    const router = await worker.createRouter({
      mediaCodecs: config.mediasoup.router.mediaCodecs,
    });

    const room: Room = {
      id: roomId,
      router,
      peers: new Map(),
      createdAt: new Date(),
    };

    this.rooms.set(roomId, room);
    console.log(`Room ${roomId} created with router`);

    return room;
  }

  getRoom(roomId: string): Room | undefined {
    return this.rooms.get(roomId);
  }

  async closeRoom(roomId: string): Promise<void> {
    const room = this.rooms.get(roomId);
    if (room) {
      await room.router.close();
      this.rooms.delete(roomId);
      console.log(`Room ${roomId} closed`);
    }
  }

  getRoomCount(): number {
    return this.rooms.size;
  }
}