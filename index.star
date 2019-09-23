def index(status, msg, date):
    retval = "HTTP/1.1 " + status + "\r\n"
    retval += "Server: gnetlark\r\n"
    retval += "Date: " + date + "\r\n"
    retval += "Content-Type: text/html; charset=utf-8\r\n"

    body = "<!doctype html><html><head><title>"
    body += "Gnetlark"
    body += "</title><style>"
    body += """
body {
  margin: 4em;
  color: #202020;
  background: #a0c0f0;
  font-family: courier;
}
h2 {
  font-family: sans-serif;
  border-bottom: 2px solid black;
  padding-bottom: 3px;
}
"""
    body += "</style></head><body>"
    body += """
<h2>Gnetlark</h2>
<p>Server is up and running.</p>
"""
    body += "</body></html>"

    if len(msg) > 0:
      body += "Got message:\n"
      body += msg + "\n"

    bodylen = len(body)
    if bodylen > 0:
      retval += "Content-Length: " + str(bodylen) + "\r\n"
      retval += "\r\n"
      retval += body

    return retval

def error(status, msg, date):
    retval = "HTTP/1.1 " + status + "\r\n"
    retval += "Server: gnetlark\r\n"
    retval += "Date: " + date + "\r\n"

    body = ""

    if len(msg) > 0:
      body += "Server error: " + msg + "\n"

    bodylen = len(body)
    if bodylen > 0:
      retval += "Content-Length: " + str(bodylen) + "\r\n"
      retval += "\r\n"
      retval += body

    return retval
