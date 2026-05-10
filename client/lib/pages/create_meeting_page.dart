import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../services/meeting_service.dart';
import 'meeting_page.dart';

class CreateMeetingPage extends StatefulWidget {
  const CreateMeetingPage({super.key});

  @override
  State<CreateMeetingPage> createState() => _CreateMeetingPageState();
}

class _CreateMeetingPageState extends State<CreateMeetingPage> {
  final _titleController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _isLoading = false;

  @override
  void dispose() {
    _titleController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _createMeeting() async {
    if (_titleController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入会议主题')),
      );
      return;
    }

    setState(() => _isLoading = true);

    final meetingService = context.read<MeetingService>();
    final meeting = await meetingService.createMeeting(
      _titleController.text,
      password: _passwordController.text.isNotEmpty
          ? _passwordController.text
          : null,
    );

    setState(() => _isLoading = false);

    if (meeting != null && mounted) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (_) => MeetingPage(meetingId: meeting.id),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('创建会议'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextField(
              controller: _titleController,
              decoration: const InputDecoration(
                labelText: '会议主题',
                hintText: '请输入会议主题',
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
                onPressed: _isLoading ? null : _createMeeting,
                child: _isLoading
                    ? const CircularProgressIndicator(color: Colors.white)
                    : const Text('创建会议'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}