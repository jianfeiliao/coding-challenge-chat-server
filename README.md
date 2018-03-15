# Coding Challenge Chat Server
A simple chat server implementation for [a coding challenge](https://gist.github.com/jianfeiliao/dba0f0e7da7c978741fd67b091e21288).

# Assumptions & Future Improvements
This solution make assumption about the user input would be valid, so therefore there is no input validation.

These are some of the possible future improvements that can be added:
* Adding a basic CI/CD loop with some kind of build script/tooling.
* Adding unit tests & integration tests etc.
* Handle reconnect. (currently if you disconnect and try to reconnect as the same user, the server won't allow it)
* Better indication for when the user can enter input.

And much more as if this will become a real thing.

# Old Version
An older version of this solution is available at [here](/old_version). This is basically where it left off during the remote session.

# Running the server
Make sure you have [Go](https://golang.org/) installed. This is tested against my Go 1.9.1 installation on Windows 10. Also port 8080 needs to be available.

You can run the server as below:
```
$ go run *.go
2018/03/14 19:12:20 Starting chat server
```

# Connecting to the server
You can use something like `netcat` or `telnet` to connect to the server. This is tested using this version of [netcat 1.11](https://eternallybored.org/misc/netcat/) on the Windows command prompt, as well as [Docker quickstart terminal](https://docs.docker.com/toolbox/toolbox_install_windows/).

You can connect to the server as below:
```
C:\>nc localhost 8080
< Welcome to my chat server! What's your name?
```

You can disconnect anytime by entering CTRL-C.

# Example interaction output

### Server
```
$ go run *.go
2018/03/14 19:35:36 Starting chat server
2018/03/14 19:36:16 User disconnected before establishing a valid connection
2018/03/14 19:36:22 User fei has connected
2018/03/14 19:36:49 User jake has connected
2018/03/14 19:39:52 User bob has connected
2018/03/14 19:40:32 User jake has disconnected
```

### User 1
```
$ nc localhost 8080
< Welcome to my chat server! What's your name?
fei
< You are the only one online right now
ok, so lonely
helloooo
< [19:36:49] *jake has joined the chat*
< [19:37:00] <jake> hello fei
what's up jake!
this is a cool coding challenge!
i can send you a beep too @jake
< [19:37:53] <jake> that's cool
sending more messages
sending more messages again
that should be more than 10 messages
< [19:39:52] *bob has joined the chat*
< [19:40:06] <bob> hey guys
welcome to the party
< [19:40:31] <jake> sorry guys, g2g and ttyl
< [19:40:32] *jake has left the chat*
```

### User 2
```
$ nc localhost 8080
< Welcome to my chat server! What's your name?
jake
< You are connected with users: [fei]
< [19:36:22] *fei has joined the chat*
< [19:36:34] <fei> ok, so lonely
< [19:36:38] <fei> helloooo
hello fei
< [19:37:10] <fei> what's up jake!
< [19:37:25] <fei> this is a cool coding challenge!
< [19:37:47] <fei> i can send you a beep too @jake
that's cool
< [19:38:38] <fei> sending more messages
< [19:38:42] <fei> sending more messages again
< [19:38:54] <fei> that should be more than 10 messages
< [19:39:52] *bob has joined the chat*
< [19:40:06] <bob> hey guys
< [19:40:17] <fei> welcome to the party
sorry guys, g2g and ttyl
$
```

### User 3
```
C:\>.\nc.exe localhost 8080
< Welcome to my chat server! What's your name?
fei
< That nickname is already taken, please pick a different one:
bob
< You are connected with users: [fei, jake]
< [19:36:38] <fei> helloooo
< [19:36:49] *jake has joined the chat*
< [19:37:00] <jake> hello fei
< [19:37:10] <fei> what's up jake!
< [19:37:25] <fei> this is a cool coding challenge!
< [19:37:47] <fei> i can send you a beep too @jake
< [19:37:53] <jake> that's cool
< [19:38:38] <fei> sending more messages
< [19:38:42] <fei> sending more messages again
< [19:38:54] <fei> that should be more than 10 messages
hey guys
< [19:40:17] <fei> welcome to the party
< [19:40:31] <jake> sorry guys, g2g and ttyl
< [19:40:32] *jake has left the chat*
```
