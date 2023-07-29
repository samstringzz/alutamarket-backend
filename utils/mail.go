package utils

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendMail(reciever string) {
	// Set up SendGrid client
	apiKey := os.Getenv("SENDGRID_KEY")
	client := sendgrid.NewSendClient(apiKey)

	// Compose the email message
	from := mail.NewEmail("Sender Name", os.Getenv("SENDER_EMAIL"))
	to := mail.NewEmail("Recipient Name", reciever)
	subject := "Test Email"
	content := mail.NewContent("text/plain", "This is a test email sent using SendGrid.")
	message := mail.NewV3MailInit(from, subject, to, content)

	// Send the email
	response, err := client.Send(message)
	if err != nil {
		log.Fatal("Error sending email:", err)
	}

	// Check the response status
	fmt.Println("Email sent. Status code:", response.StatusCode)
	fmt.Println(response)
}

func GenerateOTP() string {
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed)) // Seed the random number generator with the current time

	otpLength := 6 // Length of the OTP
	min := 100000  // Minimum value of the OTP (inclusive)
	max := 999999  // Maximum value of the OTP (inclusive)

	otp := strconv.Itoa(rand.Intn(max-min+1) + min) // Generate a random number within the specified range

	if len(otp) < otpLength {
		otp = fmt.Sprintf("%0*s", otpLength, otp) // Pad the OTP with leading zeros if necessary
	}

	return otp
}

func SendOtpMessage(otp string, reciever string) (string, error) {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	os.Setenv("TWILIO_ACCOUNT_SID", accountSID)
	os.Setenv("TWILIO_AUTH_TOKEN", authToken)

	client := twilio.NewRestClient()
	params := &api.CreateMessageParams{
		Body: new(string),
		From: new(string),
		To:   &reciever,
	}

	*params.Body = "This is the OTP for your registration: " + otp + " (expires in 5 minutes)"
	*params.From = "+16183614700"

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	if resp.Sid != nil {
		fmt.Println(*resp.Sid)
		return *resp.Sid, nil
	} else {
		fmt.Println(resp.Sid)
		return "", errors.New("failed to send OTP message")
	}
}
