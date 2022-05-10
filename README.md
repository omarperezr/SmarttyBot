# SmarttyBot
## _The smartest multiplatform bot_

[![Go Report Card](https://goreportcard.com/badge/github.com/omarperezr/SmarttyBot)](https://goreportcard.com/report/github.com/omarperezr/SmarttyBot) [![Build Status](https://app.travis-ci.com/omarperezr/SmarttyBot.svg?branch=main)](https://app.travis-ci.com/omarperezr/SmarttyBot) [![Maintainability](https://api.codeclimate.com/v1/badges/5bdbd92a6997070a88e4/maintainability)](https://codeclimate.com/github/omarperezr/SmarttyBot/maintainability)

SmarttyBot is a multiplatform bot that uses natural language to execute functions using discord, telegram, email and gitlab 

- Set up all your accounts and get all the api tokens
- Create and train your [WitAI] model
- _SmarttyBot will now connect every platform to make your life easier_

## Features
- Connect any account with any other account to create unique solutions and automations
- Easily add your own functions with Python or Go
- Train your own WitAI model and easily connect your WitAI actions with your own methods
- Monitor your accounts and set up notifications to any other account
- Create your own telegram markup for menus
- _The only limit is your imagination_

## Tutorial

### Step 1: WitAI Setup
To set up WitAI you have to do 3 things:
- Login on [WitAI] with facebook
- Create your app
- Train your app

And that is it, with that you will have set up the natural language module, for more details go to https://wit.ai/docs/quickstart

After that:
- Go to https://wit.ai/apps/
- Select your app
- Go to settings
- Copy the `Server Access Token` into you .env
```sh
WIT_API_KEY="1234ASDFGHJ"
```

### Step 2: Telegram Setup
You will need to create a bot
- Go to https://telegram.me/botfather on telegram
- Send /newbot to create a new Telegram bot, 
- When asked, enter a name for the bot
- Give the Telegram bot a unique username. Note that the bot name must end with the word "bot" (case-insensitive).
- Copy and paste the Telegram bot's access token to the .env file.

```sh
TELEGRAM_API_KEY="123456TELEGRAMMMMMMAPI"
```
### Step 3: Discord Setup
To set up a discord bot follow this steps https://discordpy.readthedocs.io/en/stable/discord.html

Once you are done
- Go to 'Bot' 
- Select 'Reset token' 
- Copy the token that appears into the .env file
```sh
DISCORD_API_KEY="DDDIIISSSCCCOOORRRRDDDAAAAPPPIII"
```

### Step 4: Email setup
Simply add your email credentials and server settings like this
```sh
EMAIL_ACCOUNT="gmailaccount@gmail.com"
EMAIL_PASWORD="emailPassword"
SMTP_SERVER="smtp.gmail.com"
SMTP_PORT=587
IMAP_SERVER="imap.gmail.com"
IMAP_PORT=993
```

## Installation & Execution
SmarttyBot requires [Go] v1.18+ to run.

```git
git clone github.com/omarperezr/SmarttyBot
```

```sh
cd SmarttyBot
```

```go
go build
```

```sh
./main
```

## Donate
#### Want to contribute?
[![Kofi](https://az743702.vo.msecnd.net/cdn/kofi3.png?v=0)](https://ko-fi.com/omarperezr)

## License

GNU General Public License v3.0

**Free Software, Hell Yeah!**

[//]: # (Links used)

   [Go]: <https://go.dev/>
   [WitAI]: <https://wit.ai/>
