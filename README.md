# clearingway

Clearingway listens to a Discord channel of your choosing for the message `/verify <world> <first-name> <last-name>`.

When it hears this message, it tries to find the relevant character in the Lodestone. If it finds them, it then parses their
fflogs and tries to assign them a few roles:

1. A role for the highest current parse they have in any relevant encounter (Gold, Orange, Purple, Blue, Green, Grey),
2. A role for every relevant encounter they've cleared ("P1S-Cleared," "P2S-Cleared," "P3S-Cleared," etc.)
3. A combo legend role purely for flexing purposes for every ultimate they've cleared (The Legend, The Double Legend, The Triple Legend, The Tetra Legend)

## Running

Clearingway requires the following environment variables to start:

* **DISCORD_TOKEN**: You have to create a [Discord bot for Clearingway](https://discord.com/developers/applications). Once you've done so, you can add the bot token here.
* **DISCORD_CHANNEL_ID**: The channel ID on which to listen for the `!clears` message.
* **FFLOGS_CLIENT_ID**: The client ID from [fflogs](https://www.fflogs.com/api/clients/).
* **FFLOGS_CLIENT_SECRET**: The client secret from [fflogs](https://www.fflogs.com/api/clients/).

