package mailer

import "fmt"

func NewTicketMessage(ticketID int64, title, description, submitterName string, adminEmails []string) Message {
	return Message{
		To:       adminEmails[0],
		CC:       adminEmails[1:],
		Subject:  fmt.Sprintf("New Ticket #%d: %s", ticketID, title),
		Body:     newTicketHTMLBody(ticketID, title, description, submitterName),
		TextBody: newTicketTextBody(ticketID, title, description, submitterName),
	}
}

func newTicketHTMLBody(ticketID int64, title, description, submittedBy string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="id" xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="x-apple-disable-message-reformatting">
  <title>Tiket Dukungan Baru #%%d</title>
</head>
<body style="margin: 0; padding: 0; background-color: #f4f6f8; -webkit-text-size-adjust: 100%%%%; -ms-text-size-adjust: 100%%%%; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;">

  <!-- Wrapper Outer -->
  <table role="presentation" width="100%%%%" cellpadding="0" cellspacing="0" border="0" bgcolor="#f4f6f8">
    <tr>
      <td align="center" style="padding: 40px 16px;">

        <!-- Main Card Container -->
        <table role="presentation" width="100%%%%" cellpadding="0" cellspacing="0" border="0" style="max-width: 580px; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08); border: 1px solid #e2e8f0;">
          
          <!-- Header Accent Strip -->
          <tr>
            <td height="4" bgcolor="#4f46e5" style="background-color: #4f46e5; font-size: 0; line-height: 0;">&nbsp;</td>
          </tr>

          <!-- Header Body -->
          <tr>
            <td style="padding: 32px 40px 24px 40px;">
              <table role="presentation" width="100%%%%" cellpadding="0" cellspacing="0" border="0">
                <tr>
                  <td valign="middle">
                    <!-- Brand Section dengan Logo SVG -->
                    <table role="presentation" cellpadding="0" cellspacing="0" border="0">
                      <tr>
                        <td valign="middle" style="padding-right: 12px;">
                          <img src="https://helpdesk.rs-erba.go.id/favicon-96x96.png" alt="Logo RS Ernaldi Bahar" width="40" height="40" style="display: block; border: 0; outline: none; text-decoration: none; width: 40px; height: 40px;">
                        </td>
                        <td valign="middle">
                          <p style="margin: 0; color: #4f46e5; font-size: 11px; font-weight: 700; letter-spacing: 1.5px; text-transform: uppercase;">
                            RS ERNALDI BAHAR
                          </p>
                          <h1 style="margin: 2px 0 0 0; color: #0f172a; font-size: 20px; font-weight: 700; line-height: 1.2;">
                            IT Helpdesk Support
                          </h1>
                        </td>
                      </tr>
                    </table>
                  </td>
                  <td align="right" valign="top">
                    <!-- Status Badge -->
                    <span style="display: inline-block; padding: 6px 12px; background-color: #fef3c7; color: #92400e; font-size: 11px; font-weight: 700; border-radius: 20px; text-transform: uppercase; letter-spacing: 0.5px;">
                      Perlu Review
                    </span>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Hero Ticket Box -->
          <tr>
            <td style="padding: 0 40px 24px 40px;">
              <table role="presentation" width="100%%%%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f8fafc; border: 1px dashed #cbd5e1; border-radius: 10px; padding: 20px;">
                <tr>
                  <td align="center">
                    <p style="margin: 0; color: #64748b; font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 1px;">
                      Nomor Referensi Tiket
                    </p>
                    <p style="margin: 4px 0 0 0; color: #4f46e5; font-size: 32px; font-weight: 800; letter-spacing: -0.5px; line-height: 1;">
                      #%d
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Divider -->
          <tr>
            <td style="padding: 0 40px;">
              <div style="border-top: 1px solid #f1f5f9; height: 1px;"></div>
            </td>
          </tr>

          <!-- Details Content -->
          <tr>
            <td style="padding: 24px 40px;">
              <table role="presentation" width="100%%%%" cellpadding="0" cellspacing="0" border="0">
                
                <!-- Submitter -->
                <tr>
                  <td valign="top" style="padding-bottom: 20px;">
                    <p style="margin: 0; color: #94a3b8; font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.5px;">
                      Pelapor / Pengirim
                    </p>
                    <p style="margin: 4px 0 0 0; color: #0f172a; font-size: 14px; font-weight: 600;">
                      %s
                    </p>
                  </td>
                </tr>

                <!-- Subject/Title -->
                <tr>
                  <td valign="top" style="padding-bottom: 20px;">
                    <p style="margin: 0; color: #94a3b8; font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.5px;">
                      Judul Kendala
                    </p>
                    <p style="margin: 4px 0 0 0; color: #0f172a; font-size: 15px; font-weight: 600; line-height: 1.4;">
                      %s
                    </p>
                  </td>
                </tr>

                <!-- Description Box -->
                <tr>
                  <td valign="top">
                    <p style="margin: 0; color: #94a3b8; font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.5px;">
                      Deskripsi Masalah
                    </p>
                    <div style="margin-top: 8px; background-color: #f8fafc; border-left: 3px solid #4f46e5; border-radius: 0 6px 6px 0; padding: 12px 16px;">
                      <p style="margin: 0; color: #334155; font-size: 14px; line-height: 1.6; white-space: pre-line;">
                        %s
                      </p>
                    </div>
                  </td>
                </tr>

              </table>
            </td>
          </tr>

          <!-- Call to Action Button -->
          <tr>
            <td style="padding: 8px 40px 32px 40px;" align="center">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                <tr>
                  <td align="center" bgcolor="#4f46e5" style="border-radius: 6px;">
                    <a href="https://helpdesk.rs-erba.go.id" target="_blank" style="display: inline-block; padding: 12px 28px; color: #ffffff; font-size: 14px; font-weight: 600; text-decoration: none; border-radius: 6px; background-color: #4f46e5;">
                      Buka Helpdesk Dashboard &rarr;
                    </a>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="background-color: #f8fafc; border-top: 1px solid #e2e8f0; padding: 24px 40px; text-align: center;">
              <p style="margin: 0; color: #64748b; font-size: 12px; line-height: 1.6;">
                Email ini dikirim secara otomatis oleh sistem IT Helpdesk.<br>
                <strong>RS Ernaldi Bahar</strong> &copy; All rights reserved.
              </p>
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>`, ticketID, submittedBy, title, description)
}

func newTicketTextBody(ticketID int64, title, description, submittedBy string) string {
	return fmt.Sprintf(
		"PEMBERITAHUAN TIKET BARU\n"+
			"========================\n\n"+
			"Nomor Referensi : #%d\n"+
			"Pelapor         : %s\n"+
			"Judul           : %s\n\n"+
			"Deskripsi Masalah\n"+
			"-----------------\n"+
			"%s\n\n"+
			"---\n"+
			"Email ini dikirim secara otomatis oleh sistem IT Helpdesk RS Ernaldi Bahar.\n"+
			"Mohon tidak membalas email ini.",
		ticketID, submittedBy, title, description,
	)
}
