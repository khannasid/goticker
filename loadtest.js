import ws from 'k6/ws';
import { check } from 'k6';

export default function () {
  const url = 'ws://localhost:8080/ws';
  const params = { tags: { my_tag: 'hello' } };

  const response = ws.connect(url, params, function (socket) {
    socket.on('open', function open() {
      // socket.send('Hello');
    });

    socket.on('message', function (message) {
      // Validate we are receiving JSON
      const msg = JSON.parse(message);
      check(msg, { 'Price is valid': (m) => m.price > 0 });
    });

    socket.on('close', () => console.log('disconnected'));
    
    // Keep connection open for 10 seconds
    socket.setTimeout(function () {
      socket.close();
    }, 10000);
  });

  check(response, { 'status is 101': (r) => r && r.status === 101 });
}