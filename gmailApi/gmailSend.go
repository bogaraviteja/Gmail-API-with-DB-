package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gmail/database"
	"gmail/gmailService"

	models "gmail/models"
	"io/ioutil"
	"log"
	"net/http"

	gmail "google.golang.org/api/gmail/v1"
)

var db = database.Db()

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func chunkSplit(body string, limit int, end string) string {

	var charSlice []rune

	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	var result string = ""

	for len(charSlice) >= 1 {

		result = result + string(charSlice[:limit]) + end

		charSlice = charSlice[limit:]

		if len(charSlice) < limit {
			limit = len(charSlice)
		}
	}
	return result
}

func createMessageWithAttachment(from string, to string, subject string, content string, fileDir string, fileName string, Signature string) *gmail.Message {

	fileBytes, err := ioutil.ReadFile(fileDir + fileName)
	if err != nil {
		log.Fatalf("Unable to read file for attachment: %v", err)
	}

	fileMIMEType := http.DetectContentType(fileBytes)

	fileData := base64.StdEncoding.EncodeToString(fileBytes)

	boundary := randStr(32, "alphanum")

	messageBody := []byte("Content-Type: multipart/mixed; boundary=" + boundary + " \n" +
		"MIME-Version: 1.0\n" +
		"to: " + to + "\n" +
		"from: " + from + "\n" +
		"subject: " + subject + "\n\n" +

		"--" + boundary + "\n" +
		"Content-Type: text/plain; charset=" + string('"') + "UTF-8" + string('"') + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: 7bit\n\n" +
		content + "\n\n" +
		"--" + boundary + "\n" +

		"Content-Type: " + fileMIMEType + "; name=" + string('"') + fileName + string('"') + " \n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: base64\n" +
		"Content-Disposition: attachment; filename=" + string('"') + fileName + string('"') + " \n\n" +
		chunkSplit(fileData, 76, "\n") +
		"--" + boundary + "--")

	raw := base64.URLEncoding.EncodeToString(messageBody)

	return &gmail.Message{Raw: raw}
}

func main() {

	// db.Create(&models.Person{Name: "Raviteja", Gender: "M", Email: "bogaraviteja@gmail.com", Address: "Hyderabad", Pincode: 500083})

	srv, err := gmailService.Service()

	msgContent := `Hello!
                    This is a automated mail from Gmail API, 
					please don't reply!  
                    Good Bye!`

	subject := "Email WITH ATTACHMENT from GMail API"

	signature := `Name : Raviteja Boga
				  phone : 0123456789
				  Location : Hyderabad`

	email, err := db.Model(&models.Person{}).Select("name,email").Rows()
	if err != nil {
		panic(err)
	}
	defer email.Close()
	for email.Next() {
		var e, name string
		email.Scan(&name, &e)
		fmt.Println(e, name)

		messageWithAttachment := createMessageWithAttachment("ravitejaboga1998@gmail.com", e, subject, msgContent, "C:/Users/DELL/OneDrive/Desktop/FANGATE/candence/", "Pakistan_April2014.pdf", signature)

		m, err := srv.Users.Messages.Send("me", messageWithAttachment).Do()
		if err != nil {
			panic(err.Error())
			return
		}

		fmt.Println(m.ServerResponse.Header)

		db.Create(&models.SentEmails{Name: name, Email: e, Subject: subject, Content: msgContent, Signature: signature})
	}
}
