package main

import (
	"fmt"
	"io"
	"strings"
)

import (
	"github.com/PuerkitoBio/goquery"
)

// parse the specified io.Reader
func NewBotFromReader(r io.Reader) (*Bot, error) {
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return nil, err
	}

	// initialize the bot container
	bot := new(Bot)
	bot.Document = doc
	bot.Configs = map[string]string{}
	bot.Buttons = map[string]*Button{}
	bot.Menus = map[string]*Menu{}
	bot.Dialogs = map[string]*Dialog{}
	bot.Inputs = map[string]*Input{}

	// compile the head.meta to configs hashmap
	doc.Find("meta").Each(func(_ int, s *goquery.Selection) {
		key, _ := s.Attr("name")
		value, _ := s.Attr("content")

		bot.Configs[key] = value
	})

	// populate the navs
	doc.Find("menu,nav").Each(func(i int, m *goquery.Selection) {
		menu := &Menu{}
		menu.ID = m.AttrOr("id", fmt.Sprintf("menu%d", i))
		menu.Title = m.AttrOr("title", "?")
		menu.Buttons = []*Button{}
		m.Find("a,button").Each(func(i int, b *goquery.Selection) {
			btn := &Button{}
			btn.ID = b.AttrOr("id", fmt.Sprintf("button%d", len(bot.Buttons)+1))
			btn.Title = b.Text()
			btn.Href = b.AttrOr("href", "")
			btn.Embed = b.AttrOr("embed", "false")
			btn.Reset = (b.AttrOr("reset", "false") == "true")

			menu.Buttons = append(menu.Buttons, btn)
			bot.Buttons[btn.ID] = btn

			b.SetAttr("id", btn.ID)
		})
		m.SetAttr("id", menu.ID)
		bot.Menus[menu.ID] = menu
	})

	// populate forms
	doc.Find("dialog,form").Each(func(i int, d *goquery.Selection) {
		dialog := &Dialog{}
		dialog.Node = d
		dialog.ID = d.AttrOr("id", fmt.Sprintf("dialog%d", i))
		dialog.Action = d.AttrOr("action", "")
		dialog.Help = d.AttrOr("help", "")
		fn := func(i int, n *goquery.Selection) {
			input := &Input{}
			input.NS = dialog.ID
			input.ID = n.AttrOr("id", fmt.Sprintf("input%d", len(bot.Inputs)+1))
			input.Name = n.AttrOr("name", input.ID)
			input.Value = n.AttrOr("value", "")
			input.If = n.AttrOr("if", "")
			input.Options = []*Button{}
			input.Type = n.AttrOr("type", "text")
			n.Find("option").Each(func(i int, o *goquery.Selection) {
				btn := &Button{}
				btn.ID = fmt.Sprintf("input-%s-%d", input.ID, len(bot.Inputs)+1)
				btn.Href = fmt.Sprintf("input=%s&answer=yes&value=%s", input.ID, o.AttrOr("value", ""))
				btn.Embed = "no"
				btn.Title = o.Text()
				btn.Reset = false
				input.Options = append(input.Options, btn)
				o.SetAttr("id", btn.ID)
			})
			bot.Inputs[input.ID] = input
			n.SetAttr("id", input.ID)
		}
		d.Find("input,select").Each(func(i int, n *goquery.Selection) {
			switch strings.ToLower(goquery.NodeName(n)) {
			case "input":
				fn(i, n)
			case "select":
				n.SetAttr("type", "options")
				fn(i, n)
			}
		})
		// fmt.Println(d.Find("div").Children().Get(1))
		d.Find("div").Children().Each(func(i int, s *goquery.Selection) {
			s.SetAttr("if", s.Parent().AttrOr("if", ""))
			s.Parent().AfterSelection(s.Clone())
			s.Remove()
		})
		bot.Dialogs[dialog.ID] = dialog
		d.SetAttr("id", dialog.ID)
	})

	return bot, nil
}
