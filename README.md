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
- ~~When the server receives a message from a client over a regular HTTP POST request, it should broadcast that message to all other connected clients via the SSE connection.~~

### Persistence
- Cache instead of map
- ~~DB for users and auth~~

### User Authentication
- ~~Continue using your existing login system to authenticate users.~~ 
- ~~When a user logs in, generate a unique token for them and send it back.~~
- The client should store this token and use it to authenticate __all__ connections.

### Establish SSE Connection
- ~~When a user opens the chat, the client should initiate an SSE connection to the server.~~
- ~~The server should verify the user's token and if it's valid, accept the connection.~~
- ~~SSE should include at least: ping, chat invites, incoming messages~~

### GUI
- Style the app, pick a futuristic dark theme
- Should be simple and intuitive
- Should remain snappy
- Custom font - jetbrains.com/lp/mono/

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
- ~~Ensure that the SSE connection is secure (https://)~~
- Validate nd sanitize all incoming messages to prevent cross-site scripting (XSS) attacks.

### *GPTs*:
- Consider for content moderation assistance
- Consider for chat participant - query, image, auto-response
