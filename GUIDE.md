# Guide

There are several common problems that happen during jams that need to be taken care of.

## Decide on target platforms.

Jam organizers need to decide on a target platform, otherwise you can end up with builds that runs only on a Windows Phone.

Usually a good targets are: Web and Windows, however it's highly recommended to get Mac and Linux builds.

## Jammers

The most important thing as a jammer, is that you are responsible such that other people can download and play your game.

### Registering a team

This should happen at the start of the jam, otherwise it'll get hectic at the end of the jam when there are other burning problems.

* Login to jamvote to ensure that they are present in the system.
* Register with your **real name**, such that the organizers can easily solve problems with users.
* Create a team and add all members to the team. If a person is not part of a team, they won't be able to vote.

### Production build

Ensure at the start of the jam that you can make a production build. There are plenty of things that can go wrong.

Test the production build on a system that does not have the development tools. Requiring people to install "tool X", usually means they won't be able to get the game running.

Compress production builds with `.zip`, since it's the only one that works across platforms. `7z`, `rar`, `tar`, `gzip` etc. require additional tools depending on the OS.

Upload to a public location. Jammers sometimes forget to change their sharing settings to allow downloading by unregistered people.

Good way is to let someone outside the team to test the download and run the "production build". _e.g. There have been issues where people have used some audio/images that triggered Google Drive-s automatic "copy protection"._

Start uploading and creating the final binary 1 hour before voting and presentations. There's plenty of things that can go wrong from compiling shaders, updating descriptions to uploading to a server.

When you are unable to get a working production build that other people can play, create and upload a video. This can be sufficient to get the game across and get some feedback on the game.

When you fix a bug, ensure you update all your uploads and links.

## Voting

Create a dedicated "emergency" communication channel where you can notify people that a game won't start. There is usually at least one per event.

You can vote for all games, however it'll randomly assign 3 to start with... and then one at a time. This is to ensure every game gets sufficient number of votes as fast as possible.

Decide what the "voting" approach is when a game doesn't start or there is only video. Do you give "0" points or do you let base their opinion on the video.

It's possible to change your votes on the "Voting" page by clicking "Edit" in front of the teams name. For example when they get a hotfix before the end of the jam.

Try to leave as much feedback as possible.

## Organizers

Make sure that teams, users are properly linked and approved for the jam. Linked means that the names on the teams page and in users match. Approved means that a user is able to vote in the jam. This requires some manual oversight and pestering jammers to login to the site and create their team.