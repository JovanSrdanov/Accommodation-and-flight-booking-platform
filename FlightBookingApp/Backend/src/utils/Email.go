package utils

import (
	"FlightBookingApp/model"
	"log"
	"net/smtp"
)

func SendConfirmationEmail(account model.Account, emailVerPassword string) error {
	domName := "http://localhost:4200"
	subject := "Email Verification"
	HTMLbody :=
		`
			<html>
				<h1>Hi ` + account.Username + `,</h1>
				<p>Thanks for creating a an account on our website. Please verify your email address by clicking the button below</p>
				<a href="` + domName + `/api/account/emailver/` + account.Username + `/` + emailVerPassword + `" 			 	style="box-sizing:border-box;text-decoration:none;background-color:#007bff;border:solid 1px #007bff;border-radius:4px;color:#ffffff;font-size:16px;font-weight:bold;margin:0;padding:9px 25px;display:inline-block;letter-spacing:1px">
        Verify email address
    		</a>
			</html>
		`
	err := send(subject, HTMLbody, account)
	if err != nil {
		return err
	}

	return nil
}

func send(subject, HTMLbody string, account model.Account) error {
	// sender data
	to := []string{account.Email}
	// smtp - Simple Mail Transfer Protocol
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	// Set up authentication information
	// TODO Stefan: Email sending not workin when the email and password are taken from the .env file
	//senderMail := os.Getenv("SENDER_EMAIL")
	//senderAppPassword := os.Getenv("SENDER_APP_PASSWORD")

	auth := smtp.PlainAuth("", "xmlprojekat1@gmail.com", "cmyrhmgnrenxhwuc", host)
	msg := []byte(
		"From: " + "FTN Airlines" + ": <" + "xmlprojekat1@gmail.com" + ">\r\n" +
			"To: " + account.Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME: MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			HTMLbody)
	err := smtp.SendMail(address, auth, "xmlprojekat1@gmail.com", to, msg)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
