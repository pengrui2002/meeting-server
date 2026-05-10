"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.Recorder = void 0;
const path = __importStar(require("path"));
class Recorder {
    constructor() {
        this.recordings = new Map();
    }
    async startRecording(roomId) {
        const recordingId = `rec_${Date.now()}`;
        const filePath = path.join('/tmp', `${recordingId}.mp4`);
        this.recordings.set(recordingId, {
            id: recordingId,
            roomId,
            startTime: new Date(),
            filePath,
            status: 'recording',
        });
        console.log(`Recording started: ${recordingId} for room ${roomId}`);
        return recordingId;
    }
    stopRecording(recordingId) {
        const recording = this.recordings.get(recordingId);
        if (recording) {
            recording.status = 'stopped';
            console.log(`Recording stopped: ${recordingId}`);
        }
        return recording;
    }
    getRecording(recordingId) {
        return this.recordings.get(recordingId);
    }
    getRecordingsByRoom(roomId) {
        return Array.from(this.recordings.values()).filter(r => r.roomId === roomId);
    }
}
exports.Recorder = Recorder;
//# sourceMappingURL=recorder.js.map