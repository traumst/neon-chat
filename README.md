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
- same user can be invited to the same chat 2+ times
- last open session is the only session a user can have
- add/remove chat/msg should be done by sse, not api response

### GUI
- top and bottom controls
- user settings page

### Moderation
- Add ability to mute/block users in chat
- Add ability for users to mute/report other users

### User Notifications
- mute / unmute chat
- new chat invite
- new msg in chat

### User Authentication
- middleware auth
- email verification
- 2FA

### Persistence
- message store
- cache instead of map

### Chat features
- Search chats by: 
    1. chat name
    2. invited user name
    3. message content
- Search messages by:
    1. content
    2. author

### Message Broadcasting: 
- track user deltas: chats, messages
- only serve deltas
- buffer for unstable connection/s
- test overload of a conn channel

### Client storage
- local / innodb
- store chats with history on client
- load only chat deltas

## Later

### *Security Considerations*
- Validate nd sanitize all incoming messages to prevent cross-site scripting (XSS) attacks.

### *GPTs*:
- Consider for content moderation assistance
- Consider for chat participant - query, image, auto-response
