# Organizing a jam

Here is a list of things to consider while running a game jam. We
include recommendations and reasoning, but as situations vary you are
free to ignore them at your own peril.

## Set up emergency channels

Things will go wrong. Set up a communication channel (e.g. a Discord
server or group chat) where participants can report issues and
organizers can post announcements during the jam.

## Jammers should register early

There's so much chaos at the end of the jam so it's better to get the
bureaucracy out of the way early. Have everyone sign in to jamvote,
create their teams, and add members at the start of the jam. This way
people won't have to worry about registering at the last minute.

## Decide on target platforms

Announce a target platform before the jam starts, otherwise you can end
up with builds that only run on a Windows Phone.

Usually good targets are: Web and Windows, however it's highly
recommended to also support Mac and Linux builds.

## Timeline reminders

Remind people early and often about the deadlines. People are tired,
stressed, and focused on their projects; so time can slip out of their
mind. Key moments to send reminders:

- When the jam starts (confirm theme, deadline, and submission process)
- Halfway through
- A few hours before the deadline
- 1 hour before the deadline (start uploading now!)

## Production build

The production build is one of the biggest sources of last-minute
problems. Cover these points with your jammers:

**Build early.** Remind people at the start of the jam that they should
make a production build. There are plenty of things that can go wrong
during that stage.

**Test on a clean machine.** The build should work on a system without
the development tools installed. Requiring people to install "tool X"
usually means they won't be able to get the game running.

**Use `.zip` for compression.** It's the only format that works across
all platforms without extra tools. `7z`, `rar`, `tar`, `gzip` etc.
require additional software depending on the OS.

**Upload to a public location.** Jammers often forget to change their
sharing settings. The best way to verify is to have someone outside the
team test the download link. _There have been issues where people have
used audio/images that triggered Google Drive's automatic "copy
protection"._

**Start uploading early.** Remind people to start uploading at least 1
hour before the deadline. There are plenty of things that can go wrong
from compiling shaders, updating descriptions, to uploading to a server.

## When a game is broken

Have clear guidelines for what happens when a game is not playable due
to a technical issue. Options include:

- Allow a hotfix window (e.g. 30 minutes after the deadline)
- Let the team show a video of the game instead
- Have reviewers score based on whatever they can observe

## When a game is unfinished

Have guidelines for what happens when a game is unfinished. There are
still ways for people to show off what they completed:

- Create a video of what was completed — this still allows feedback
  from voting
- Skip voting on the entry and let the team do a live presentation
  instead
