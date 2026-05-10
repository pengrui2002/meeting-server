import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../services/meeting_service.dart';
import 'meeting_page.dart';

class JoinMeetingPage extends StatefulWidget {
  const JoinMeetingPage({super.key});

  @override
  State<JoinMeetingPage> createState() => _JoinMeetingPageState();
}

class _JoinMeetingPageState extends State<JoinMeetingPage> {
  final _meetingIdController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _isLoading = false;

  @override
  void dispose() {
    _meetingIdController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _joinMeeting() async {
    if (_meetingIdController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入会议号')),
      );
      return;
    }

    setState(() => _isLoading = true);

    final meetingService = context.read<MeetingService>();
    final participant = await meetingService.joinMeeting(
      _meetingIdController.text,
      'user_${DateTime.now().millisecondsSinceEpoch}',
      password: _passwordController.text.isNotEmpty
          ? _passwordController.text
          : null,
    );

    setState(() => _isLoading = false);

    if (participant != null && mounted) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (_) => MeetingPage(meetingId: _meetingIdController.text),
        ),
      );
    } else if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('加入会议失败，请检查会议号和密码')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('加入会议'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextField(
              controller: _meetingIdController,
              decoration: const InputDecoration(
                labelText: '会议号',
                hintText: '请输入会议号',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _passwordController,
              decoration: const InputDecoration(
                labelText: '会议密码（可选）',
                hintText: '请输入会议密码',
                border: OutlineInputBorder(),
              ),
              obscureText: true,
            ),
            const SizedBox(height: 32),
            SizedBox(
              width: double.infinity,
              height: 48,
              child: ElevatedButton(
                onPressed: _isLoading ? null : _joinMeeting,
                child: _isLoading
                    ? const CircularProgressIndicator(color: Colors.white)
                    : const Text('加入会议'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}