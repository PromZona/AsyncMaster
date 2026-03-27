# AsyncMaster

Async Master is a telegram bot that allows you to play text RPGs with your friends.
Pick a master of the game, invite players to bot and start your play.
Functionality is minimal to support any types of games:

- Each player chooses a Player name, which adds layer of anonymity 
- Messages between players and master
- Master can make requests to players, where can ask to roll dice

In this iteration of the program, players controll factions, which they decide on while registration process.
Each faction has name, description and resources.

# How to run

1) Install and setup postgresql
2) Create Database
4) Fill *.env* file
```
BOT_TOKEN=<Telegram Token>
BOT_USER_PASSWORD="<Choose a password for players>"
BOT_MASTER_PASSWORD="<Choose a password for master>"

# AsyncMaster uses GOOSE library to handle migrations
# 
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=<Connection string to db>
GOOSE_MIGRATION_DIR=./migrations
```
5) Build and run
```shell
go build
./AsyncMaster
```

# Mock Runtime
AsyncMaster supports flag * -mock * to start testing in a console, without telegram.
It allows you to create *fake* telegram users to test handles.
This makes it convinient to develop features, without creating a lot of telegram account

To start mock runtime
```shell
./AsyncMaster -mock
```

## Supported commands
### /server
Server command mock telegram sides of things. Because we do not want to buy 10 sim cards to simulate different users,
this command will provide unique ChatID of the user and name.
```shell
# Creates user for testing, telegram_name will be used to make user to do something
/server create_user <fake_telegram_name> <player name>
```

### /user
User commands simulate user writing from the telegram to the bot.

There are 2 ways user can communicate: write some text to chat and press a button.
Therefore, we have command templates for those scenarious
```shell
# User sends text in chat with bot
/user <telegram_name> "text goes here"

# User presses button. Button contains unique identifier, and may contain some data
/user <telegram_name> <callback_string>|<data_1>:<data_2>
```

## Example Usage

```shell
# Create 2 users in Mock Runtime
/server create_user John Vitus
/server create_user Victor Horrow

# User 1 registration
/user John "123"
/user John "Vitus"
/user John "Tigers"
/user John "Stripe tigers. Very powerful"

# User 2 registration
/user Victor "123"
/user Victor "Horrow"
/user Victor "Necro Guild"
/user Victor "Super scary necromancers. And a lot of money"

# User 1 *presses* button that requests faction data
/user John factions_list
```

Mock runtime also used in tests. So if you want to check more examples of its usage, you can check *tests* directory

# Testing

There are only integration tests for now. Tests use Mock Runtime for convinience of testing several users.
Tests are integraion ones and test the whole pipeline of inputs from user to DB. Therefore, tests require from you to setup testing database in .env.test file

1) Fill *./internal/app/tests/.env.test* file. Provide testing database connection string.
Do not change password, those are hardcoded into tests (yeah, yeah, I know)

2) Run all tests
```shell
go test ./internal/app/tests/
```

or run tests separately like golang documentation propose

# Thanks to libraries
- Telebot
- goose
