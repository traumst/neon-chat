# Neon chat

Minimalistic chat app built using server-components on go + htmx.
Default colorscheme is a personal preference, and light/dark mode are... khm... on the horizon.
This app is being build as an excercise for me to
- learn go practically
- try sqlite db via sqlx, write custom migration module
- explore server sent events for live messaging
- further explore htmx, locality of behaviour, etc
- familiarize with tailwind
- <s>prove react is overrated</s>
- finish a project for a change
- maintain go project over time, optimize, extend functionality, fix bugs, etc.


<details>
    <summary>
        <b>Roadmap</b>
    </summary>
## TODOs
|-------------------------------------------------------------|
|---  from: 2024-10-29  --------------------------------------|
|-------------------------------------------------------------|
+ user info card (2 weeks)
    + should include:
        + avatar
        + name
        + contact info like email
        + mutual chats
    + optional order by:
        + most notifications in chat (default) - every message for 1-on-1
        + most unviewed messages
        + most recent activivity
        + most recently joined
+ store sessions in db, keep between restarts (2 days)
+ add api throttling (1 week)
+ sanitize user input against XSS attacks (1 week)
+ need [Regular Maintenance: VACUUM and Analyze] (1 week)
    + run `PRAGMA incremental_vacuum;` periodically to reclaim space
    + db does not shrink, need to [do a VACUUM](https://www.sqlite.org/lang_vacuum.html)
+ db
    + add db query-response caching (2 weeks)
    + consider [conn pool for db](https://github.com/jmoiron/sqlx/issues/300) (1 week)
    + add timestamps to every table:
        * created
        * updated
+ docker setup: pod + persistent storage [1 week]
+ deployment github action

### Next up

+ add log trace to all methods (1 day)
+ logrotate (1 week)
+ log levels with [slog in GO 1.23](https://pkg.go.dev/log/slog#Debug) (1 week)
+ change chat title (1 week)
+ collapsible / resizable left panel (2 weeks)
+ add contacts page / address book (2 weeks)
    * limit who can invite add / you

## Backlog

### Research
- stress test a bufferred vs unbuffered channel

### Message Broadcasting:
- msg should distribute to user connection, even if chat is closed
- pagination, track user deltas in chats and messages
- buffer outgoing events for unstable connection/s

### User Notifications
- new chat invite
- new msg in chat
- @user - display user card on hover
- setting on/off
    - mute / unmute chat
    - blacklist / whitelist users

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

### Testing
- extend unit tests
- add integration tests

### Moderation
- Add ability for users to mute/report other users
- Add ability to ban users from chat

### User Authentication
- throttling api actions - 10x exec time
- provide authType as form input
- google auth
- *2FA / MFA*
- *crypto wallet*

### Web Calls
- voice chat
- video chat
- effects

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

### *Mini Games*
- embed games into chat, to start
    - XO
    - Battle Ships
    - Chess
    - Go

### *GPTs*:
- Consider for content moderation assistance
- Consider for chat participant - query, image, auto-response

</details>

## Technical Details

### Prerequisites

I assume you already have golang installed.
You can check if it is by running `go version` in a terminal.
If not, you can refer to [official istallation instruction](https://go.dev/doc/install) for you pc/mac/linux.

There's [tailwind.config.js](./tailwind.config.js) in the root.
Technically we can simply make a 3rd party call and get entire tailwindcss.min.js from a CDN.
And we do - it's much easier for development.
Still, current setup will build tailwind and produce updated minimal css.
[Compiled css](./static/css/tailwind.css) comes out at roughly 15kb,
about 1/20 of the default minified [cdn provided talwindcss.js](https://cdn.tailwindcss.com/3.4.5) wich is around 300kb.

To actually compile this, you'd need to install tailwind. There are too many options.
I recommend downloading the [stadalone tailwind cli](https://tailwindcss.com/blog/standalone-cli) and avoiding npm completely.
But if you have nodejs installed and feel more comfortable with it,
you can [install tailwind via npm](https://tailwindcss.com/docs/installation).

> <b><u>Note</u></b><br>
[run.sh script](./run.sh) specifies path to tailwind executable and needs to be updated,
`~/code/bin/tailwindcss` is specific to my system, so unless you put it in exacty the
same place - you'd need to replace it with `npx tailwind` or something similar.

Finally, app expects to have `.env` file in the root directory, which you have to create.
There's `.env.template` file that you can easily copy-paste and fill up.
App may still start without this file, as some defaults are provided,
but the behaviour in this case is unpredicable. Broken state should fatal exit at some point in short future.

### Preparing run script

Successfull launch of the app requires:
* passing static checks and tests
* compiling tailwind
* launching the app
* etc.

`run.sh` file is just a plaintext shell script. We need to make it executable by running
`chmod +x run.sh` in project root in terminal. After that, we can start the app by running `sh run.sh` or `./run.sh`.

Contents of [run.sh](./run.sh) are more or less:
```sh
go vet ./... && staticcheck ./...
time go test ./...
~/code/bin/tailwindcss -i static/css/input.css -o static/css/tailwind.css
go run main.go
```

> <b><u>Note</u></b><br>
by default, db file is created in the root folder where executable runs

> <b><u>Note</u></b><br>
for db file to be deleted on start, must uncomment appropriate lines in `run.sh`

> <b><u>Note</u></b><br>
path to <b>tailwind</b> executable needs to be updated to match your system

### Run script

Execute a script from the terminal, you should then see output similar to:
```sh
> ./run.sh
Running tests
ok      neon-chat/src/handler     (cached)
ok      neon-chat/src/model/app   (cached)
Starting server...
2024/05/01 00:07:58 Application is starting...
...
2024/05/01 00:07:58 Starting server at port [8080]
```
App should now be available at http://localhost:8080

> <b><u>Note</u></b><br>
Executing `./run.sh` also builds the short tailwind.css.
Making changes to tailwind classes requires rerun to display properly.
Unless, we load the entire tailwind.min.css from CDN - then we have all classes available.

## Sqlite DB

You can imagine this entire app as an upside-down pyramid standing firmly
on the pinacle of sqlite. I can't believe I am only discovering this now.

> <b><u>Note</u></b><br>
When you first run the app with `./run.sh`, the `chat.db` file should be created.

> <b><u>Note</u></b><br>
You can specify a few test users in the `.env` file to pre-create users with auth in DB for tests.

### DB Options

There are a number of PRAGMAs defined in code. In general, pragmas adjust db settings. Changing these requires at least a restart, while some pragmas cannot be applied to a pre-existing db file at all. Some of the relevant pragmas are listed below, for detailed explanation and full list of all available pragmas, refer to [official documentation](https://www.sqlite.org/pragma.html#pragma_journal_mode).
* PRAGMA journal_mode = WAL;
* PRAGMA synchronous = NORMAL;
* PRAGMA locking_mode = NORMAL;
* PRAGMA foreign_keys = ON;
* PRAGMA journal_size_limit = 67108864;
* PRAGMA page_size = 4096;
* PRAGMA cache_size = 2000;
* PRAGMA mmap_size = 134217728;

### CLI Options

We can interact with our db from the terminal. In order to do so we would need a connector. Install [sqlite3 cli](https://www.sqlite.org/cli.html), it is very simple and it's what I use at the moment.

Connect with sqlite3 cli tool to the db file
``` sql
> sqlite3 chat.db
SQLite version 3.43.2 2023-10-10 13:08:14
Enter ".help" for usage hints.
sqlite>
```

Sqlite is just a normal sql, so we can run a normal select like below
```sql
sqlite> SELECT * FROM users WHERE status='active' LIMIT 2;
1|ABCDE|abcd@gmail.com|basic|active
2|NEW12|newt@gmail.com|basic|active
sqlite>
```

Data is provided without the headers - probably good default for connectors.
But when I look at it I prefer to have column names shown as well.
This can be achieved by running the two commands below in sqlite3 cli:
```sql
sqlite> .headers ON
sqlite>
```
```sql
sqlite> .mode columns
sqlite>
```
Here's the full reference for [sqlite dot commands](https://sqlite.org/cli.html#special_commands_to_sqlite3_dot_commands_).

These commands do not produce any output but if we run the same query again, the output will now be nicely formatted:
```sql
sqlite> SELECT * FROM users WHERE status='active' LIMIT 2;
id  name   email           type   status
--  -----  --------------  -----  ------
1   ABCDE  abcd@gmail.com  basic  active
2   NEW12  newt@gmail.com  basic  active
sqlite>
```

### Run from terminal

There are a couple of annoying settings when working directly in `sqlite>`, which I list below. Running your commands directly in the terminal solves all of them.

Some of the `sqlite>` tool's interface limitations:
* No command history from previous sessions
* Opitons like `.headers ON` are set for entire session or globally
* Global options can interfere with go's sqlx connector
* Queries require `;` at the end

```sh
sqlite3 chat.db "select * from users" -cmd ".headers ON" -cmd ".mode columns"
```

or with the shorthand
```sh
sqlite3 chat.db  -header -column  "select * from users;"
```

With this, all of queries are executed in a single session, we can find it in the history, and both `-cmd`s values are cleaned up for the next query. Semicolon is also not required for a single query. No global setting to mess with running code later.

## System requirements

I don't normally specify such info, because this is absolutely subjective and could probably work "fine" on less that half of minimal requirements. But then we start testing the limits of garbage collection and resource allocation. This is russian roulette and I love it! Test it out and report the lowest resource consumption you are able to achieve while serving a hot steamy load.

### Pre-measure estimates

What **I** think a good measure is 
* ~1000 active/concurrent users 
* 100-300 mixed calls per second
* responses always under 200ms

#### Minimal Requirements

Obviously our machine is expected to be bottlenecked by the CPU the most here, expect occasional hang or reboot.

* 1 CORE    - with arbitrary performance
* 256MB RAM - for app memory and db caches with ~30% spair
* 2GB DISK  - for dependencies, db file, logs with ~60% spair

#### Recommended Requirements

Here we should not cross ~85% cpu anymore, memory should be at max ~50% in use

* 2+ CORES   - to better utilize go concurrency
* 512MB+ RAM - lower % usage = faster RAM lookup
* 4GB+ DISK  - same as RAM for SSD, otherwise it's only log space for HHD