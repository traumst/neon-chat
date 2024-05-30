# prplchat

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
ok      prplchat/src/handler     (cached)
ok      prplchat/src/model/app   (cached)
Starting server...
2024/05/01 00:07:58 Application is starting...
...
2024/05/01 00:07:58 Starting server at port [8080]
```
And the app should be available at http://localhost:8080

## Right now
bugs and refactor

## TODOs

### Known bugs
- last open session is the only session a user can have

### Persistence
- message store
- cache instead of map

### Message Broadcasting: 
- track user deltas: chats, messages
- only serve deltas
- buffer for unstable connection/s

### Testing
- extend unit tests
- add integration tests

### Moderation
- Add ability for users to mute/report other users
- Add ability to ban users from chat

### User Authentication
- middleware auth
- provide authType as form input
- email auth
- google auth
- *2FA / MFA*

### Extend functionality
- @users in chat
- @messages in chat
- *web call*

### User Notifications
- setting on/off
    - mute / unmute chat
- new chat invite
- new msg in chat

### Search
- fuzzy name matching methods:
    - common key [C530, V500] - fast - mostly latin
        - https://www.sqlite.org/lang_corefunc.html#soundex
    - word embeding - %VALUE%
    - edit distance [cindy-cindi=1] - only latin - in-memory
    - statistical similarity - slow
- search chats by: 
    - chat name
    - invited user name
    - message content
- search messages by:
    - content
    - author

### GUI
- user settings page
    * add alternative auth
    * light/dark mode
    * mute list - per user - per chat
- top controls
    * local time now iso
    * local session expiration iso
    * settings link
- bottom controls
    * status 🟢🟡🔴
    * light/dark mode switch
    * mute user
- collapsible sub menus
    * active user - logout, setting
    * open chat - close, delete
    * chat members

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
