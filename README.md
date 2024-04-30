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
