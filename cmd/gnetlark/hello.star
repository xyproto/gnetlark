def index(status, msg, method, path, date):
    return "HTTP/1.1 " + status + "\r\nServer: gnetlark\r\nDate: " + date + "\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n" + "Hello, World!"
