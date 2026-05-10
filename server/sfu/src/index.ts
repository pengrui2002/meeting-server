import express from 'express';
import { RouterManager } from './router';
import { TransportManager } from './transport';
import { PeerManager } from './peer';
import { Recorder } from './recorder';
import { createSignalingRouter } from './signaling';
import { config } from './config';

async function main() {
  console.log('Starting SFU server...');

  // 初始化 Router 管理器
  const routerManager = new RouterManager();
  await routerManager.initialize();
  console.log('Router manager initialized');

  // 初始化其他管理器
  const transportManager = new TransportManager();
  const peerManager = new PeerManager();
  const recorder = new Recorder();

  // 创建 Express 应用
  const app = express();
  app.use(express.json());

  // 注册信令路由
  const signalingRouter = createSignalingRouter(routerManager, transportManager, peerManager);
  app.use(signalingRouter);

  // 健康检查
  app.get('/health', (req: any, res: any) => {
    res.json({
      status: 'ok',
      rooms: routerManager.getRoomCount(),
      peers: peerManager.getPeerCount(),
    });
  });

  // 启动服务器
  app.listen(config.port, () => {
    console.log(`SFU server listening on port ${config.port}`);
    console.log(`MediaSoup worker config: ${config.mediasoup.numWorkers} workers`);
    console.log(`RTC port range: ${config.mediasoup.worker.rtcMinPort}-${config.mediasoup.worker.rtcMaxPort}`);
  });
}

main().catch((error) => {
  console.error('Failed to start SFU server:', error);
  process.exit(1);
});