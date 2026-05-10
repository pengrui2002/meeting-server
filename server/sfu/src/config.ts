export const config = {
  port: parseInt(process.env.SFU_PORT || '3000'),
  mediasoup: {
    numWorkers: 4,
    worker: {
      rtcMinPort: 40000,
      rtcMaxPort: 40100,
    },
    router: {
      mediaCodecs: [
        {
          kind: 'audio' as const,
          mimeType: 'audio/opus',
          clockRate: 48000,
          channels: 2,
        },
        {
          kind: 'video' as const,
          mimeType: 'video/VP8',
          clockRate: 90000,
        },
        {
          kind: 'video' as const,
          mimeType: 'video/VP9',
          clockRate: 90000,
        },
        {
          kind: 'video' as const,
          mimeType: 'video/H264',
          clockRate: 90000,
          parameters: {
            'packetization-mode': 1,
          },
        },
      ],
    },
  },
};