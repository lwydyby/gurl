# 简介
这是一个用go写的类似curl的命令行工具,使用方式和curl一致

# 安装
```bash
go install github.com/lwydyby/gurl/gurl@latest
```

# 使用说明

```bash
gurl --help
gurl [options...] <url>

Usage:
  gurl [flags]

Flags:
      --compressed                   Request compressed response
  -d, --data string                  HTTP POST data
      --data-raw string              HTTP POST data, '@' allowed
      --data-urlencode stringArray   HTTP POST data url encoded
  -f, --form stringArray             Specify multipart MIME data
  -G, --get                          Put the post data in the URL and use GET
  -H, --header stringArray           Pass custom header(s) to server
  -h, --help                         help for gurl
  -i, --include                      Include the HTTP response headers in the output. The HTTP response headers can include things like server name, cookies, date of the document, HTTP version and more.
  -k, --insecure                     Allow insecure server connections when using SSL
  -L, --location                     Follow redirects
  -X, --request string               Specify request method to use
      --url string                   URL to work with
```