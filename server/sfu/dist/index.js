"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const router_1 = require("./router");
const transport_1 = require("./transport");
const peer_1 = require("./peer");
const recorder_1 = require("./recorder");
const signaling_1 = require("./signaling");
const config_1 = require("./config");
async function main() {
    console.log('Starting SFU server...');
    // 初始化 Router 管理器
    const routerManager = new router_1.RouterManager();
    await routerManager.initialize();
    console.log('Router manager initialized');
    // 初始化其他管理器
    const transportManager = new transport_1.TransportManager();
    const peerManager = new peer_1.PeerManager();
    const recorder = new recorder_1.Recorder();
    // 创建 Express 应用
    const app = (0, express_1.default)();
    app.use(express_1.default.json());
    // 注册信令路由
    const signalingRouter = (0, signaling_1.createSignalingRouter)(routerManager, transportManager, peerManager);
    app.use(signalingRouter);
    // 健康检查
    app.get('/health', (req, res) => {
        res.json({
            status: 'ok',
            rooms: routerManager.getRoomCount(),
            peers: peerManager.getPeerCount(),
        });
    });
    // 启动服务器
    app.listen(config_1.config.port, () => {
        console.log(`SFU server listening on port ${config_1.config.port}`);
        console.log(`MediaSoup worker config: ${config_1.config.mediasoup.numWorkers} workers`);
        console.log(`RTC port range: ${config_1.config.mediasoup.worker.rtcMinPort}-${config_1.config.mediasoup.worker.rtcMaxPort}`);
    });
}
main().catch((error) => {
    console.error('Failed to start SFU server:', error);
    process.exit(1);
});
//# sourceMappingURL=index.js.map