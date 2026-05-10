import { types } from 'mediasoup';
type Router = types.Router;
interface Room {
    id: string;
    router: Router;
    peers: Map<string, any>;
    createdAt: Date;
}
export declare class RouterManager {
    private workers;
    private rooms;
    private workerIndex;
    initialize(): Promise<void>;
    createRoom(roomId: string): Promise<Room>;
    getRoom(roomId: string): Room | undefined;
    closeRoom(roomId: string): Promise<void>;
    getRoomCount(): number;
}
export {};
//# sourceMappingURL=router.d.ts.map