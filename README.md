BKIT
====
> An advanced tiny bots scripting engine using the power of XML/HTML.

Features
====
- Portable & Tiny
- Embedded a tiny `Expression Engine` to support simple scripting
- Supports the following tags `meta, text, label, p, line, inline, span, embed, resource, media, div, menu, nav, dialog, form, template`
- Supports custom replies from custom templates

Status
====
> `BKIT` is in its first release, we wanted to keep it as simple as possible, we are planning to add more features in the future

Tags
===
- Configs: `<meta />`
- Text: `<text></text>`, `<label></label>`, `<span></span>`, `<p></p>`, `<line></line>`, `<inline></inline>`
- Media: `<media src="" type="" />`, `<resource src="" type="" />`, `<embed type="" src="" />`
- Menu: `<menu></menu>`, `<nav></nav>`
- Form: `<form></form>`, `<dialog></dialog>`
- Template: `<template></template>`

Demo
=====
> a very simple messenger bot that will collects the user info 
```html
<html>
	<head>
        <!-- here we define the main-menu element, `main` is the id of the main menu -->
        <meta name="main-menu" content="main" />
        
        <!-- whether the bkit server should verify the incoming request "from messenger itself" or not -->
        <meta name="messenger-verify" content="true" />
        
        <!-- set the facebook app secret -->
        <meta name="messenger-app-secret" content="1a1176bf562a6645b0d168d7b0a6088e" />

        <!-- the verify token in the bot settings in facebook messenger platform -->
        <meta name="messenger-verify-token" content="123456" />
        
        <!-- the page token -->
		<meta name="messenger-page-token" content="EAAZAouu3OSM0BAPDpcHxXFVXZArZBVD8dsp6YlRqcdU63AjoEzKwCjZBsTjQt8lWAedZBY3MheC1UBZBflD5SXj1EQ3SOxa1Dc8wZBLqeu8D2sQRtWpuhs4laeZBKFbwsZBnh0PtIyUTiI3jLceTd7gnpOkf43FoEjtZBsHHcvFwzhW5dCZB7Bw2ZB41" />
        
        <!-- this content will be displayed when a new user opens the messenger window -->
        <meta name="description" content="Welcome to robusta!, I’m ro-bot, not your average customer care chatbot! Check the menu"/>
	</head>
	<body>
        <menu id="main" title="Under your command, Sir :)">
			<a href="#about" reset="true">About</a>
			<a href="#collect" reset="true">Collect Data</a>
        </menu>
        
        <div id="about">
            <text>I'm bkit, the bots-kit engine</text>
            <text>;)</text>
        </div>

        <form id="collect" action="http://some/backend">
            <text>enter your name</text>
            <input name="user_name" />

            <text>'Hi ' + user_name + ' ;)'</text>
            <text>Select a type ...</text>
            <select name="type">
                <option value="a">Type 1</option>
                <option value="b">Type 2</option>
            </select>
            <text if="type == 'a'">You selected Type 1</text>
            <text if="type == 'b'">You selected Type 2</text>
            <text>Thank You!</text>
        </form>
    </body>
</html>
```

Installation
