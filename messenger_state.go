package main

// type User struct {
//     writer  *messenger.Response
//     sender  messenger.Sender
// }
//
// func

// func (m Messenger) IsExpectingUserInput() bool {
// 	dialog := m.GetCurrentDialog()
// 	if dialog == nil {
// 		return false
// 	}
// 	return m.GetDialogData(dialog.ID)["current_input"] != ""
// }
//
// func (m Messenger) SendGreeting() error {
// 	if desc, ok := m.bot.Configs["description"]; ok  {
// 		return m.client.GreetingSetting(desc.(string))
// 	}
// 	return nil
// }
//
// // SendMainMenu - sending the main messenger menu
// func (m Messenger) SendMainMenu(menu *Menu) error {
// 	if menu == nil {
// 		return nil
// 	}
// 	callToActions := []messenger.CallToActionsItem{}
// 	for _, btn := range menu.Buttons {
// 		cta := messenger.CallToActionsItem{}
// 		cta.Title = btn.Title
// 		if strings.HasPrefix(btn.Href, "http://") || strings.HasPrefix(btn.Href, "https://") {
// 			cta.Type = "web_url"
// 			cta.MessengerExtension = (btn.Embed != "no") && (btn.Embed != "")
// 			cta.WebviewHeightRatio = btn.Embed
// 		} else {
// 			cta.Type = "postback"
// 			cta.Payload = fmt.Sprintf("trgt=%s&src=%s", btn.Href[1:], btn.ID)
// 		}
// 		callToActions = append(callToActions, cta)
// 	}
// 	return m.client.CallToActionsSetting("existing_thread", callToActions)
// }
//
// func (m Messenger) SendText(r *messenger.Response, text string) error {
// 	return r.Text(text)
// }
//
// func (m Messenger) SendMedia(r *messenger.Response, typ, src string) error {
// 	return r.Attachment(messenger.AttachmentType(typ), src)
// }
//
// func (m Messenger) SendInput(r *messenger.Response, input *Input) error {
// 	if input == nil {
// 		return nil
// 	}
// 	if input.Type == "options" {
// 		callToActions := []messenger.StructuredMessageButton{}
// 		for _, btn := range input.Options {
// 			cta := messenger.StructuredMessageButton{}
// 			cta.Title = btn.Title
// 			cta.Type = "postback"
// 			cta.Payload = btn.Href
// 			callToActions = append(callToActions, cta)
// 		}
// 		return r.ButtonTemplate(menu.Title, &callToActions)
// 	}
// 	return nil
// }
//
// func (m Messenger) SendBasicMenu(r *messenger.Response, menu *Menu) error {
// 	if menu == nil {
// 		return nil
// 	}
// 	callToActions := []messenger.StructuredMessageButton{}
// 	for _, btn := range menu.Buttons {
// 		cta := messenger.StructuredMessageButton{}
// 		cta.Title = btn.Title
// 		if strings.HasPrefix(btn.Href, "http://") || strings.HasPrefix(btn.Href, "https://") {
// 			cta.Type = "web_url"
// 			cta.MessengerExtensions = (btn.Embed != "no") && (btn.Embed != "")
// 			cta.WebviewHeightRatio = btn.Embed
// 		} else {
// 			cta.Type = "postback"
// 			cta.Payload = fmt.Sprintf("trgt=%s&src=%s", btn.Href[1:], btn.ID)
// 		}
// 		callToActions = append(callToActions, cta)
// 	}
// 	return r.ButtonTemplate(menu.Title, &callToActions)
// }
//
// func (m Messenger) GetMainMenu() *Menu {
// 	mainMenu := m.bot.Configs["main-menu"]
// 	if nil == mainMenu {
// 		return nil
// 	}
// 	menu := m.bot.Menus[mainMenu.(string)]
// 	return menu
// }
