import * as fs from 'fs';
import * as path from 'path';

interface Recording {
  id: string;
  roomId: string;
  startTime: Date;
  filePath?: string;
  status: 'pending' | 'recording' | 'stopped';
}

export class Recorder {
  private recordings: Map<string, Recording> = new Map();

  async startRecording(roomId: string): Promise<string> {
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

  stopRecording(recordingId: string): Recording | undefined {
    const recording = this.recordings.get(recordingId);
    if (recording) {
      recording.status = 'stopped';
      console.log(`Recording stopped: ${recordingId}`);
    }
    return recording;
  }

  getRecording(recordingId: string): Recording | undefined {
    return this.recordings.get(recordingId);
  }

  getRecordingsByRoom(roomId: string): Recording[] {
    return Array.from(this.recordings.values()).filter(r => r.roomId === roomId);
  }
}