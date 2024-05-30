# My telegram bot to get homelab notification

## Goal

After some issue with my Raspberry Pi, AKA Cpu throttle, I wanted a way to get notification and metrics of my servers of my homelab.  
There is some way to have it with well known softwares. But it is funnier to build it from scratch.  

Here is the plan:  
Create a Master-Agent pattern communicating by gRPC in golang. Why ? Because I know none of those subjects.

## Plan

- [ ] Create the agent with the gRPC communication, may test with python scripts.
   > Note that these agents will be **gRPC servers** and the master will send requests to them. 
- [ ] Connect agent to the master so he can know who is connected. (gave custom name or hostname)
- [ ] Create the bot commands via the master.
- [ ] System to have command with steps: /metrics -> which server ? Master send new buttons to select by the user, just like botFather with bot management.

- [ ] How to easily manage commands/function associated ? map ? objects ?
