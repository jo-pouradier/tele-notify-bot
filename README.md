# My telegram bot to get homelab notification

## Goal

After some issue with my Raspberry Pi, AKA Cpu throttle, I wanted a way to get notification and metrics of my servers of my homelab.  
There is some way to have it with well known softwares. But it is funnier to build it from scratch.  

Here is the plan:  
Create a Master-Agent pattern communicating by gRPC in golang. Why ? Because I know none of those subjects.

## Plan

1. Create the agent with the gRPC communication, may test with python scripts.
    Note that these agents will be **gRPC servers** and the master will send requests to them. 
