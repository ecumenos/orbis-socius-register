package smtp_test

import (
	"testing"

	"github.com/ecumenos/orbis-socius-register/cmd/api/smtp"
	"github.com/stretchr/testify/require"
)

func TestSender_Send(t *testing.T) {
	t.Skip()
	s := smtp.NewSender(&smtp.SenderConfig{
		From:     "",
		Password: "",
		SMTPHost: "",
		SMTPPort: 587,
	})
	err := s.Send([]string{""}, "greeting", "text/plain", "Hello World!")
	require.NoError(t, err)
}
