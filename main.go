package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-co-op/gocron"
	gomail "gopkg.in/mail.v2"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	sendEmail("Hola")
	runCronJobs()
}

func runCronJobs() {
	dsn := flag.String("dsn", fmt.Sprintf("%v:@/%v?parseTime=true", db_username, db_name), "MySQL data source name")
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minutes().Do(func() {
		db, err := openDB(*dsn)
		if err != nil {
			sendEmail(err.Error())
			return
		}

		db.Close()
	})

	s.StartBlocking()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func sendEmail(err string) {
	m := gomail.NewMessage()

	m.SetHeader("From", "diego@andina.dev")
	m.SetHeader("To", "andrianodna@gmail.com")
	m.SetHeader("Subject", "Error en la db!")
	m.SetBody("text/plain", fmt.Sprintf("Error al enviar al pingear db, %v", err))

	d := gomail.NewDialer("smtp.zoho.com", 465, "diego@andina.dev", fmt.Sprintf("%s", email_password))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

}
