# Neon chat

Minimalistic chat app built using server-components on go + htmx.
Default colorscheme is a personal preference, and light/dark mode are... khm... on the horizon.
This app is being build as an excercise for me to 
- learn go practically
- try sqlite db via sqlx
- further explore htmx
- explore server sent events for live messaging
- familiarize with tailwind
- <s>prove react is overrate</s>
- finish a project for a change

## How to run

### Prerequisites

First run is special. We're now going to do a couple of things that we would not have to do on subsequent runs. At least not entirely.

I assume you already have golang installed.
You can check if it is by running `go version` in a terminal.
If not, you can refer to [official istallation instruction](https://go.dev/doc/install) for you pc/mac/linux.

There's [tailwind.config.js](./tailwind.config.js) in the root, but it's not in use YET. Meaning we don't have custom tailwind definition, only relying on built-in utility classes. Technically we can simply make a 3rd party call and get entire tailwindcss.min.js from a CDN. Still, in the current setup we built tailwind file to produce minimal required css. [Compiled css](./static/css/tailwind.css) comes out at roughly 15kb, about 1/20 of the default minified [cdn provided talwindcss.js](https://cdn.tailwindcss.com/3.4.5) wich is around 300k and works exactly the same.

To actually compile this, you'd need to install tailwind. There's plenty of options. I recommend downloading the [stadalone tailwind cli](https://tailwindcss.com/blog/standalone-cli) and avoiding npm completely. But if you have nodejs installed and feel more comfortable with it, you can [install tailwind via npm](https://tailwindcss.com/docs/installation). Note that [run.sh script](./run.sh) specifies path to tailwind executable. That line may need to be updated to match your setup.

Finally, app expects to have `.env` file in the root directory, which you have to create. There's `.env.template` file that you can easily copy-paste and fill up. App may still start without this file, as some defaults are provided. But the behaviour in this case is unpredicable, thus is a broken state and should fatal exit at some point.

### Preparing run script

Successfull launch of the app requires:
* passing tests
* compiling tailwind
* (optional) purging the db
* launching the app
* etc.

`run.sh` file is just a plaintext shell script. We need to make it executable by running `chmod +x run.sh` in project root in terminal. After that, we can start the app by running `sh run.sh` or `./run.sh`.

Contents of [run.sh](./run.sh) is more or less:
```
> cat run.sh
echo "Running tests"
time go test ./...

# for dev purposes
#echo "Dropping db file..."
# rm chat.db
# (la chat.db && echo "...Dropped db successfully.") || echo "...Data not dropped."

echo "Building tailwind..."
~/code/bin/tailwindcss -i static/css/input.css -o static/css/tailwind.css

echo "Starting server..."
go run main.go
```

Notes:
1. by default, db file is created in the root folder where executable runs
2. for db file to be deleted on start, must uncomment appropriate lines
3. tailwind executable call will need to be updated to match your system

### Run script

Executing a script from the terminal, you should then see output similar to:
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
App should now be available at http://localhost:8080

> <b><u>Note</u></b><br>After initial successfull run there's no real reason to execute run script. Unless css is messed up or db is corrupted I just use `go run main.go` or `go test ./...` directly.

## TODOs

### Bugs

* user name change should update active user info on left panel
* controllers should only touch templates, never models
* fucking session transactions

### Next up

+ UI improvements
    * @ other messages
        - "reply" button on message
            + click should paste quote into input at carret
        - update text processing
            + quote source FK relationship
            + link to messages:
                + all or nothing when editing
                + removing any char from quote removes entire quote
                + default unavailable message
        - add quote message html template
            + like message_li but lighter
        - click on quoted message
            + scroll original message into view if available
    * user info card - avatar, name, email?, mutual chats?, 
    * @ other users - display user card on hover
    * collapsible menus

## Backlog

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
- throttling api actions - 10x exec time
- middleware auth
- provide authType as form input
- google auth
- *2FA / MFA*
- *crypto wallet*

### Extend functionality
- change chat title
- add contacts page / address book
- limit who invites to contacts
- black and whitelist
- zoom, web calls
- introduce Tmpl to replace default templating engine

### User Notifications
- setting on/off
    - mute / unmute chat
- new chat invite
- new msg in chat

### Fuzzy search
- search chats by: 
    - by chat name
    - by members
    - by message content
- search messages:
    - by messages content
    - by messages author
    - in chat
- approaches
    - common key [C530, V500] - fast - mostly latin
        - https://www.sqlite.org/lang_corefunc.html#soundex
    - word embeding - %VALUE%
    - edit distance [cindy-cindi=1] - only latin - in-memory
    - statistical similarity - slow

### GUI
- user settings page
    * add alternative auth
    * light/dark mode
- top controls
    * local time now iso
    * local session expiration iso
    * settings link
- bottom controls
    * status 🟢🟡🔴
    * light/dark mode switch
- collapsible sub menus
    * active user - logout, setting
    * open chat - close, delete
    * chat members - expel
    * msg options - delete, reply

## Never gonna happen, but sounds nice

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
