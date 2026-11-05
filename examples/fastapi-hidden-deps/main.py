import os
from http.server import BaseHTTPRequestHandler, HTTPServer

if not os.getenv("DATABASE_URL"):
    raise SystemExit("missing DATABASE_URL")
if not os.getenv("APP_MODE") == "demo":
    raise SystemExit("APP_MODE must be demo")

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"ok")
    def log_message(self, *_): pass

HTTPServer(("0.0.0.0", 8000), Handler).serve_forever()
