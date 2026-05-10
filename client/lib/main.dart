import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'pages/home_page.dart';
import 'services/meeting_service.dart';
import 'services/webrtc_service.dart';

void main() {
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => MeetingService()),
        ChangeNotifierProvider(create: (_) => WebRTCService()),
      ],
      child: const MeetingApp(),
    ),
  );
}

class MeetingApp extends StatelessWidget {
  const MeetingApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: '视频会议',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        useMaterial3: true,
      ),
      home: const HomePage(),
    );
  }
}