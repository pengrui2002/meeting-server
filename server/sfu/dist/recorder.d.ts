interface Recording {
    id: string;
    roomId: string;
    startTime: Date;
    filePath?: string;
    status: 'pending' | 'recording' | 'stopped';
}
export declare class Recorder {
    private recordings;
    startRecording(roomId: string): Promise<string>;
    stopRecording(recordingId: string): Recording | undefined;
    getRecording(recordingId: string): Recording | undefined;
    getRecordingsByRoom(roomId: string): Recording[];
}
export {};
//# sourceMappingURL=recorder.d.ts.map