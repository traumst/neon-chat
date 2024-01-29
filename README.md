# MinChat

Minimalistic chat app built using server-components on go+htmx.
As simple as: register, create chat rooms, invite other users, chat.

## TODOs
User Authentication: 
    * Continue using your existing login system to authenticate users. 
    * When a user logs in, generate a unique token for them and send it back. 
    * The client should store this token and use it to authenticate the SSE connection.

Establish SSE Connection: 
    * When a user opens the chat, the client should initiate an SSE connection to the server. 
    * The server should verify the user's token and if it's valid, accept the connection.
    * SSE should include at least: ping, chat invites, incoming messages

Message Broadcasting: 
    * When the server receives a message from a client over a regular HTTP POST request, 
    * It should broadcast that message to all other connected clients via the SSE connection.

GUI
    * Style the app, pick a futuristic dark theme
    * Should be simple and intuitive
    * Should remain snappy
    * Custom font

_Security Considerations_: 
    * Ensure that the SSE connection is secure (https://) and that we validate 
        and sanitize all incoming messages to prevent cross-site scripting (XSS) attacks.

_GPT_:
    * Consider for content moderation assistance
    * Consider for virtual member
