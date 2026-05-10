import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import '../config/app_config.dart';
import '../models/meeting.dart';

class MeetingService extends ChangeNotifier {
  Meeting? _currentMeeting;
  final List<Participant> _participants = [];
  bool _isLoading = false;

  Meeting? get currentMeeting => _currentMeeting;
  List<Participant> get participants => _participants;
  bool get isLoading => _isLoading;

  Future<Meeting?> createMeeting(String title, {String? password}) async {
    _isLoading = true;
    notifyListeners();

    try {
      final response = await http.post(
        Uri.parse('${AppConfig.apiBase}/meeting/create'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'title': title,
          'password': password ?? '',
        }),
      );

      if (response.statusCode == 200) {
        _currentMeeting = Meeting.fromJson(jsonDecode(response.body));
      }
    } catch (e) {
      debugPrint('Error creating meeting: $e');
    }

    _isLoading = false;
    notifyListeners();
    return _currentMeeting;
  }

  Future<Meeting?> getMeeting(String id) async {
    try {
      final response = await http.get(
        Uri.parse('${AppConfig.apiBase}/meeting/get?id=$id'),
      );

      if (response.statusCode == 200) {
        _currentMeeting = Meeting.fromJson(jsonDecode(response.body));
        notifyListeners();
        return _currentMeeting;
      }
    } catch (e) {
      debugPrint('Error getting meeting: $e');
    }
    return null;
  }

  Future<Participant?> joinMeeting(String meetingId, String userId, {String? password}) async {
    try {
      final response = await http.post(
        Uri.parse('${AppConfig.apiBase}/meeting/join'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'meeting_id': meetingId,
          'user_id': userId,
          'password': password ?? '',
        }),
      );

      if (response.statusCode == 200) {
        final participant = Participant.fromJson(jsonDecode(response.body));
        _participants.add(participant);
        notifyListeners();
        return participant;
      }
    } catch (e) {
      debugPrint('Error joining meeting: $e');
    }
    return null;
  }

  void setCurrentMeeting(Meeting? meeting) {
    _currentMeeting = meeting;
    notifyListeners();
  }
}