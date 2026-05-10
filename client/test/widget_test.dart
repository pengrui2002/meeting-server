import 'package:flutter_test/flutter_test.dart';
import 'package:meeting_client/main.dart';

void main() {
  testWidgets('App smoke test', (WidgetTester tester) async {
    await tester.pumpWidget(const MeetingApp());
    expect(find.text('视频会议'), findsOneWidget);
  });
}
