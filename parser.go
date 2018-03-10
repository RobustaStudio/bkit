package main

import (
	"fmt"
	"io"
	"log"
	"strings"

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
		if key == "" {
			return
		}
		bot.Configs[key] = value
	})

	// paginate the menu/navs buttons
	// messanger maximum navigation links are 3 - 4, we will make them 3 by maximum.
	doc.Find("menu,nav").Each(func(_ int, m *goquery.Selection) {
		if m.AttrOr("id", "") == bot.Configs["main-menu"] {
			return
		}
		p, cnt := m, 1
		more := bot.Configs["pager-more"]
		if more == "" {
			more = "More"
		}
		m.Find("a,button").Each(func(i int, b *goquery.Selection) {
			if (i > 0) && (i%2 == 0) && m.Find("a,button").Length() > 3 {
				newId := fmt.Sprintf("%s-%d", m.AttrOr("id", ""), cnt)
				p.AppendHtml(fmt.Sprintf("<a href='#%s'>%s</a>", newId, more))
				p.AfterHtml(fmt.Sprintf("<menu id='%s' title='%s'></menu>", newId, m.AttrOr("title", "")))
				p = bot.Document.Find("#" + newId)
				cnt++
			}
			if p != m {
				p.AppendSelection(b.Clone())
				b.Remove()
			}
		})
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
			input.Title = n.AttrOr("title", "Fill/Choose ...")
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
			n.SetAttr("title", input.Title)
		}
		d.Find("div").Children().Each(func(i int, s *goquery.Selection) {
			s.SetAttr("if", s.Parent().AttrOr("if", ""))
			s.Parent().BeforeSelection(s.Clone())
		})
		d.Find("div").Remove()
		d.Find("input,select").Each(func(i int, n *goquery.Selection) {
			switch strings.ToLower(goquery.NodeName(n)) {
			case "input":
				fn(i, n)
			case "select":
				n.SetAttr("type", "options")
				fn(i, n)
			}
		})
		bot.Dialogs[dialog.ID] = dialog
		d.SetAttr("id", dialog.ID)
	})

	// populate the navs
	doc.Find("menu,nav").Each(func(i int, m *goquery.Selection) {
		menu := &Menu{}
		menu.ID = m.AttrOr("id", fmt.Sprintf("menu%d", i))
		menu.Title = m.AttrOr("title", "Choose ...")
		menu.Buttons = []*Button{}
		menu.Inline = (m.AttrOr("inline", "false") == "true")
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
		m.SetAttr("title", menu.Title)
		bot.Menus[menu.ID] = menu
	})

	// finding errors
	doc.Find("menu,nav").Find("a,button").Each(func(i int, b *goquery.Selection) {
		href := b.AttrOr("href", "")
		title := b.AttrOr("title", b.Text())
		btnErr := "Invalid `href` (href=" + href + ") attribute of the button(" + title + ") in the Menu/Nav(" + b.Parent().AttrOr("id", "") + ")"

		if string(href[0]) != "#" && !strings.HasPrefix(href, "http://") && !strings.HasPrefix(href, "https://") {
			bot.Errors = append(bot.Errors, "[WrongID Format] "+btnErr)
		}

		if href == "" {
			bot.Errors = append(bot.Errors, "[EMPTY] "+btnErr)
		}

		if string(href[0]) == "#" && doc.Find(href).Length() == 0 {
			bot.Errors = append(bot.Errors, "[NOT FOUND] "+btnErr)
		}
	})

	if len(bot.Errors) > 0 {
		for _, msg := range bot.Errors {
			log.Println("[CompilerError]", msg)
		}
		log.Fatal("Please fix the above errors")
	} else {
		log.Println("No errors found ...")
	}

	return bot, nil
}
