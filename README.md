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

We need a storage.

Chats cannot be deleted

User cannot be removed from chat

Oh and maybe it's time to start sanitizing raw user input

And UI looks like shit, have to style it

Currently I only store username as string in memory.
They just login by submitting this username.
Instead I want users to sign up fist, then login with user+pass.
Let's make an email-based signup with email message verification.
Later I'd want to add an SSO, like google or metamask. 
We should take this into account, but not build this yet.
We need to properly store encrypted credentials.
Let's use sqlite via sqlx lib for this.

### Message Broadcasting: 
- ~~When the server receives a message from a client over a regular HTTP POST request, it should broadcast that message to all other connected clients via the SSE connection.~~

### Persistence
- Cache instead of map
- DB, at least for users

### User Authentication
- Continue using your existing login system to authenticate users. 
- When a user logs in, generate a unique token for them and send it back. 
- The client should store this token and use it to authenticate the SSE connection.

### Establish SSE Connection
- ~~When a user opens the chat, the client should initiate an SSE connection to the server.~~
- The server should verify the user's token and if it's valid, accept the connection.
- ~~SSE should include at least: ping, chat invites, incoming messages~~

### GUI
- Style the app, pick a futuristic dark theme
- Should be simple and intuitive
- Should remain snappy
- Custom font

### GUI Interactivity
- Search chats by: 
    1. chat name
    2. invited user name
    3. message content
- Search messages by:
    1. content
    2. author

## LATERs

### Moderation
- Add ability to mute/block users in chat
- Add ability for users to mute/report other users

### *Security Considerations*
- Ensure that the SSE connection is secure (https://) and that we validate 
    and sanitize all incoming messages to prevent cross-site scripting (XSS) attacks.

### *GPT*:
- Consider for content moderation assistance
- Consider for virtual member
