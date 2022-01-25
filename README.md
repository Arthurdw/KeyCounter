# Key Counter

I randomly had the great idea to make a custom keyboard layout which should be ideal for me. (way more optimized for programming/vim)
But to do this I need data, lots and lots of it. Since I want to make it for myself I had to make a little efficient background script that would silently collect the amount of times I type certain characters so I can place these in the most ideal locations on my layout.

## What is this intended to do?

It has one very simple task, collect all my keystrokes efficiently in the background

## Installation

This is intended to run as a background job. *(I've done it through the systemctl daemon)
Installation just requires a build of this repo, and then for the main file to be moved to `/bin`.

### Further installation with systemctl

If you want to *(like me)* run this as a background task constantly as a service. You can easily do this with a few more commands.
First of all copy the `key-counterd.service` to `/etc/systemd/system`. Then execute the following commands:
```bash
$ sudo systemctl daemon reload
$ sudo systemctl start key-counterd.service
$ sudo systemctl status key-counterd.service
$ sudo systemctl enable key-counterd.service
```

## Where is my data saved?

The data is saved in `/etc/key-counter/`, there are two files `data.csv` *(Which contains the general keys data)* and the `combinations.csv` file *(Which contains all keyboard combinations, eg CTRL+C)*.
