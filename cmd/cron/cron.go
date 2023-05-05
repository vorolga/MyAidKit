package main

import (
	"log"
	"main/internal/composites"
	"main/internal/models"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"
	"time"

	cron "github.com/robfig/cron/v3"
)

func main() {
	postgresDBC, err := composites.NewPostgresDBComposite()
	if err != nil {
		log.Fatal("postgres db composite failed")
	}

	loc := time.UTC
	scheduler := cron.New(cron.WithLocation(loc))

	defer scheduler.Stop()

	scheduler.AddFunc("30 2 * * *", func() { DeleteNotifications(postgresDBC) })
	scheduler.AddFunc("0 * * * *", func() { SendNotifications(postgresDBC) })

	go scheduler.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func DeleteNotifications(postgresDBC *composites.PostgresDBComposite) {
	sqlScript := "DELETE FROM notification_user WHERE time::timestamptz < now()::timestamptz"
	_, err := postgresDBC.DB.Exec(sqlScript)
	if err != nil {
		log.Fatal(err)
	}

	loc := time.UTC
	log.Println(time.Now().In(loc).Format("2006-01-02 15:04:05") + " DeleteNotifications\n")
}

func SendNotifications(postgresDBC *composites.PostgresDBComposite) {
	loc := time.UTC
	currentTime := time.Now().In(loc).Format("2006-01-02 15:04")

	from := "vorrovvorrov@gmail.com"
	password := os.Getenv("EMAILPASSWORD")

	host := "smtp.gmail.com"
	port := "587"

	sqlScript := "select id_from, to_is_user, id_to_user, name_to, name_medicine, email, id_family from users join notification_user on id_from = users.id where time::timestamptz = to_timestamp($1, 'YYYY-MM-DD HH24:MI')::timestamptz;"
	rows, err := postgresDBC.DB.Query(sqlScript, currentTime)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		notificationFrom := models.NotificationsFrom{
			IDFrom:       0,
			ToIsUser:     false,
			IDTo:         0,
			NameTo:       "",
			NameMedicine: "",
			Email:        "",
			IDFamily:     0,
		}

		if err = rows.Scan(&notificationFrom.IDFrom, &notificationFrom.ToIsUser, &notificationFrom.IDTo,
			&notificationFrom.NameTo, &notificationFrom.NameMedicine, &notificationFrom.Email, &notificationFrom.IDFamily); err != nil {
			log.Fatal(err)
		}

		if notificationFrom.IDFrom == notificationFrom.IDTo && notificationFrom.ToIsUser {
			toList := []string{notificationFrom.Email}
			msg := "Вам неообходимо принять " + notificationFrom.NameMedicine + "\r\n" +
				"Зайдите в приложение и отметьте количество выпитых таблеток или просто нажмите \"Принять\", если выпитое лекарство не является таблеткой.\r\n" +
				"https://myaidkit.ru"

			body := []byte(msg)

			authSMTP := smtp.PlainAuth("", from, password, host)
			err = smtp.SendMail(host+":"+port, authSMTP, from, toList, body)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			if !notificationFrom.ToIsUser {
				sqlScript = "select email from users where id_family = $1 and users.is_adult = true;"
				userRows, err := postgresDBC.DB.Query(sqlScript, notificationFrom.IDFamily)
				if err != nil {
					log.Fatal(err)
				}
				func() {
					defer userRows.Close()

					for userRows.Next() {
						var email string
						if err = userRows.Scan(&email); err != nil {
							log.Fatal(err)
						}

						toList := []string{email}
						msg := "Члену семьи " + notificationFrom.NameTo + " необходимо принять " + notificationFrom.NameMedicine + "\r\n" +
							"Зайдите в приложение и отметьте количество выпитых таблеток или просто нажмите \"Принять\", если выпитое лекарство не является таблеткой.\r\n" +
							"https://myaidkit.ru"
						body := []byte(msg)

						authSMTP := smtp.PlainAuth("", from, password, host)
						err = smtp.SendMail(host+":"+port, authSMTP, from, toList, body)
						if err != nil {
							log.Fatal(err)
						}
					}
				}()
			} else {
				sqlScript = "select id, email from users where id_family = $1 and (users.is_adult = true or id = $2);"
				userRows, err := postgresDBC.DB.Query(sqlScript, notificationFrom.IDFamily, notificationFrom.IDTo)
				if err != nil {
					log.Fatal(err)
				}
				func() {
					defer userRows.Close()

					for userRows.Next() {
						var email string
						var id int64
						if err = userRows.Scan(&id, &email); err != nil {
							log.Fatal(err)
						}

						if id == notificationFrom.IDTo {
							toList := []string{email}
							msg := "Вам неообходимо принять " + notificationFrom.NameMedicine + "\r\n" +
								"Зайдите в приложение и отметьте количество выпитых таблеток или просто нажмите \"Принять\", если выпитое лекарство не является таблеткой.\r\n" +
								"https://myaidkit.ru"

							body := []byte(msg)

							authSMTP := smtp.PlainAuth("", from, password, host)
							err = smtp.SendMail(host+":"+port, authSMTP, from, toList, body)
							if err != nil {
								log.Fatal(err)
							}

							continue
						}

						toList := []string{email}
						msg := "Члену семьи " + notificationFrom.NameTo + " необходимо принять " + notificationFrom.NameMedicine
						body := []byte(msg)

						authSMTP := smtp.PlainAuth("", from, password, host)
						err = smtp.SendMail(host+":"+port, authSMTP, from, toList, body)
						if err != nil {
							log.Fatal(err)
						}
					}
				}()
			}
		}
	}

	log.Println(time.Now().In(loc).Format("2006-01-02 15:04:05") + " SendNotifications\n")
}
