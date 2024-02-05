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
...working on adding channels and some concurrency to distribute
Ok, now I want to have a loop in a go routine, that listens on a channel.
I need to hold a list of Clients, contaning: 
* user 
* open chat
* connection - r and w

I will then loop over channel messages and distribute them between users in chats. 

### Message Broadcasting: 
- When the server receives a message from a client over a regular HTTP POST request, it should broadcast that message to all other connected clients via the SSE connection.

### User Authentication
- Continue using your existing login system to authenticate users. 
- When a user logs in, generate a unique token for them and send it back. 
- The client should store this token and use it to authenticate the SSE connection.

### Establish SSE Connection
- When a user opens the chat, the client should initiate an SSE connection to the server. 
- The server should verify the user's token and if it's valid, accept the connection.
- SSE should include at least: ping, chat invites, incoming messages

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

## MAYBEs

### *Security Considerations*
- Ensure that the SSE connection is secure (https://) and that we validate 
    and sanitize all incoming messages to prevent cross-site scripting (XSS) attacks.

### *GPT*:
- Consider for content moderation assistance
- Consider for virtual member
