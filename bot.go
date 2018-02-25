package main

import (
	"github.com/PuerkitoBio/goquery"
)

type Bot struct {
	Document *goquery.Document
	Configs  map[string]string
	Buttons  map[string]*Button
	Inputs   map[string]*Input
	Menus    map[string]*Menu
	Dialogs  map[string]*Dialog
}

type Menu struct {
	ID      string
	Title   string
	If      string
	Buttons []*Button
}

type Button struct {
	ID    string
	Title string
	Href  string
	Embed string
	If    string
	Reset bool
}

type Dialog struct {
	Node   *goquery.Selection
	ID     string
	Action string
	Help   string
	If     string
	Inputs []*Input
}

type Input struct {
	NS      string
	ID      string
	Title   string
	Name    string
	Type    string
	Options []*Button
	Value   string
	If      string
}
