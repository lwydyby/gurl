package gurl

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
)

type Curl struct {
	Method        string
	Get           bool
	Header        []string
	Data          string
	DataRaw       string
	Form          []string
	URL           string
	Location      bool
	DataUrlencode []string

	Compressed bool
	Include    bool
	Insecure   bool
	Err        error
}

const (
	bodyEmpby     = "empty"
	bodyURLEncode = "data-urlencode"
	bodyForm      = "form"
	bodyData      = "data"
	bodyDataRaw   = "data-raw"
)

func (c *Curl) Request() *http.Request {
	// Set request method
	if !strings.Contains(c.URL, "http") {
		c.URL = "http://" + c.URL
	}
	c.initMethod()
	var req *http.Request
	var err error
	// Create request
	if c.Method == "GET" {
		req, err = http.NewRequest(c.Method, c.URL, nil)
	} else {
		var body bytes.Buffer
		switch c.findHighestPriority() {
		case bodyURLEncode:
			body.WriteString(strings.Join(c.getWWWForm(), "&"))
		case bodyForm:
			body.WriteString(strings.Join(c.getForm(), "&"))
		case bodyData:
			body.WriteString(c.Data)
		case bodyDataRaw:
			body.WriteString(c.getDataRaw())
		}
		req, err = http.NewRequest(c.Method, c.URL, &body)
	}

	if err != nil {
		panic(err)
	}

	// Set headers
	headers := c.getHeader()
	for i := 0; i < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	return req
}

func (c *Curl) initMethod() {
	if len(c.Method) == 0 && c.Get {
		c.Method = "GET"
		return
	}

	if len(c.Method) != 0 {
		return
	}

	if len(c.Data) > 0 {
		c.Method = "POST"
		return
	}

	c.Method = "GET"
}

func (c *Curl) getHeader() []string {
	if len(c.Header) == 0 {
		return nil
	}

	header := make([]string, len(c.Header)*2)
	index := 0
	for _, v := range c.Header {
		pos := strings.IndexByte(v, ':')
		if pos == -1 {
			continue
		}

		header[index] = v[:pos]
		index++
		header[index] = v[pos+1:]
		index++
	}

	return header
}

func (c *Curl) findHighestPriority() string {
	if len(c.DataUrlencode) != 0 {
		return bodyURLEncode
	}
	if len(c.Form) != 0 {
		return bodyForm
	}
	if len(c.Data) != 0 {
		return bodyData
	}
	if len(c.DataRaw) != 0 {
		return bodyDataRaw
	}

	return bodyEmpby
}

func (c *Curl) getDataRaw() string {
	if len(c.DataRaw) == 0 || !strings.Contains(c.DataRaw, "@") {
		return c.DataRaw
	}
	return getFromFile(c.DataRaw[1:])
}

func (c *Curl) getWWWForm() []string {
	if len(c.DataUrlencode) == 0 {
		return []string{}
	}

	form := make([]string, len(c.DataUrlencode)*2)
	index := 0
	for _, v := range c.DataUrlencode {
		pos := strings.IndexByte(v, '=')
		if pos == -1 {
			continue
		}

		form[index] = v[:pos]
		index++
		form[index] = v[pos+1:]
		index++
	}

	return form
}

func (c *Curl) getForm() []string {
	if len(c.Form) == 0 {
		return []string{}
	}

	form := make([]string, len(c.Form)*2)
	index := 0
	for _, v := range c.Form {
		pos := strings.IndexByte(v, '=')
		if pos == -1 {
			continue
		}

		form[index] = v[:pos]
		fieldValue := v[pos+1:]
		// 删除value中的双引号
		fieldValue = strings.Trim(fieldValue, "\"")
		if len(fieldValue) > 0 && fieldValue[0] == '@' {
			form[index] = getFromFile(fieldValue[1:])
		} else {
			form[index] = fieldValue
		}

		index++
	}

	return form
}

func getFromFile(path string) string {
	body, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	value, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return string(value)
}

func (c *Curl) getURL() string {
	return c.URL
}
