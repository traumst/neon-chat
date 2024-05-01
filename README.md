# MinChat

Minimalistic chat app built using server-components on go + htmx.
As simple as: register, create chat rooms, invite other users, chat.
This app is being build as an excercise for me to 
- learn go
- further explore htmx abilities
- explore server sent events as websocket alternative
- fiddle with tailwind
- <s>prove react is overrated</s>

## How to run

1. I assume you already have golang installed.
    You can check if it is by running `go version` in a terminal.
    If not, you can refer to [official instruction](https://go.dev/doc/install) for you pc/mac/linux
2. Create `.env` file in the root of the project.
    An example for available values can be found in `.env.template`
3. Make `run.sh` executable by running `chmod +x run.sh` in the terminal.
4. Execute `sh run.sh` in the terminal.

You should then see output similar to:
```
> sh run.sh
Running tests
ok      go.chat/src/handler     (cached)
ok      go.chat/src/model/app   (cached)
Starting server...
2024/05/01 00:07:58 Application is starting...
2024/05/01 00:07:58     setting up logger...
2024/05/01 00:07:58     parsing config...
2024/05/01 00:07:58     parsed config: {LoadLocal:false,Port:8080,Sqlite:chat.db}
2024/05/01 00:07:58     connecting db...
2024/05/01 00:07:58   opening db file [chat.db] [24576]
2024/05/01 00:07:58     init app state...
2024/05/01 00:07:58     init controllers...
2024/05/01 00:07:58 Starting server at port [8080]
```
And the app should be available at http://localhost:8080

## TODOs

### Known bugs
- last open session is the only session a user can have

### GUI
- top controls
    1. local time now iso
    2. local session expiration iso
    3. settings link
- bottom controls
    1. status 🟢🟡🔴
    2. light/dark mode switch
    3. mute user
- user settings page
    1. rename
    2. avatar?
    3. alternative auth options
    4. light/dark mode
    5. mute list - per user - per chat

### Moderation
- Add ability for users to mute/report other users
- Add ability to ban users from chat

### User Notifications
- setting on/off
    - mute / unmute chat
- new chat invite
- new msg in chat

### Search
- fuzzy name matching methods
    - https://www.sqlite.org/lang_corefunc.html#soundex
    - word embeding - %VALUE%
    - common key [C530, V500] - fast - mostly latin
    - edit distance [cindy-cindi=1] - only latin - in-memory
    - statistical similarity - slow
- search chats by: 
    1. chat name
    2. invited user name
    3. message content
- search messages by:
    1. content
    2. author

### User Authentication
- middleware auth
- email verification
- 2FA

### Persistence
- message store
- cache instead of map

### Message Broadcasting: 
- track user deltas: chats, messages
- only serve deltas
- buffer for unstable connection/s
- test overload of a conn channel

## Later

### Client storage
- local / innodb
- store chats with history on client
- load only chat deltas

### *Security Considerations*
- Validate nd sanitize all incoming messages to prevent cross-site scripting (XSS) attacks.

### *GPTs*:
- Consider for content moderation assistance
- Consider for chat participant - query, image, auto-response
