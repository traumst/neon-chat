# Neon chat

Minimalistic chat app built using server-components on go + htmx.
Default colorscheme is a personal preference, with light/dark mode on the horizon.
This app is being build as an excercise for me to 
- learn go
- try sqlx
- further explore htmx abilities
- explore server sent events as websocket alternative
- fiddle with tailwind
- <s>prove react is overrated</s>
- finish a project for a change

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
ok      neon-chat/src/handler     (cached)
ok      neon-chat/src/model/app   (cached)
Starting server...
2024/05/01 00:07:58 Application is starting...
...
2024/05/01 00:07:58 Starting server at port [8080]
```
And the app should be available at http://localhost:8080

## TODOs

+ Bugs
    * User name change should update left panel

+ UI improvements
    * collapsible menus
    * avatars on messages
    * @ other messages
    * @ other users

+ Fuzzy search
    * messages by author
    * messages by content
    * in chat

### Message Broadcasting: 
- msg should distribute to user connection, even if chat is closed
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
- introduce Tmpl to replace default templating engine

### User Notifications
- setting on/off
    - mute / unmute chat
- new chat invite
- new msg in chat

### Search
- search chats by: 
    - chat name
    - message content
- search messages by:
    - content
    - author
- fuzzy search methods:
    - common key [C530, V500] - fast - mostly latin
        - https://www.sqlite.org/lang_corefunc.html#soundex
    - word embeding - %VALUE%
    - edit distance [cindy-cindi=1] - only latin - in-memory
    - statistical similarity - slow

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
    * msg options

## Later

### Client storage
- local / innodb
- store chats with history on client
- load only chat deltas

### *Security Considerations*
- Validate nd sanitize all incoming messages to prevent cross-site scripting (XSS) attacks.

### *Mini Games*
- embed games into chat, to start
    - XO
    - Battle Ships
    - Chess
    - Go

### *GPTs*:
- Consider for content moderation assistance
- Consider for chat participant - query, image, auto-response
