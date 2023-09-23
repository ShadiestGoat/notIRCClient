# Not IRC Client

## What is this?

Not an IRC client, let me tell you that much...

This is an 'educational' client for the [bare-bones server](https://github.com/ShadiestGoat/notIRCServer).

The idea is to create a non-authenticated server/client pair, where the client is a TUI/CUI.

For me personally, the bigger lesson was in syntax trees & markdown, and how it all works. This client has it's own markdown renderer, which taught me a lot of things.

## Quality

This project is a bit of a disorganized mess, since I was making the entire thing in a rush, and changing a lot all the time. Don't take this or the server as a representative of my quality of work T_T

## Configuration

This is a client made to not be configurable. Sounds *odd*, but the reason is that it's intended to be shipped to used with a few pre-baked values. As such, the build script (`build.sh`) is whats used to 'configure' the builds. Edit the build script to specify the location of your server, and run it with env for `AUTHOR_COLOR` (a 6-valued hex color, without any prefix like `0x` or `#`), and `AUTHOR_NAME` (the clients's username)
