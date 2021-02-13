# arkbot

arkbot is a small bot to monitor a set of [ARK](https://survivetheark.com)
servers and report on which of them are online, what map they are serving, and
how many players are online.

## Running the bot

The bot is most easily run as a docker container. Configure your bot token,
channel id, and list of servers as environment variables.

```bash
docker run -d \
    -e ARKBOT_TOKEN="XXXXXXXXXXXXXXXXXXX" \
    -e ARKBOT_CHANNEL="01234567890" \
    -e ARKBOT_SERVERS="my.server.com:12345,my.other.server.com:56789" \
    --name=arkbot justinian/arkbot
```

## Thanks

This repository mostly just steals from [my discord dice bot][2] and [the
discordgo airhorn example][3]. Many thanks to bwmarrin for Discordgo and the
great examples.

[2]: https://github.com/justinian/discorddice
[3]: https://github.com/bwmarrin/discordgo/tree/master/examples/airhorn
