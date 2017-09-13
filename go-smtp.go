package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/go-gomail/gomail"
	"bytes"
	"io"
)

func main() {

	//sendMail()
	http.HandleFunc("/v1/email", sendEmail)
	http.ListenAndServe(":9000", nil)

}

type Emailbody struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Msg     string   `json:"msg"`
}

type FileAttachment struct {
	Filename    string
	FileContent []byte
}

func sendEmail(x http.ResponseWriter, r *http.Request) {
	x.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		fmt.Printf("ReceivedPOST")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {

			fmt.Printf("\nBadboy")
		}

		result := bytes.Split(body, [] byte("filename="))
		//Message (to, subject, message JSON)
		var fileAttachment []FileAttachment
		main_message := result[0]
		index_message := bytes.Index(main_message, []byte("{"))
		final_message := main_message[index_message: ]
		email_message := bytes.Split(final_message, [] byte("------WebKit"))
		fmt.Printf("FinalMessage\n%s", email_message[0])
		emailbody := new(Emailbody)
		if err := json.Unmarshal(email_message[0], &emailbody); err != nil {
			fmt.Printf("\nBadMarshall")
		}

		//Decode Files from POST

		len_divisions := len(result)
		var tempfileAttachment FileAttachment

		i := 1
		for (i < len_divisions) {
			main_content := bytes.Split(result[i], [] byte("Content-"))
			filename := string(bytes.Trim(main_content[0], "\"\r\n"))
			index := bytes.Index(main_content[1], [] byte("\r\n"))
			data_to_be_iterated := main_content[1]
			split_content := data_to_be_iterated[index + 4: ]
			file_content := bytes.Split(split_content, [] byte("\r\n------WebKit"))
			file_data := file_content[0] //This is where the data lies
			tempfileAttachment.Filename = filename
			tempfileAttachment.FileContent = file_data
			fileAttachment = append(fileAttachment, tempfileAttachment)
			
      //File Writer <change to buffer writer if needed>
			/*err = ioutil.WriteFile(filename, file_data, 0666)
			if err != nil {
				panic(err)
			}*/
			i++
		}

		m := gomail.NewMessage()
		recipients := emailbody.To
		print(len(recipients))
		addresses := make([]string, len(recipients))
		for i, recipient := range recipients {
			addresses[i] = m.FormatAddress(recipient, "")
		}

		m.SetHeader("From", "noreply@yourdomain.com")
		m.SetHeader("To", addresses...)
		m.SetHeader("Subject", emailbody.Subject)
		m.SetBody("text/html", emailbody.Msg)

		no_of_attachments := len(fileAttachment)
		
		if no_of_attachments == 1 {

			m.Attach(fileAttachment[0].Filename, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(fileAttachment[0].FileContent)
				return err
			}))
		}
		if no_of_attachments == 2 {

			m.Attach(fileAttachment[0].Filename, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(fileAttachment[0].FileContent)
				return err
			}))
			m.Attach(fileAttachment[1].Filename, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(fileAttachment[1].FileContent)
				return err
			}))
		}
		if no_of_attachments == 3 {

			m.Attach(fileAttachment[0].Filename, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(fileAttachment[0].FileContent)
				return err
			}))
			m.Attach(fileAttachment[1].Filename, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(fileAttachment[1].FileContent)
				return err
			}))
			m.Attach(fileAttachment[2].Filename, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(fileAttachment[2].FileContent)
				return err
			}))
		}

		d := gomail.NewDialer(SMTP_URL, 25, USERNAME, PASSWORD)
		if err := d.DialAndSend(m); err != nil {
			fmt.Println(err)
		}
	}
}

const (
	USERNAME = "smtp_usrname"
	PASSWORD = "smtp_password"
	SMTP_URL = "smtp.server.com"
)
