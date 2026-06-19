package mailer

import "fmt"

func NewWelcomeUserMessage(name, email, password string) Message {
	return Message{
		To:       email,
		Subject:  "Welcome to IT Helpdesk — Your Account is Ready",
		Body:     welcomeUserHTMLBody(name, email, password),
		TextBody: welcomeUserTextBody(name, email, password),
	}
}

func welcomeUserHTMLBody(name, email, password string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="x-apple-disable-message-reformatting">
  <title>Welcome to IT Helpdesk</title>
</head>
<body style="margin:0;padding:0;background-color:#f1f5f9;-webkit-text-size-adjust:100%%;-ms-text-size-adjust:100%%;">

  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" bgcolor="#f1f5f9">
    <tr>
      <td align="center" style="padding:48px 20px;">

        <table role="presentation" width="600" cellpadding="0" cellspacing="0" border="0"
          style="max-width:600px;width:100%%;border-radius:16px;overflow:hidden;box-shadow:0 8px 32px rgba(0,0,0,0.10);">

          <!-- Header -->
          <tr>
            <td bgcolor="#0f172a" style="background-color:#0f172a;padding:40px 48px;text-align:center;">
              <p style="margin:0;color:#475569;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;font-size:11px;font-weight:700;letter-spacing:3px;text-transform:uppercase;">IT Helpdesk</p>
              <h1 style="margin:12px 0 0;color:#f8fafc;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;font-size:26px;font-weight:700;letter-spacing:-0.5px;line-height:1.3;">Welcome, %s!</h1>
            </td>
          </tr>

          <!-- Accent bar -->
          <tr>
            <td bgcolor="#6366f1" style="background-color:#6366f1;height:4px;font-size:4px;line-height:4px;">&nbsp;</td>
          </tr>

          <!-- Body -->
          <tr>
            <td bgcolor="#ffffff" style="background-color:#ffffff;padding:40px 48px 32px;">

              <p style="margin:0 0 28px;color:#64748b;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;font-size:15px;line-height:1.7;">Your IT Helpdesk account has been created. Use the credentials below to sign in for the first time.</p>

              <!-- Credentials -->
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0">
                <tr>
                  <td bgcolor="#f8fafc" style="background-color:#f8fafc;border:1px solid #e2e8f0;border-radius:10px;padding:24px 28px;">

                    <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0">
                      <tr>
                        <td style="padding-bottom:16px;border-bottom:1px solid #e2e8f0;">
                          <p style="margin:0;color:#94a3b8;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;font-size:10px;font-weight:700;letter-spacing:2.5px;text-transform:uppercase;">Email</p>
                          <p style="margin:6px 0 0;color:#0f172a;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;font-size:16px;font-weight:600;">%s</p>
                        </td>
                      </tr>
                      <tr>
                        <td style="padding-top:16px;">
                          <p style="margin:0;color:#94a3b8;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;font-size:10px;font-weight:700;letter-spacing:2.5px;text-transform:uppercase;">Temporary Password</p>
                          <p style="margin:6px 0 0;color:#6366f1;font-family:'Courier New',Courier,monospace;font-size:20px;font-weight:700;letter-spacing:2px;">%s</p>
                        </td>
                      </tr>
                    </table>

                  </td>
                </tr>
              </table>

              <!-- Spacer -->
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0">
                <tr><td style="height:28px;line-height:28px;">&nbsp;</td></tr>
              </table>

              <!-- Warning -->
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0">
                <tr>
                  <td bgcolor="#fef9c3" style="background-color:#fef9c3;border:1px solid #fde047;border-radius:10px;padding:16px 20px;">
                    <p style="margin:0;color:#854d0e;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;font-size:13px;font-weight:700;line-height:1.5;">&#9888;&nbsp; Security Notice</p>
                    <p style="margin:6px 0 0;color:#713f12;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;font-size:13px;line-height:1.7;">Please <strong>change your password immediately</strong> after your first login. Do not share your credentials with anyone.</p>
                  </td>
                </tr>
              </table>

            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td bgcolor="#f8fafc" style="background-color:#f8fafc;border-top:1px solid #e2e8f0;padding:24px 48px;text-align:center;">
              <p style="margin:0;color:#94a3b8;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif;font-size:12px;line-height:1.8;">This is an automated notification. Please do not reply to this email.<br>&copy; IT Helpdesk. Ernaldi Bahar Hospital. All rights reserved.</p>
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
		"WELCOME TO IT HELPDESK\n"+
			"======================\n\n"+
			"Hello %s,\n\n"+
			"Your account has been created. Use the credentials below to sign in:\n\n"+
			"Email             : %s\n"+
			"Temporary Password: %s\n\n"+
			"IMPORTANT: Please change your password immediately after your first login.\n"+
			"Do not share your credentials with anyone.\n\n"+
			"---\n"+
			"This is an automated notification from the Helpdesk system.\n"+
			"Please do not reply to this email.",
		name, email, password,
	)
}
