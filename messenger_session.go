package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/paked/messenger"
)

// MessengerSession - contains the messenger specific session data
type MessengerSession struct {
	store    SessionStore
	session  *Session
	response *messenger.Response
	sender   messenger.Sender
	client   *messenger.Messenger
	bot      *Bot
}

// NewMessengerSession - create a new messenger session instance
func NewMessengerSession(b *Bot, s SessionStore, w *messenger.Response, u messenger.Sender) *MessengerSession {
	return &MessengerSession{
		store:    s,
		session:  s.Acquire(fmt.Sprintf("messenger/%d", u.ID)),
		response: w,
		sender:   u,
		bot:      b,
	}
}

// MatchTemplate - find the macthed template and sends its content
func (m *MessengerSession) MatchTemplate(txt string) (error, bool) {
	var matched *goquery.Selection
	m.bot.Document.Find("template").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		matches := s.AttrOr("matches", "")
		rgx, err := regexp.Compile(matches)
		if err != nil {
			return false
		}
		if rgx.MatchString(txt) {
			picked := RandInt(int64(s.Children().Length()))
			s.Children().EachWithBreak(func(i int, s *goquery.Selection) bool {
				if int64(i) == picked {
					matched = s
					return false
				}
				return true
			})
			return false
		}
		return true
	})
	if matched != nil {
		switch goquery.NodeName(matched) {
		case "div", "dialog", "form":
			return m.SendNodesOf(matched), true
		case "menu", "nav":
			return m.SendBasicMenu(m.bot.Menus[matched.AttrOr("id", "")]), true
		}
	}
	return nil, false
}

// SendText - send some text to the user
func (m *MessengerSession) SendText(text string) error {
	return m.response.Text(text)
}

// SendMedia - send a media data to the user
func (m *MessengerSession) SendMedia(typ, src string) error {
	return m.response.Attachment(messenger.AttachmentType(typ), src)
}

// SendInput - send an input to t
func (m *MessengerSession) SendInput(input *Input) error {
	if input == nil {
		return nil
	}
	if input.Type == "options" {
		callToActions := []messenger.StructuredMessageButton{}
		for _, btn := range input.Options {
			cta := messenger.StructuredMessageButton{}
			cta.Title = btn.Title
			cta.Type = "postback"
			cta.Payload = btn.Href
			callToActions = append(callToActions, cta)
		}
		return m.response.ButtonTemplate("Select an option ...", &callToActions)
	}
	return nil
}

// SendBasicMenu - send a menu to the user
func (m *MessengerSession) SendBasicMenu(menu *Menu) error {
	if menu == nil {
		return nil
	}
	callToActions := []messenger.StructuredMessageButton{}
	for _, btn := range menu.Buttons {
		cta := messenger.StructuredMessageButton{}
		cta.Title = btn.Title
		if strings.HasPrefix(btn.Href, "http://") || strings.HasPrefix(btn.Href, "https://") {
			cta.Type = "web_url"
			cta.MessengerExtensions = (btn.Embed != "no") && (btn.Embed != "")
			cta.WebviewHeightRatio = btn.Embed
		} else {
			cta.Type = "postback"
			cta.Payload = fmt.Sprintf("trgt=%s&btn=%s", btn.Href[1:], btn.ID)
		}
		callToActions = append(callToActions, cta)
	}
	return m.response.ButtonTemplate(menu.Title, &callToActions)
}

// SendNodesOf - send the child nodes of the specified node
func (m *MessengerSession) SendNodesOf(s *goquery.Selection) error {
	var err error

	m.session.CurrentContainer = s.AttrOr("id", "")
	parentId := s.AttrOr("id", "")

	// elements loop ...
	s.Children().EachWithBreak(func(i int, s *goquery.Selection) bool {
		next := true
		m.session.ExpectingUserInput = false
		m.session.CurrentElement = s.AttrOr("id", fmt.Sprintf("node-%s-element-%d", parentId, i))
		s.SetAttr("id", m.session.CurrentElement)
		if strings.Contains(m.session.ElementsHistory, fmt.Sprintf(";%d;", i)) {
			return true
		}
		m.session.ElementsHistory += fmt.Sprintf(";%d;", i)
		ifexpr := s.AttrOr("if", "")
		if (ifexpr != "" && EvalIfExpression(ifexpr, m.session.Data)) || ifexpr == "" {
			switch strings.ToLower(goquery.NodeName(s)) {
			case "media", "embed", "resource":
				err = m.SendMedia(s.AttrOr("type", "image"), s.AttrOr("src", ""))
			case "text", "label", "p", "span", "inline", "line":
				err = m.SendText(GetExpressionValue(s.Text(), m.session.Data))
			case "input", "select":
				err = m.SendInput(m.bot.Inputs[s.AttrOr("id", "")])
				m.session.ExpectingUserInput = true
				next = false
			case "menu", "nav":
				err = m.SendBasicMenu(m.bot.Menus[s.AttrOr("id", "")])
			default:
				log.Println("[unkown]", "un-implemented HTML tag")
			}
		}
		return next
	})

	// reached the end !!
	if s.Children().Last().AttrOr("id", "") == m.session.CurrentElement && !m.session.ExpectingUserInput {
		m.session.EOF = true
		if action := s.AttrOr("action", ""); action != "" {
			err = m.Submit(action)
		}
		m.session.Clear()
	}

	return err
}

// IsExpectingUserInput - whether it is expecting a user input or not
func (m *MessengerSession) IsExpectingUserInput() bool {
	return m.session.ExpectingUserInput
}

// SetData - set session data
func (m *MessengerSession) SetData(k, v string) {
	m.session.Data[k] = v
}

// GetData - get session data
func (m *MessengerSession) GetData(k string) interface{} {
	return m.session.Data[k]
}

// GetAllData - get all session data
func (m *MessengerSession) GetAllData() map[string]interface{} {
	return m.session.Data
}

// SetCurrentContainer - set the current session container
func (m *MessengerSession) SetCurrentContainer(id string) {
	m.session.CurrentContainer = id
}

// GetCurrentContainer - get the current session container
func (m *MessengerSession) GetCurrentContainer() string {
	return m.session.CurrentContainer
}

// SetCurrentElement - set the current element
func (m *MessengerSession) SetCurrentElement(id string) {
	m.session.CurrentElement = id
}

// GetCurrentElement - get the current element
func (m *MessengerSession) GetCurrentElement() string {
	return m.session.CurrentElement
}

// IsEOF - reached the end of the form?
func (m *MessengerSession) IsEOF() bool {
	return m.session.EOF
}

// Submit - submits the current session data
func (m *MessengerSession) Submit(target string) error {
	vals := url.Values{}
	for k, v := range m.session.Data {
		vals.Set(k, v.(string))
	}
	vals.Set("fb_profile_id", fmt.Sprintf("%d", m.sender.ID))
	resp, err := http.PostForm(target, vals)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// Forget - forget the current user session
func (m *MessengerSession) Forget() {
	m.store.Forget(fmt.Sprintf("messenger/%d", m.sender.ID))
	m.session = m.store.Acquire(fmt.Sprintf("messenger/%d", m.sender.ID))
}
