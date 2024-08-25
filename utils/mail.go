package utils

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	// api "github.com/twilio/twilio-go/rest/api/v2010"
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

	otpLength := 5 // Length of the OTP
	min := 10000   // Minimum value of the OTP (inclusive)
	max := 99999   // Maximum value of the OTP (inclusive)

	otp := strconv.Itoa(rand.Intn(max-min+1) + min) // Generate a random number within the specified range

	if len(otp) < otpLength {
		otp = fmt.Sprintf("%0*s", otpLength, otp) // Pad the OTP with leading zeros if necessary
	}

	return otp
}

func SendOTPMessage(phoneNumber, otp string) error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	rapidAPIKey := os.Getenv("RAPID_API_KEY")
	token := os.Getenv("ACCESS_TOKEN")
	url := fmt.Sprintf("https://smsapi-com3.p.rapidapi.com/sms.do?access_token=%s", token)

	// Prepare the payload
	payload := fmt.Sprintf(`{
        "to": "%s",
        "message": "Your OTP code is: %s",
        "from": "Aluta market",
        "normalize": "",
        "group": "",
        "encoding": "",
        "flash": "",
        "test": "",
        "details": "",
        "date": "",
        "date_validate": "",
        "time_restriction": "follow",
        "allow_duplicates": "",
        "idx": "",
        "check_idx": "",
        "max_parts": "",
        "fast": "",
        "notify_url": "",
        "format": "json"
    }`, phoneNumber, otp)

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return err
	}

	// Set the required headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-RapidAPI-Key", rapidAPIKey)
	req.Header.Add("X-RapidAPI-Host", "smsapi-com3.p.rapidapi.com")

	// Send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(res)
	fmt.Println(string(body))

	return nil
}
