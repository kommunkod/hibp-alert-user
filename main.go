package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"slices"
	"strings"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Get information on the latest breach
	latestBreach, err := QueryLatestBreach()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if latestBreach.IsFabricated && config.Ignore.Fabricated {
		fmt.Println("Latest breach is fabricated and should be ignored. Exiting...")
		// TODO: Add breach to config
		os.Exit(0)
	}

	if latestBreach.IsSensitive && config.Ignore.Sensitive {
		fmt.Println("Latest breach is sensitive and should be ignored. Exiting...")
		os.Exit(0)
	}

	if latestBreach.IsSpamList && config.Ignore.SpamList {
		fmt.Println("Latest breach is a spam list and should be ignored. Exiting...")
		os.Exit(0)
	}

	if latestBreach.IsMalware && config.Ignore.Malware {
		fmt.Println("Latest breach is malware and should be ignored. Exiting...")
		os.Exit(0)
	}

	if latestBreach.IsRetired && config.Ignore.Retired {
		fmt.Println("Latest breach is retired and should be ignored. Exiting...")
		os.Exit(0)
	}

	fmt.Printf("Latest Breach Name: %s\n", latestBreach.Name)

	if latestBreach.Name == config.LatestBreach {
		fmt.Println("No new breaches found! Exiting...")
		os.Exit(0)
	}

	fmt.Println("New breach found!")

	notifyUsers := map[string][]string{}
	userPresentIn := map[string][]string{}

	uniqueBreaches := map[string]bool{}

	for _, domain := range config.Domains {
		fmt.Printf("Querying the API for %s...\n", domain)
		breaches, err := QueryDomainBreaches(domain, config)
		if err != nil {
			fmt.Printf("Error querying the API for %s: %s\n", domain, err)
			continue
		}

		for user, presentIn := range breaches {
			username := user + "@" + domain
			for _, breach := range presentIn {
				uniqueBreaches[breach] = true

				if _, ok := userPresentIn[username]; !ok {
					userPresentIn[username] = []string{}
				}

				userPresentIn[username] = append(userPresentIn[username], breach)

				if !slices.Contains(config.NotifiedBreaches, breach) {
					if _, ok := notifyUsers[username]; !ok {
						notifyUsers[username] = []string{}
					}

					notifyUsers[username] = append(notifyUsers[username], breach)
				}
			}
		}
	}

	fmt.Printf("%d users are going to be notified\n", len(notifyUsers))

	html, err := os.ReadFile("email_template.html")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	emailBody := string(html)

	emailBody = strings.ReplaceAll(emailBody, "{{bg_color}}", config.Email.Colors.Background)
	emailBody = strings.ReplaceAll(emailBody, "{{text_color}}", config.Email.Colors.Text)
	emailBody = strings.ReplaceAll(emailBody, "{{body_header}}", config.Email.Body.Header)

	bodyText := ""
	for _, text := range config.Email.Body.Texts {
		bodyText += fmt.Sprintf("<p style=\"line-height: 1.4em;\">%s</p>", text)
	}

	emailBody = strings.ReplaceAll(emailBody, "{{body_text}}", bodyText)

	previousBreachesText := ""

	for _, text := range config.Email.Body.PreviousBreachTexts {
		previousBreachesText += fmt.Sprintf("<p style=\"line-height: 1.4em;\">%s</p>", text)
	}

	emailBody = strings.ReplaceAll(emailBody, "{{previous_breaches_text}}", previousBreachesText)

	for user, breaches := range notifyUsers {
		if user != "lars@scheibling.se" {
			continue
		}
		body := string(emailBody)

		newBreachRows := ""
		for _, breach := range breaches {
			newBreachRows += fmt.Sprintf(`<tr>
				<td>%s</td>
				<td><a href="https://haveibeenpwned.com/PwnedWebsites#%s">More Information</a></td>
			</tr>`, breach, breach)
		}

		body = strings.ReplaceAll(body, "{{new_breach_rows}}", newBreachRows)

		previousBreachesRows := ""
		for _, breach := range userPresentIn[user] {
			previousBreachesRows += fmt.Sprintf(`<tr>
				<td>%s</td>
				<td><a href="https://haveibeenpwned.com/PwnedWebsites#%s">More Information</a></td>
			</tr>`, breach, breach)
		}

		body = strings.ReplaceAll(body, "{{previous_breaches_rows}}", previousBreachesRows)

		body = strings.ReplaceAll(body, "{{username}}", user)

		fullBody := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s", config.Smtp.Sender, user, config.Email.Subject, body)

		err = smtp.SendMail(
			fmt.Sprintf("%s:%d", config.Smtp.Host, config.Smtp.Port),
			smtp.PlainAuth("", config.Smtp.User, config.Smtp.Pass, config.Smtp.Host),
			config.Smtp.Sender,
			[]string{user},
			[]byte(fullBody),
		)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Hold")
	}

	for _, breach := range config.NotifiedBreaches {
		uniqueBreaches[breach] = true
	}

	uniqueBreaches[latestBreach.Name] = true

	config.NotifiedBreaches = []string{}

	for breach := range uniqueBreaches {
		config.NotifiedBreaches = append(config.NotifiedBreaches, breach)
	}

	config.LatestBreach = latestBreach.Name

	err = SaveConfig(config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)

}
