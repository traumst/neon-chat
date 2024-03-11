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

* User registraction sucks
* User cannot be removed from chat
* User cannot leave chat - admin can delete, user can leave

Oh and maybe it's time to start sanitizing raw user input

And UI looks like shit, have to style it

DB is close. Need to add user registration and authorization, thinking about 2 tables,
* user
    * id
    * name
    * salt?
* auth
    * id
    * user id
    * type
    * hash

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
