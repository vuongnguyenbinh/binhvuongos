#!/usr/bin/env python3
"""One-time script to get Google Drive OAuth2 refresh token.
Run this locally, then add GOOGLE_REFRESH_TOKEN to .env on server.
"""
import json
import http.server
import urllib.parse
import urllib.request
import webbrowser
import threading
import os

CLIENT_ID = os.environ.get("GOOGLE_CLIENT_ID", "")
CLIENT_SECRET = os.environ.get("GOOGLE_CLIENT_SECRET", "")
REDIRECT_URI = "http://localhost:8089"
SCOPE = "https://www.googleapis.com/auth/drive.file"

auth_code = None

class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        global auth_code
        query = urllib.parse.urlparse(self.path).query
        params = urllib.parse.parse_qs(query)
        auth_code = params.get("code", [None])[0]
        self.send_response(200)
        self.send_header("Content-Type", "text/html")
        self.end_headers()
        self.wfile.write(b"<h1>OK! You can close this tab.</h1>")

    def log_message(self, format, *args):
        pass

# Start local server
server = http.server.HTTPServer(("localhost", 8089), Handler)
thread = threading.Thread(target=server.handle_request)
thread.start()

# Open browser for auth
auth_url = (
    f"https://accounts.google.com/o/oauth2/auth?"
    f"client_id={CLIENT_ID}&redirect_uri={REDIRECT_URI}"
    f"&response_type=code&scope={SCOPE}&access_type=offline&prompt=consent"
)
print(f"\nOpen this URL in browser:\n{auth_url}\n")
webbrowser.open(auth_url)

thread.join()
server.server_close()

if not auth_code:
    print("ERROR: No auth code received")
    exit(1)

# Exchange code for tokens
data = urllib.parse.urlencode({
    "code": auth_code,
    "client_id": CLIENT_ID,
    "client_secret": CLIENT_SECRET,
    "redirect_uri": REDIRECT_URI,
    "grant_type": "authorization_code",
}).encode()

req = urllib.request.Request("https://oauth2.googleapis.com/token", data=data, method="POST")
resp = urllib.request.urlopen(req)
tokens = json.loads(resp.read())

print("\n=== TOKENS ===")
print(f"Access Token: {tokens.get('access_token', 'N/A')[:50]}...")
print(f"Refresh Token: {tokens.get('refresh_token', 'N/A')}")
print(f"\nAdd this to your .env:")
print(f"GOOGLE_REFRESH_TOKEN={tokens.get('refresh_token', '')}")
