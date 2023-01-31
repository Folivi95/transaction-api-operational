---
name: Story kick-off (KO) and definition of done (DOD)
about: Detailed information and updates on a story
title: ''
labels: story
assignees: ''

---

## What?

<!-- _What are you trying to achieve?_ What's the big picture at the end of the rainbow? Does this really help our user do their job? Now is a good time to decide. -->

### Example criteria

| Given | When | Then |
| ----- | ------ | ----- |
| Chris James is the king | He writes code | Everyone's dreams come true |

## How?

- How are you going to proceed? 
- What's the smallest change _that can be integrated_, and that gets you towards your final goal? Think about "steel thread" through the system. Could you maybe hard-code some values first to move forward?
- What's the next change after that? And after that? 
- Can you see a path made of small steps that get us to our destination? 

Write down your high-level path to getting this work done here. If you can't, maybe have a conversation with someone about how to proceed.

## How will you know it works, and it works well?

- What metrics will you need to look at to know it's working well?
- What logging do you need for when things go wrong?

# Definition of Done

## Testing 
- [ ] Unit
- [ ] Integration  
- [ ] Acceptance
- [ ] Performance

## Operability 

- [ ] Metrics
- [ ] Logging 
- [ ] Alerts 

## Documentation 

- [ ] Do you need an ADR?
- [ ] Have you added a new external dependency? If yes, add to the context diagram
- [ ] Have you introduced a new flow through the system? Consider adding a sequence diagram
- [ ] Have you added a new API? Document an example request
- [ ] Have you changed the way the system is built or requires developers to do extra setup? Document in the README.md
- [ ] Have you considered how this change may affect on-call? If so, update the runbook. 

## Internal quality

- [ ] Are you happy with the code?
- [ ] Is someone else happy with the code?

## Delivering value

- [ ] Is it live?
- [ ] Have you demonstrated the change to whoever needs it?

