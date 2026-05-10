import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_webrtc/flutter_webrtc.dart';
import 'package:provider/provider.dart';
import '../services/webrtc_service.dart';

class MeetingPage extends StatefulWidget {
  final String meetingId;

  const MeetingPage({super.key, required this.meetingId});

  @override
  State<MeetingPage> createState() => _MeetingPageState();
}

class _MeetingPageState extends State<MeetingPage> {
  final Map<String, RTCVideoRenderer> _renderers = {};
  bool _isInitialized = false;
  bool _isRecording = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _initRenderers();
    });
  }

  Future<void> _initRenderers() async {
    final webrtc = context.read<WebRTCService>();
    await webrtc.joinRoom(widget.meetingId, 'user_${DateTime.now().millisecondsSinceEpoch}');
    setState(() => _isInitialized = true);
  }

  Future<RTCVideoRenderer> _getRenderer(MediaStream stream) async {
    if (!_renderers.containsKey(stream.id)) {
      final renderer = RTCVideoRenderer();
      await renderer.initialize();
      _renderers[stream.id] = renderer;
    }
    return _renderers[stream.id]!;
  }

  @override
  void dispose() {
    for (final renderer in _renderers.values) {
      renderer.dispose();
    }
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.black,
        title: Text('会议号: ${widget.meetingId}'),
        actions: [
          IconButton(
            icon: const Icon(Icons.copy),
            onPressed: () {
              Clipboard.setData(ClipboardData(text: widget.meetingId));
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('会议号已复制')),
              );
            },
          ),
        ],
      ),
      body: Consumer<WebRTCService>(
        builder: (context, webrtc, child) {
          if (!_isInitialized) {
            return const Center(child: CircularProgressIndicator());
          }
          return Column(
            children: [
              Expanded(
                child: _buildVideoGrid(webrtc),
              ),
              _buildControlBar(webrtc),
            ],
          );
        },
      ),
    );
  }

  Widget _buildVideoGrid(WebRTCService webrtc) {
    final streams = webrtc.remoteStreams.values.toList();

    if (streams.isEmpty && webrtc.localStream == null) {
      return const Center(
        child: Text(
          '等待其他人加入...',
          style: TextStyle(color: Colors.white70),
        ),
      );
    }

    return GridView.builder(
      padding: const EdgeInsets.all(8),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 16 / 9,
        crossAxisSpacing: 8,
        mainAxisSpacing: 8,
      ),
      itemCount: streams.length + (webrtc.localStream != null ? 1 : 0),
      itemBuilder: (context, index) {
        if (webrtc.localStream != null && index == 0) {
          return _buildVideoView(webrtc.localStream!, '我', true);
        }
        final remoteIndex = webrtc.localStream != null ? index - 1 : index;
        return _buildVideoView(streams[remoteIndex], '用户$index', false);
      },
    );
  }

  Widget _buildVideoView(MediaStream? stream, String label, bool isLocal) {
    if (stream == null) {
      return Container(
        color: Colors.grey[900],
        child: const Center(
          child: Icon(Icons.videocam_off, color: Colors.white54, size: 48),
        ),
      );
    }

    return FutureBuilder<RTCVideoRenderer>(
      future: _getRenderer(stream),
      builder: (context, snapshot) {
        if (!snapshot.hasData) {
          return Container(
            color: Colors.grey[900],
            child: const Center(
              child: CircularProgressIndicator(),
            ),
          );
        }
        return ClipRRect(
          borderRadius: BorderRadius.circular(8),
          child: Stack(
            fit: StackFit.expand,
            children: [
              RTCVideoView(
                snapshot.data!,
                objectFit: RTCVideoViewObjectFit.RTCVideoViewObjectFitCover,
              ),
              Positioned(
                bottom: 8,
                left: 8,
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: Colors.black54,
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    label,
                    style: const TextStyle(color: Colors.white, fontSize: 12),
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildControlBar(WebRTCService webrtc) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 24),
      color: Colors.grey[900],
      child: SafeArea(
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _buildControlButton(
              icon: webrtc.isMuted ? Icons.mic_off : Icons.mic,
              label: webrtc.isMuted ? '取消静音' : '静音',
              onPressed: () => webrtc.toggleMute(),
              isActive: !webrtc.isMuted,
            ),
            _buildControlButton(
              icon: webrtc.isVideoEnabled ? Icons.videocam : Icons.videocam_off,
              label: webrtc.isVideoEnabled ? '关闭视频' : '开启视频',
              onPressed: () => webrtc.toggleVideo(),
              isActive: webrtc.isVideoEnabled,
            ),
            _buildControlButton(
              icon: Icons.cameraswitch,
              label: '切换摄像头',
              onPressed: () => webrtc.switchCamera(),
              isActive: true,
            ),
            _buildControlButton(
              icon: Icons.call_end,
              label: '挂断',
              onPressed: () {
                webrtc.leaveRoom();
                Navigator.pop(context);
              },
              isActive: false,
              isEndCall: true,
            ),
            _buildControlButton(
              icon: _isRecording ? Icons.stop_circle : Icons.fiber_manual_record,
              label: _isRecording ? '停止录制' : '开始录制',
              onPressed: () {
                setState(() {
                  _isRecording = !_isRecording;
                });
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(_isRecording ? '开始录制' : '停止录制'),
                  ),
                );
              },
              isActive: !_isRecording,
              isRecording: _isRecording,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildControlButton({
    required IconData icon,
    required String label,
    required VoidCallback onPressed,
    required bool isActive,
    bool isEndCall = false,
    bool isRecording = false,
  }) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        IconButton(
          onPressed: onPressed,
          icon: Icon(icon),
          style: IconButton.styleFrom(
            backgroundColor: isEndCall
                ? Colors.red
                : (isActive ? Colors.white24 : Colors.white10),
            padding: const EdgeInsets.all(16),
          ),
          iconSize: 28,
          color: isRecording ? Colors.red : Colors.white,
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: const TextStyle(color: Colors.white70, fontSize: 12),
        ),
      ],
    );
  }
}