package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// RE implement the session store & session
import (
	"github.com/PuerkitoBio/goquery"
	"github.com/paked/messenger"
)

type Messenger struct {
	bot    *Bot
	client *messenger.Messenger
	store  SessionStore
}

// NewMessenger - Initialize new Messenger wrapper
func NewMessenger(b *Bot) *Messenger {
	m := &Messenger{
		bot: b,
		client: messenger.New(messenger.Options{
			Verify:      b.Configs["messenger-verify"] == "yes",
			AppSecret:   b.Configs["messenger-app-secret"],
			VerifyToken: b.Configs["messenger-verify-token"],
			Token:       b.Configs["messenger-page-token"],
		}),
		store: NewSessionStore(),
	}
	m.Boot()
	return m
}

// SendGreeting - send the greeting data/messages
func (m Messenger) SendGreeting() error {
	if desc, ok := m.bot.Configs["description"]; ok {
		return m.client.GreetingSetting(desc)
	}
	return nil
}

// GetMainMenu - find the main menu
func (m Messenger) GetMainMenu() *Menu {
	mainMenu := m.bot.Configs["main-menu"]
	if "" == mainMenu {
		return nil
	}
	menu := m.bot.Menus[mainMenu]
	return menu
}

// SendMainMenu - sending the main messenger menu
func (m Messenger) SendMainMenu(menu *Menu) error {
	if menu == nil {
		return nil
	}
	callToActions := []messenger.CallToActionsItem{}
	for _, btn := range menu.Buttons {
		cta := messenger.CallToActionsItem{}
		cta.Title = btn.Title
		if strings.HasPrefix(btn.Href, "http://") || strings.HasPrefix(btn.Href, "https://") {
			cta.Type = "web_url"
			cta.MessengerExtension = (btn.Embed != "no") && (btn.Embed != "")
			cta.WebviewHeightRatio = btn.Embed
		} else {
			cta.Type = "postback"
			cta.Payload = fmt.Sprintf("trgt=%s&btn=%s", btn.Href[1:], btn.ID)
		}
		callToActions = append(callToActions, cta)
	}
	return m.client.CallToActionsSetting("existing_thread", callToActions)
}

// Boot - Setup the default handlers
func (m Messenger) Boot() {
	// send the main menu
	if err := m.SendMainMenu(m.GetMainMenu()); err != nil {
		log.Println(err)
		return
	}

	// send greetings
	if err := m.SendGreeting(); err != nil {
		log.Println(err)
		return
	}

	// handling the free text input
	m.client.HandleMessage(func(msg messenger.Message, r *messenger.Response) {
		// setup the session
		session := NewMessengerSession(
			m.bot,
			m.store,
			r,
			msg.Sender,
		)

		// setting up states "seen, typing ... etc"
		r.SenderAction("mark_seen")
		r.SenderAction("typing_on")
		defer r.SenderAction("typing_off")

		// expecting a user input ...
		if session.IsExpectingUserInput() {
			input := m.bot.Inputs[session.GetCurrentElement()]
			ans := msg.Text
			if input.Type == "file" {
				if len(msg.Attachments) < 1 {
					r.Text("Please upload a valid file", messenger.MessagingType("RESPONSE"))
					return
				}
				ans = msg.Attachments[0].Payload.URL
			}
			session.SetData(input.Name, ans)
			session.SendNodesOf(m.bot.Document.Find("#" + session.GetCurrentContainer()))
			return
		}

		// from templates
		err, found := session.MatchTemplate(msg.Text)
		if !found || err != nil {
			log.Println("[template]", err)
			session.SendText("Sorry, but I cannot understand your input")
		}
	})

	// handling the postbacks "link/button" clicks
	m.client.HandlePostBack(func(p messenger.PostBack, r *messenger.Response) {
		// setup the session
		session := NewMessengerSession(
			m.bot,
			m.store,
			r,
			p.Sender,
		)

		// setting up states "seen, typing ... etc"
		r.SenderAction("mark_seen")
		r.SenderAction("typing_on")
		defer r.SenderAction("typing_off")

		// get started button ?
		if p.Payload == "get_started" {
			if err := m.SendGreeting(); err != nil {
				log.Println("[greeting]", err)
				return
			}
			if err := session.SendBasicMenu(m.GetMainMenu()); err != nil {
				log.Println("[menu]", err)
				return
			}
			return
		}

		// the required action
		needle, _ := url.ParseQuery(p.Payload)
		log.Println("[postback]", needle, ", ShouldBeInput: ", session.IsExpectingUserInput())

		// reset?
		if needle.Get("btn") != "" && m.bot.Buttons[needle.Get("btn")] != nil && m.bot.Buttons[needle.Get("btn")].Reset {
			session.Forget()
		}

		// an answer?
		if session.IsExpectingUserInput() && needle.Get("answer") == "yes" && needle.Get("input") != "" {
			input := m.bot.Inputs[needle.Get("input")]
			if input == nil {
				log.Println("[input]", "unknown input with id (", needle.Get("input"), ")")
				return
			}
			session.SetData(input.Name, needle.Get("value"))
			session.SendNodesOf(m.bot.Document.Find("#" + session.GetCurrentContainer()))
			return
		}

		// an action?
		trgt := m.bot.Document.Find("#" + needle.Get("trgt"))
		switch goquery.NodeName(trgt) {
		case "div", "dialog", "form":
			log.Println("[sending nodes] has error?", session.SendNodesOf(trgt))
		case "menu", "nav":
			log.Println("[sending menu] has error?", session.SendBasicMenu(m.bot.Menus[trgt.AttrOr("id", "")]))
		default:
			r.Text("Sorry, I couldn't understand you", messenger.MessagingType("RESPONSE"))
			log.Println("[unknown]", "undefined node type", needle.Get("trgt"))
		}
	})
}

// ServeHTTP - handle the http requests
func (m Messenger) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.client.Handler().ServeHTTP(res, req)
}
