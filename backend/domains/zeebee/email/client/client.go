package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"strings"

	"golang.org/x/oauth2"

	commonemail "zeebee-backend/common/email"
	"zeebee-backend/domains/zeebee/constants"
)

type Client interface {
	SendEmail(ctx context.Context, msg commonemail.Message) (err error)
	// AuthCodeUrl(ctx context.Context) (authCodeUrl string)
	RefreshToken(ctx context.Context) (token *oauth2.Token, err error)
	GetEmail(ctx context.Context) (email string, err error)
	GetEmailType() (emailType constants.EmailType)
}

func buildRawMessage(msg commonemail.Message) (string, error) {
	var buf bytes.Buffer

	write := func(s string) { buf.WriteString(s) }
	encodedSubject := mime.BEncoding.Encode("utf-8", msg.Subject)
	write(fmt.Sprintf("From: %s\r\n", msg.From))
	if len(msg.To) > 0 {
		write(fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ", ")))
	}
	if len(msg.Cc) > 0 {
		write(fmt.Sprintf("Cc: %s\r\n", strings.Join(msg.Cc, ", ")))
	}
	write(fmt.Sprintf("Subject: %s\r\n", encodedSubject))
	write("MIME-Version: 1.0\r\n")

	writer := multipart.NewWriter(&buf)
	defer writer.Close()
	boundaryMixed := writer.Boundary()
	write(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%q\r\n", boundaryMixed))
	write("\r\n")

	hdr := make(textproto.MIMEHeader)
	if msg.IsHTML {
		hdr.Set("Content-Type", `text/html; charset="utf-8"`)
	} else {
		hdr.Set("Content-Type", `text/plain; charset="utf-8"`)
	}

	hdr.Set("Content-Transfer-Encoding", "7bit")
	partWriter, err := writer.CreatePart(hdr)
	if err != nil {
		return "", err
	}
	_, _ = partWriter.Write([]byte(msg.Body))

	for _, att := range msg.Attachments {
		h := make(textproto.MIMEHeader)
		if att.MimeType == "" {
			att.MimeType = "application/octet-stream"
		}
		h.Set("Content-Type", fmt.Sprintf("%s; name=%q", att.MimeType, att.Filename))
		h.Set("Content-Transfer-Encoding", "base64")
		h.Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, att.Filename))

		pw, err := writer.CreatePart(h)
		if err != nil {
			return "", err
		}

		// Write base64 encoding and wrap every 76 characters (RFC recommended)
		content, err := io.ReadAll(att.Reader)
		if err != nil {
			return "", err
		}
		b64 := base64.StdEncoding.EncodeToString(content)
		for i := 0; i < len(b64); i += 76 {
			end := i + 76
			if end > len(b64) {
				end = len(b64)
			}
			_, _ = pw.Write([]byte(b64[i:end] + "\r\n"))
		}
	}
	raw := base64.URLEncoding.EncodeToString(buf.Bytes())
	return raw, nil
}
