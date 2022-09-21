package mailx

import (
	"bytes"
	"io"
)

// @author valor.

// Sender sends emails via *smtp.Client
type Sender struct {
	smtpClient
	from   string
	signer Signer
}

type Signer func([]byte) ([]byte, error)

// SetSigner Set Signer
func (s *Sender) SetSigner(signer Signer) {
	s.signer = signer
}

// Send sends the given emails.
func (s *Sender) Send(m *Message) error {
	rcpt, err := m.rcpt()
	if err != nil {
		return err
	}
	if m.header.from.Address == "" {
		m.header.from.Address = s.from
	}
	return s.send(s.from, rcpt, m)
}

// SendOne sends a message implements io.WriterTo
func (s *Sender) send(from string, to []string, msg io.WriterTo) error {
	if err := s.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err := s.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := s.Data()
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	if _, err = msg.WriteTo(buf); err != nil {
		_ = w.Close()
		return err
	}
	mes := buf.Bytes()
	if s.signer != nil {
		mes, err = s.signer(mes)
		if err != nil {
			_ = w.Close()
			return err
		}
	}

	if _, err = w.Write(mes); err != nil {
		_ = w.Close()
		return err
	}
	return w.Close()
}

// Close sends the QUIT command and closes the connection to the server.
func (s *Sender) Close() error {
	return s.Quit()
}
