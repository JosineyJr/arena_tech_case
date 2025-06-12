import http, { IncomingMessage, ServerResponse } from 'http';
import { parse } from 'url';

let server: http.Server | null = null;
export function start(port = 3000) {
  if (!server) {
    server = http.createServer((req: IncomingMessage, res: ServerResponse) => {
      const parsedUrl = parse(req.url || '', true);

      if (req.method === 'GET' && parsedUrl.pathname === '/ip/location') {
        res.statusCode = 404;
        res.setHeader('Content-Type', 'application/json');
        res.end(JSON.stringify({ error: 'Not Found' }));
        return;
      }

      // Default 404 response
      res.statusCode = 404;
      res.setHeader('Content-Type', 'application/json');
      res.end(JSON.stringify({ error: 'Not Found' }));
    });
  }

  server.listen(port, () => {
    console.log(`Server listening at http://localhost:${port}`);
  });
}

export function stop() {
  if (server) {
    server.close();
  }
  server = null;
}
