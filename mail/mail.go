package mail

import (
	"time"
	"net/textproto"
)

type Mail struct {
	ReceivedAt time.Time
	Sender string
	Recipients []string
	Header textproto.MIMEHeader
	Data []byte

	Multipart bool
	Parts []MailPart
}

type MailPart struct {
	Header textproto.MIMEHeader
	Disposition string
	DispositionParams map[string]string
	Data []byte
}
