class Meeting {
  final String id;
  final String title;
  final String status;

  Meeting({
    required this.id,
    required this.title,
    required this.status,
  });

  factory Meeting.fromJson(Map<String, dynamic> json) {
    return Meeting(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      status: json['status'] ?? 'pending',
    );
  }
}

class Participant {
  final String id;
  final String meetingId;
  final String userId;
  final String role;

  Participant({
    required this.id,
    required this.meetingId,
    required this.userId,
    required this.role,
  });

  factory Participant.fromJson(Map<String, dynamic> json) {
    return Participant(
      id: json['id'] ?? '',
      meetingId: json['meeting_id'] ?? '',
      userId: json['user_id'] ?? '',
      role: json['role'] ?? 'participant',
    );
  }
}