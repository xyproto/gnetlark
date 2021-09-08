style = """
body {
  margin: 4em;
  color: #202020;
  background: black;
  font-family: courier;
  background-image:url(/gnetlark.svg);
  background-position: right 50px top 30px;
  background-repeat: no-repeat;
  background-attachment: fixed;
  background-size: 64px 64px;

}
h2 {
  color: red;
  font-family: sans-serif;
  border-bottom: 2px solid white;
  padding-bottom: 0.2em;
  text-shadow: black 0px 0px 5px;
  -webkit-font-smoothing: antialiased;
  margin-bottom: 2em;
}
p {
  margin-left: 2em;
  font-weight: bold;
  color: white;
}
"""

svgdata = """<svg xmlns="http://www.w3.org/2000/svg" width="359.887" height="318.059" viewBox="0 0 95.22 84.153"><g transform="translate(-46.506 -80.986)"><path d="M107.805 146.268l-9.345-3.662-26.702-4.762-.157-10.339a28.037 28.389 0 0 0 2.42 2.817 28.037 28.389 0 0 0 2.744 2.386v.007l.008-.002a28.037 28.389 0 0 0 31.13 2.002l-5.88-10.216 1.631-14.313h7.031zm2.88-36.08l1.126-14.104-20.694-4.762-5.34 16.118 10.346 14.286 1.77-4.405 4.13 7.18-.893 7.848-4.339 2.198 1.001-6.228-12.683-19.047-4.006 12.454 12.684 8.425-17.014 2.563a28.037 28.389 0 0 1-.008-.006l.667-16.11 10.013-27.84-8.01-2.197-8.01 29.304.176 11.64a28.037 28.389 0 0 1-2.123-31.463 28.037 28.389 0 0 1 31.538-13.283 28.037 28.389 0 0 1 20.808 27.429zm-12.793 7.133l-4.105-7.134h6.97zm2.864-7.134l4.045-10.073-1.147 10.073z" fill="red" fill-rule="evenodd" stroke="black" stroke-width="1.622"/><path d="M108.072 111.649l-9.374 2.997-26.782 3.896-.158 8.46a28.121 23.228 0 0 1 2.428-2.305 28.121 23.228 0 0 1 2.752-1.953v-.006l.008.002a28.121 23.228 0 0 1 31.223-1.637l-5.896 8.358 1.635 11.711h7.052zm2.888 29.522l1.129 11.54-20.756 3.896-5.357-13.188 10.378-11.689 1.775 3.604 4.144-5.874-.897-6.422-4.352-1.798 1.004 5.095-12.721 15.586-4.018-10.19 12.722-6.894-17.065-2.097a28.121 23.228 0 0 0-.008.005l.669 13.181 10.043 22.779-8.034 1.798-8.035-23.977.177-9.524a28.121 23.228 0 0 0-2.129 25.743 28.121 23.228 0 0 0 31.633 10.869 28.121 23.228 0 0 0 20.87-22.443zm-12.831-5.837l-4.118 5.837h6.99zm2.872 5.837l4.057 8.243-1.15-8.243z" fill="red" fill-rule="evenodd" stroke="black" stroke-width="1.469"/><rect ry="11.978" rx="8.242" y="99.847" x="47.776" height="48.835" width="92.68" fill="red" stroke="black" stroke-width="2.54" stroke-miterlimit="1"/><text y="129.017" x="52.702" style="line-height:125%;-inkscape-font-specification:'Fira Mono, Bold';font-variant-ligatures:normal;font-variant-caps:normal;font-variant-numeric:normal;font-feature-settings:normal;text-align:start" font-weight="700" font-size="16.933" font-family="Fira Mono" letter-spacing="0" word-spacing="0" stroke-width=".265"><tspan style="-inkscape-font-specification:'Fira Mono, Bold';font-variant-ligatures:normal;font-variant-caps:normal;font-variant-numeric:normal;font-feature-settings:normal;text-align:start" y="129.017" x="52.702">Gnetlark</tspan></text></g></svg>"""

def index(status, msg, method, path, date):
    OK = "HTTP/1.1 " + status + "\r\n"
    OK += "Server: gnetlark\r\n"
    OK += "Date: " + date + "\r\n"

    NOTOK = "HTTP/1.1 " + "404 Not found" + "\r\n"
    NOTOK += "Server: gnetlark\r\n"
    NOTOK += "Date: " + date + "\r\n"

    body = ""
    if method == "GET" and path == "/":
        retval = OK
        retval += "Content-Type: text/html; charset=utf-8\r\n"

        body += "<!doctype html><html><head><title>"
        body += "gnetlark"
        body += "</title><style>" + style + "</style></head><body>"
        body += "<h2>gnetlark</h2><p>Server is up and running.</p>"
        body += "</body></html>"

        if len(msg) > 0:
          body += "Got message:\n"
          body += msg + "\n"

    elif method == "GET" and path == "/gnetlark.svg":
        retval = OK
        retval += "Content-Type: image/svg+xml\r\n"
        body += svgdata

    elif method == "GET":
        retval = NOTOK
        retval += "Content-Type: text/plain; charset=utf-8\r\n"
        body += "Page not found: " + path

    bodylen = len(body)
    if bodylen > 0:
        retval += "Content-Length: " + str(bodylen) + "\r\n"
        retval += "\r\n"
        retval += body

    return retval

def error(status, msg, method, path, date):
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
