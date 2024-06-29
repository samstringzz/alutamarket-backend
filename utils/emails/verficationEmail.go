package emails

import (
	"bytes"
	"html/template"
	"log"

	"gopkg.in/mail.v2"
)

type EmailData struct {
	Name string
	Code string
}

func SendVerificationMail(userName string, code string) {
	emailTemplate := `
        <html>
        <body>
		    <div style="text-align: center;">
                <img src="https://example.com/path/to/your/logo.png" alt="Aluta Market" style="max-width: 200px;"/>
            </div>
            <h1>Hello, {{.Name}}</h1>
            <p>Your verification code is: {{.Code}}</p>
        </body>
        </html>`

	data := EmailData{
		Name: userName,
		Code: code,
	}

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		log.Fatalf("error parsing template: %v", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		log.Fatalf("error executing template: %v", err)
	}

	m := mail.NewMessage()
	m.SetHeader("From", "Contact@thealutamarket.com")
	m.SetHeader("To", "folajimiopeyemisax13@gmail.com")
	m.SetHeader("Subject", "Verification Code from Alutamarket")
	m.SetBody("text/html", body.String())

	d := mail.NewDialer("smtp.gmail.com", 465, "folajimiopeyemisax13@gmail.com", "pfed wvbc mwxh xooa")
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("error sending email: %v", err)
	}

	log.Println("Email sent successfully!")
}
