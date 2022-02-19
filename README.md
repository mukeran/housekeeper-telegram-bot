# housekeeper-telegram-bot

Part of ***HouseKeeper Project***.

This project is designed to automatically sign in various services, such as CCCAT.

## Setup

1. Install **Go SDK >= 1.14** and **Redis**;
2. Clone this project: `git clone https://github.com/mukeran/housekeeper-telegram-bot.git`;
3. Compile binary file: `cd housekeeper-telegram-bot && go build`;
4. Run the
   bot: `TELEGRAM_BOT_TOKEN="<YOUR-TELEGRAM-BOT-TOKEN>" ADMIN_ID="<YOUR-TELEGRAM-USER-ID>" ./housekeeper-telegram-bot`.

Note: If your Redis instance isn't running on port 6379, add `REDIS_ADDR` environment variable to specify.

## Command Descriptions

```
cccat_add - Add a CCCAT account
cccat_del - Delete a CCCAT account
cccat_update - Update Cookie user_auth of a CCCAT account
cccat_list - List all of your CCCAT account
cccat_sign - Start a CCCAT sign procedure
start - Obtain bot's main menu
help - Help information for HouseKeeper Telegram Bot
```

