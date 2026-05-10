export declare const config: {
    port: number;
    mediasoup: {
        numWorkers: number;
        worker: {
            rtcMinPort: number;
            rtcMaxPort: number;
        };
        router: {
            mediaCodecs: ({
                kind: "audio";
                mimeType: string;
                clockRate: number;
                channels: number;
                parameters?: undefined;
            } | {
                kind: "video";
                mimeType: string;
                clockRate: number;
                channels?: undefined;
                parameters?: undefined;
            } | {
                kind: "video";
                mimeType: string;
                clockRate: number;
                parameters: {
                    'packetization-mode': number;
                };
                channels?: undefined;
            })[];
        };
    };
};
//# sourceMappingURL=config.d.ts.map