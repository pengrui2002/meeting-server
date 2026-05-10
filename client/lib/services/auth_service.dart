import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import '../config/app_config.dart';
import '../models/user.dart';

class AuthService extends ChangeNotifier {
  User? _currentUser;
  String? _token;
  bool _isLoading = false;

  User? get currentUser => _currentUser;
  String? get token => _token;
  bool get isLoading => _isLoading;
  bool get isLoggedIn => _token != null;

  Future<AuthResponse?> register(String username, String password, String displayName) async {
    _isLoading = true;
    notifyListeners();

    try {
      final response = await http.post(
        Uri.parse('${AppConfig.apiBase}/auth/register'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'username': username,
          'password': password,
          'display_name': displayName,
        }),
      );

      if (response.statusCode == 200) {
        final authResponse = AuthResponse.fromJson(jsonDecode(response.body));
        _currentUser = authResponse.user;
        _token = authResponse.token;
        _isLoading = false;
        notifyListeners();
        return authResponse;
      }
    } catch (e) {
      debugPrint('Error registering: $e');
    }

    _isLoading = false;
    notifyListeners();
    return null;
  }

  Future<AuthResponse?> login(String username, String password) async {
    _isLoading = true;
    notifyListeners();

    try {
      final response = await http.post(
        Uri.parse('${AppConfig.apiBase}/auth/login'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'username': username,
          'password': password,
        }),
      );

      if (response.statusCode == 200) {
        final authResponse = AuthResponse.fromJson(jsonDecode(response.body));
        _currentUser = authResponse.user;
        _token = authResponse.token;
        _isLoading = false;
        notifyListeners();
        return authResponse;
      }
    } catch (e) {
      debugPrint('Error logging in: $e');
    }

    _isLoading = false;
    notifyListeners();
    return null;
  }

  void logout() {
    _currentUser = null;
    _token = null;
    notifyListeners();
  }

  void setToken(String token) {
    _token = token;
    notifyListeners();
  }
}