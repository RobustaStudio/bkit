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
        <meta name="messenger-app-secret" content="***************************" />

        <!-- the verify token in the bot settings in facebook messenger platform -->
        <meta name="messenger-verify-token" content="*************" />

        <!-- the facebook page token -->
		<meta name="messenger-page-token" content="****************" />

        <!-- this content will be displayed when a new user opens the messenger window -->
        <meta name="description" content="Hi, I'm your bot"/>
	</head>
	<body>
        <menu id="main" title="Under your command, Sir :)">
            <!-- reset means "clear the current session if there were some old previous data" -->
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
> save that file as `demo.html`
> then just run `bkit -html "demo.html" -https ":443" -http ":80" -server-name "bkit.domain.com"`
> Point your facebook messenger webhook to `https://bkit.domain.com/messenger`

Installation
===============
- Binaries ? go to [Releases Page](releases) and select your own distro/arch.
- Docker ? `docker run -v $(pwd)/demo.html:/demo.html --network host alash3al/bkit -html /demo.html -https ":443" -http ":80" -server-name "bkit.domain.com"`
- From Source ? `go get github.com/RobustaStudio/bkit`

Credits
==============
Copyright 2018 (c) [Robustastudio](https://robustastudio.com)
