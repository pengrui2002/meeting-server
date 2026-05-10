import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:flutter_webrtc/flutter_webrtc.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../config/app_config.dart';

class WebRTCService extends ChangeNotifier {
  RTCPeerConnection? _peerConnection;
  MediaStream? _localStream;
  final Map<String, MediaStream> _remoteStreams = {};
  WebSocketChannel? _channel;
  String? _roomId;
  String? _userId;
  bool _isConnected = false;
  bool _isMuted = false;
  bool _isVideoEnabled = true;

  RTCPeerConnection? get peerConnection => _peerConnection;
  MediaStream? get localStream => _localStream;
  Map<String, MediaStream> get remoteStreams => _remoteStreams;
  bool get isConnected => _isConnected;
  bool get isMuted => _isMuted;
  bool get isVideoEnabled => _isVideoEnabled;

  Future<void> joinRoom(String roomId, String userId) async {
    _roomId = roomId;
    _userId = userId;

    try {
      // 初始化本地媒体流
      final Map<String, dynamic> mediaConstraints = {
        'audio': true,
        'video': {
          'facingMode': 'user',
          'width': 640,
          'height': 480,
        },
      };

      _localStream = await navigator.mediaDevices.getUserMedia(mediaConstraints);

      // 创建 WebSocket 连接
      _channel = WebSocketChannel.connect(
        Uri.parse('${AppConfig.wsUrl}?room=$roomId&user=$userId'),
      );

      // 监听消息
      _channel!.stream.listen((message) {
        _handleSignalingMessage(message);
      });

      // 初始化 RTCPeerConnection
      await _initPeerConnection();

      _isConnected = true;
      notifyListeners();
    } catch (e) {
      debugPrint('Error joining room: $e');
    }
  }

  Future<void> _initPeerConnection() async {
    final Map<String, dynamic> config = {
      'iceServers': [
        {'urls': 'stun:stun.l.google.com:19302'},
      ]
    };

    _peerConnection = await createPeerConnection(config);

    // 添加本地流轨道
    if (_localStream != null) {
      _localStream!.getTracks().forEach((track) {
        _peerConnection!.addTrack(track, _localStream!);
      });
    }

    // 监听远程流
    _peerConnection!.onTrack = (RTCTrackEvent event) {
      if (event.streams.isNotEmpty) {
        _remoteStreams[event.streams[0].id] = event.streams[0];
        notifyListeners();
      }
    };
  }

  void _handleSignalingMessage(dynamic message) {
    try {
      final data = jsonDecode(message);

      switch (data['type']) {
        case 'offer':
          _handleOffer(data['sdp']);
          break;
        case 'answer':
          _handleAnswer(data['sdp']);
          break;
        case 'ice-candidate':
          _handleIceCandidate(data);
          break;
      }
    } catch (e) {
      debugPrint('Error handling signaling message: $e');
    }
  }

  Future<void> _handleOffer(String sdp) async {
    if (_peerConnection == null) return;
    await _peerConnection!.setRemoteDescription(
      RTCSessionDescription(sdp, 'offer'),
    );
    final answer = await _peerConnection!.createAnswer();
    await _peerConnection!.setLocalDescription(answer);
    _sendMessage({'type': 'answer', 'sdp': answer.sdp});
  }

  Future<void> _handleAnswer(String sdp) async {
    if (_peerConnection == null) return;
    await _peerConnection!.setRemoteDescription(
      RTCSessionDescription(sdp, 'answer'),
    );
  }

  Future<void> _handleIceCandidate(Map<String, dynamic> data) async {
    if (_peerConnection == null) return;
    final candidate = RTCIceCandidate(
      data['candidate'],
      data['sdpMid'],
      data['sdpMLineIndex'],
    );
    await _peerConnection!.addCandidate(candidate);
  }

  void _sendMessage(Map<String, dynamic> message) {
    _channel?.sink.add(jsonEncode(message));
  }

  Future<void> toggleMute() async {
    if (_localStream != null) {
      bool enabled = _localStream!.getAudioTracks().first.enabled;
      _localStream!.getAudioTracks().first.enabled = !enabled;
      _isMuted = !enabled;
      notifyListeners();
    }
  }

  Future<void> toggleVideo() async {
    if (_localStream != null) {
      bool enabled = _localStream!.getVideoTracks().first.enabled;
      _localStream!.getVideoTracks().first.enabled = !enabled;
      _isVideoEnabled = !enabled;
      notifyListeners();
    }
  }

  Future<void> switchCamera() async {
    if (_localStream != null) {
      await Helper.switchCamera(_localStream!.getVideoTracks().first);
    }
  }

  Future<void> leaveRoom() async {
    await _localStream?.dispose();
    await _peerConnection?.close();
    await _channel?.sink.close();
    _localStream = null;
    _peerConnection = null;
    _channel = null;
    _remoteStreams.clear();
    _isConnected = false;
    notifyListeners();
  }

  @override
  void dispose() {
    leaveRoom();
    super.dispose();
  }
}