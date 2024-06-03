# My telegram bot to get homelab notification

## Goal

After some issue with my Raspberry Pi, AKA Cpu throttle, I wanted a way to get notification and metrics of my servers of my homelab.  
There is some way to have it with well known softwares. But it is funnier to build it from scratch.  

Here is the plan:  
Create a Master-Agent pattern communicating by gRPC in golang. Why ? Because I know none of those subjects.

## Plan

**gRPC**:
 - [ ] Create the agent with the gRPC communication, may test with python scripts.
   > Note that these agents will be **gRPC servers** and the master will send requests to them. 
 - [ ] Connect agent to the master so he can know who is connected. (gave custom name or hostname)

**BOT**:
 - [ ] How to easily manage commands/function associated ? map ? struct ?
 - [ ] System to have command with steps: /metrics -> which server ? Master send new buttons to select by the user, just like botFather with bot management.

- [ ] Create/match the bot commands via the master to grpc communication.

## Interactions

The combinations of menu button (a list of all commands available) and of keyboard button when multiple choices are possible. The usage of CallbackQueries is not meant to respond something in the chat (according to telegram doc), but we still have the corresponding message so it's still possible to respond something! Great. With this InlineKeyboard with their callback maybe the best way to chaine messages.

One way I can think about is:
 - commands, available by menu button
 - response with InlineButton, with specifique Callback (only 64 char long)
