# MinChat

Minimalistic chat app built using server-components on go + htmx.
As simple as: register, create chat rooms, invite other users, chat.
This app is being build as an excercise for me to 
- learn go
- further explore htmx abilities
- explore SSE as websocket alternative
- <s>prove react is overrated</s>

## TODOs

### *Right now...*

Leaving chat should update users in chat

It's time to start sanitizing raw user input.
Def break the chat: ' "

* HTML hx- stuff should sit next to target
* HTML button onclick should
    * clear related textbox
    * new chat -> remove chat-list placeholder
    * new msg -> remove chat-history placeholder
* User registraction sucks

And UI looks like shit, have to style it

### Message Broadcasting: 
- track user deltas: chats, messages
- only serve deltas
- test overload of a conn channel

### Persistence
- message store
- cache instead of map

### User Authentication
- client should send session token to authenticate __all__ connections
- email verification
- 2FA

### GUI
- style the app, pick a futuristic dark theme
- should remain minimal and snappy
- apply custom font jetbrains.com/lp/mono/

### Chat features
- Search chats by: 
    1. chat name
    2. invited user name
    3. message content
- Search messages by:
    1. content
    2. author

### Client storage
- local / innodb
- store chats with history on client
- load only chat deltas

### Moderation
- Add ability to mute/block users in chat
- Add ability for users to mute/report other users

### *Security Considerations*
- Validate nd sanitize all incoming messages to prevent cross-site scripting (XSS) attacks.

## Later

### *GPTs*:
- Consider for content moderation assistance
- Consider for chat participant - query, image, auto-response
