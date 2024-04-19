# MinChat

Minimalistic chat app built using server-components on go + htmx.
As simple as: register, create chat rooms, invite other users, chat.
This app is being build as an excercise for me to 
- learn go
- further explore htmx abilities
- explore server sent events as websocket alternative
- fiddle with tailwind
- <s>prove react is overrated</s>

## TODOs

### Known bugs
* chat-create accumulates and sends first msg N times to each receiver
* last open session is the only session a user can have

### Pressing issues
* UI looks like shit, have to style it NOW
* it's time to start sanitizing raw user input.
    * Def break the chat: ' "
* HTMX hx- stuff should sit next to target

### Message Broadcasting: 
- track user deltas: chats, messages
- only serve deltas
- test overload of a conn channel

### Persistence
- message store
- cache instead of map

### User Authentication
- middleware auth
- email verification
- 2FA

### GUI
- style the app with dark theme with green/purple accents
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
