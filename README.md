# My telegram bot to get homelab notification

Part of this project are higly inspired from [grpc-go examples](https://github.com/grpc/grpc-go/tree/master/examples)


## Goal

After some issue with my Raspberry Pi, AKA Cpu throttling, I wanted a way to get notification and metrics of my servers of my homelab.  
There is some way to have it with well known softwares. But it is funnier to build it from scratch.  

Here is the plan:  
Create a Master-Agent pattern communicating by gRPC in golang. Why ? Because I know none of those subjects.

## Plan

**gRPC**:
 - [x] Add Tls
 - [x] Add authentication
 - [x] Create the agent with the gRPC communication, may test with python scripts.
   > Note that these agents will be **gRPC servers** and the master will send requests to them [See Update on this](#architecture-issue). 
 - [ ] Connect agent to the master so he can know who is connected via metadata.

**BOT**:
 - [ ] How to easily manage commands/function associated ? map ? struct ?
 - [ ] System to have command with steps: /metrics -> which server ? Master send new buttons to select by the user, just like botFather with bot management.

- [ ] Create/match the bot commands via the master to grpc communication.

## Interactions

The combinations of menu button (a list of all commands available) and of keyboard button when multiple choices are possible. The usage of CallbackQueries is not meant to respond something in the chat (according to telegram doc), but we still have the corresponding message so it's still possible to respond something! Great. With this InlineKeyboard with their callback maybe the best way to chaine messages.

One way I can think about is:
 - commands, available by menu button
 - response with InlineButton, with specifique Callback (only 64 char long)

## Registration

We ask for a new key to the master, give this key to the agent, the agent connect to themaster with some metadata like server name, version, status(?), methods(?). The master accept the connection based on the token and register this agent based on the metadata.

## Architecture Issue

The agent beeing a server implise that one port must be openned and that the server is accessible on this port, the agent should have his own TLS keys + more work on master to register and manage each agent.  
Opening a port is a problem if we want to manage servers in a private network and that the master cant be on this same network. Only the master should be a server and the agent initiate a long pooling request or a bidirectional stream (I'll go for the stream).
