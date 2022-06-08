# how - info repository tool
`how` is a command that saves information grouped by topic so that the information can be looked up later. 

For example, I can use `how -- write bluetooth 'Lets computers talk to one another'` to create a new topic entry in 
the database. Later, when I want to remember what bluetooth does, I can view the entry with `how bluetooth`:
```
$ how bluetooth
how information:

Name: bluetooth
Description: Lets computers talk to one another
```
I can also do a fuzzy search on the topic names, i.e. `how blue` will point you toward `how bluetooth`:
```
$ how blue
...

Search Result #1
Name: bluetooth
Description: Lets computers talk ...

...
```
