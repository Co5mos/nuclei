package smtp

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"net/textproto"
	"strings"
)

type (
	// SMTPMessage is a message to be sent over SMTP
	// @example
	// ```javascript
	// const smtp = require('nuclei/smtp');
	// const message = new smtp.SMTPMessage();
	// message.From('xyz@projectdiscovery.io');
	// ```
	SMTPMessage struct {
		from       string
		to         []string
		sub        string
		msg        []byte
		user       string
		pass       string
		attachment string
		attachData []byte
	}
)

// From adds the from field to the message
// @example
// ```javascript
// const smtp = require('nuclei/smtp');
// const message = new smtp.SMTPMessage();
// message.From('xyz@projectdiscovery.io');
// ```
func (s *SMTPMessage) From(email string) *SMTPMessage {
	s.from = email
	return s
}

// To adds the to field to the message
// @example
// ```javascript
// const smtp = require('nuclei/smtp');
// const message = new smtp.SMTPMessage();
// message.To('xyz@projectdiscovery.io');
// ```
func (s *SMTPMessage) To(email string) *SMTPMessage {
	s.to = append(s.to, email)
	return s
}

// Subject adds the subject field to the message
// @example
// ```javascript
// const smtp = require('nuclei/smtp');
// const message = new smtp.SMTPMessage();
// message.Subject('hello');
// ```
func (s *SMTPMessage) Subject(sub string) *SMTPMessage {
	s.sub = sub
	return s
}

// Body adds the message body to the message
// @example
// ```javascript
// const smtp = require('nuclei/smtp');
// const message = new smtp.SMTPMessage();
// message.Body('hello');
// ```
func (s *SMTPMessage) Body(msg []byte) *SMTPMessage {
	s.msg = msg
	return s
}

// Auth when called authenticates using username and password before sending the message
// @example
// ```javascript
// const smtp = require('nuclei/smtp');
// const message = new smtp.SMTPMessage();
// message.Auth('username', 'password');
// ```
func (s *SMTPMessage) Auth(username, password string) *SMTPMessage {
	s.user = username
	s.pass = password
	return s
}

// String returns the string representation of the message
// @example
// ```javascript
// const smtp = require('nuclei/smtp');
// const message = new smtp.SMTPMessage();
// message.From('xyz@projectdiscovery.io');
// message.To('xyz2@projectdiscoveyr.io');
// message.Subject('hello');
// message.Body('hello');
// message.Attachment('file.txt', 'hello');
// log(message.String());
// ```
func (s *SMTPMessage) String() string {
	var buff bytes.Buffer
	tw := textproto.NewWriter(bufio.NewWriter(&buff))

	_ = tw.PrintfLine("To: %s", strings.Join(s.to, ","))
	if s.sub != "" {
		_ = tw.PrintfLine("Subject: %s", s.sub)
	}

	_ = tw.PrintfLine("MIME-Version: 1.0")
	if s.attachment != "" {
		boundary := "my-boundary-12345"
		_ = tw.PrintfLine("Content-Type: multipart/mixed; boundary=%s", boundary)
		_ = tw.PrintfLine("\r\n--%s", boundary)
		_ = tw.PrintfLine("Content-Type: text/plain; charset=\"utf-8\"")
		_ = tw.PrintfLine("\r\n%s", s.msg)
		_ = tw.PrintfLine("\r\n--%s", boundary)
		_ = tw.PrintfLine("Content-Type: application/octet-stream; name=\"%s\"", s.attachment)
		_ = tw.PrintfLine("Content-Transfer-Encoding: base64")
		_ = tw.PrintfLine("Content-Disposition: attachment; filename=\"%s\"", s.attachment)
		encoded := base64.StdEncoding.EncodeToString(s.attachData)
		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			_ = tw.PrintfLine("%s", encoded[i:end])
		}
		_ = tw.PrintfLine("\r\n--%s--", boundary)
	} else {
		_ = tw.PrintfLine("\r\n%s", s.msg)
	}
	return buff.String()
}

// Attachment adds an attachment to the message
// @example
// ```javascript
// const smtp = require('nuclei/smtp');
// const message = new smtp.SMTPMessage();
// message.Attachment('file.txt', 'hello');
// ```
func (s *SMTPMessage) Attachment(filename string, data []byte) *SMTPMessage {
	s.attachment = filename
	s.attachData = data
	return s
}
