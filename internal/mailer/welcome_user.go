package mailer

import "fmt"

func NewWelcomeUserMessage(name, email, password string) Message {
	return Message{
		To:       email,
		Subject:  "Selamat Datang di IT Helpdesk — Akun Anda Telah Aktif",
		Body:     welcomeUserHTMLBody(name, email, password),
		TextBody: welcomeUserTextBody(name, email, password),
	}
}

func welcomeUserHTMLBody(name, email, password string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="id" xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="x-apple-disable-message-reformatting">
  <title>Selamat Datang di IT Helpdesk</title>
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
                    <span style="display: inline-block; padding: 6px 12px; background-color: #dcfce7; color: #166534; font-size: 11px; font-weight: 700; border-radius: 20px; text-transform: uppercase; letter-spacing: 0.5px;">
                      Akun Aktif
                    </span>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Welcome Banner -->
          <tr>
            <td style="padding: 0 40px 24px 40px;">
              <div style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 20px 24px;">
                <h2 style="margin: 0; color: #0f172a; font-size: 18px; font-weight: 700;">
                  Selamat Datang, %s! 👋
                </h2>
                <p style="margin: 6px 0 0 0; color: #64748b; font-size: 14px; line-height: 1.5;">
                  Akun IT Helpdesk Anda telah dibuat. Silakan gunakan kredensial di bawah ini untuk masuk ke dalam sistem pertama kali.
                </p>
              </div>
            </td>
          </tr>

          <!-- Credentials Box -->
          <tr>
            <td style="padding: 0 40px 24px 40px;">
              <table role="presentation" width="100%%%%" cellpadding="0" cellspacing="0" border="0" style="background-color: #ffffff; border: 1px solid #cbd5e1; border-radius: 10px; padding: 20px;">
                <!-- Email -->
                <tr>
                  <td style="padding-bottom: 16px; border-bottom: 1px solid #f1f5f9;">
                    <p style="margin: 0; color: #94a3b8; font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.5px;">
                      Email Akun
                    </p>
                    <p style="margin: 4px 0 0 0; color: #0f172a; font-size: 15px; font-weight: 600;">
                      %s
                    </p>
                  </td>
                </tr>
                <!-- Password -->
                <tr>
                  <td style="padding-top: 16px;">
                    <p style="margin: 0; color: #94a3b8; font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.5px;">
                      Kata Sandi Sementara
                    </p>
                    <div style="margin-top: 6px; display: inline-block; background-color: #f1f5f9; border: 1px solid #e2e8f0; padding: 8px 16px; border-radius: 6px;">
                      <span style="color: #4f46e5; font-family: 'Courier New', Courier, monospace; font-size: 18px; font-weight: 700; letter-spacing: 1.5px;">
                        %s
                      </span>
                    </div>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Security Notice -->
          <tr>
            <td style="padding: 0 40px 28px 40px;">
              <table role="presentation" width="100%%%%" cellpadding="0" cellspacing="0" border="0" style="background-color: #fefce8; border: 1px solid #fef08a; border-radius: 8px; padding: 14px 18px;">
                <tr>
                  <td valign="top">
                    <p style="margin: 0; color: #854d0e; font-size: 13px; font-weight: 700; line-height: 1.4;">
                      ⚠️ Peringatan Keamanan
                    </p>
                    <p style="margin: 4px 0 0 0; color: #713f12; font-size: 13px; line-height: 1.5;">
                      Demi keamanan akun, harap <strong>segera ubah kata sandi Anda</strong> setelah berhasil login pertama kali. Jangan berikan kredensial ini kepada siapa pun.
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Call to Action Button -->
          <tr>
            <td style="padding: 0 40px 32px 40px;" align="center">
              <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                <tr>
                  <td align="center" bgcolor="#4f46e5" style="border-radius: 6px;">
                    <a href="https://helpdesk.rs-erba.go.id/login" target="_blank" style="display: inline-block; padding: 12px 28px; color: #ffffff; font-size: 14px; font-weight: 600; text-decoration: none; border-radius: 6px; background-color: #4f46e5;">
                      Login ke IT Helpdesk &rarr;
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
</html>`, name, email, password)
}

func welcomeUserTextBody(name, email, password string) string {
	return fmt.Sprintf(
		"SELAMAT DATANG DI IT HELPDESK RS ERNALDI BAHAR\n"+
			"===============================================\n\n"+
			"Halo %s,\n\n"+
			"Akun IT Helpdesk Anda telah aktif. Gunakan kredensial berikut untuk login pertama kali:\n\n"+
			"Email                : %s\n"+
			"Kata Sandi Sementara : %s\n\n"+
			"PERINGATAN KEAMANAN:\n"+
			"Harap segera ubah kata sandi Anda setelah berhasil login pertama kali.\n"+
			"Jangan berikan kredensial ini kepada siapapun.\n\n"+
			"Login sekarang: https://helpdesk.rs-erba.go.id/login\n\n"+
			"---\n"+
			"Email ini dikirim secara otomatis oleh sistem IT Helpdesk RS Ernaldi Bahar.\n"+
			"Mohon tidak membalas email ini.",
		name, email, password,
	)
}
