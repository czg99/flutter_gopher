import 'dart:async';
import 'package:flutter/material.dart';
import 'package:{{.ProjectName}}/{{.ProjectName}}.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return const MaterialApp(
      title: 'Flutter Gopher Demo',
      home: MyHomePage(title: 'Flutter Gopher Demo Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  final lib = {{.LibClassName}}();
  int method = 0;
  String send = '';
  String recv = '';

  @override
  void initState() {
    super.initState();
  }

  Widget buildButton(String text, Color color, FutureOr<void> Function() onTap) {
    return MaterialButton(
      color: color,
      onPressed: () async {
        try {
          await onTap();
        } catch (e) {
          showToast(e.toString());
        }
      },
      child: Text(text),
    );
  }

  void showToast(String text) {
    if (context.mounted) {
      ScaffoldMessenger.of(context)
          .showSnackBar(SnackBar(content: Text(text)));
    }
  }

  void updateState(int method, String send, String recv) {
    setState(() {
      this.method = method;
      this.send = send;
      this.recv = recv;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.grey,
        title: Text(widget.title),
      ),
      body: Container(
        padding: const EdgeInsets.only(top: 20),
        alignment: Alignment.center,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: <Widget>[
            Text('Method: $method'),
            Text('Send: $send'),
            Text('Recv: $recv'),
            const SizedBox(height: 10),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                buildButton('Call go', Colors.blue, () {
                  final (method, send, recv) = lib.callGoMethod();
                  updateState(method, send, recv);
                }),
                buildButton('Call platform', Colors.green, () {
                  final (method, send, recv) = lib.callPlatformMethod();
                  updateState(method, send, recv);
                }),
              ],
            ),
            const SizedBox(height: 10),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                buildButton('Call go async', Colors.blueAccent,
                    () async {
                  final (method, send, recv) = await lib.callGoMethodAsync();
                  updateState(method, send, recv);
                }),
                buildButton('Call platform async', Colors.greenAccent,
                    () async {
                  final (method, send, recv) = await lib.callPlatformMethodAsync();
                  updateState(method, send, recv);
                }),
              ],
            )
          ],
        ),
      ),
    );
  }
}
