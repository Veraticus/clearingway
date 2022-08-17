# clearingway

Clearingway listens to a Discord channel of your choosing for the message `!clears`.

When it hears this message, it examines the sending user's guild name (which should be their full Final Fantasy XIV character name, i.e. "Atmus Coldheart")
and what server they are on based on their roles (expecting a role that exactly matches the name of a NA server, i.e. "Gilgamesh").

Clearingway sends a few graphql requests to fflogs to determine what fights that character has cleared, and assigns the Discord user:
1. A role for the highest current parse they have in any relevant encounter (Gold, Orange, Purple, Blue, Green, Grey),
2. A role for every relevant encounter they've cleared ("P1S-Cleared," "P2S-Cleared," "P3S-Cleared," etc.)
3. A combo legend role purely for flexing purposes for every ultimate they've cleared (The Legend, The Double Legend, The Triple Legend, The Tetra Legend)

## Running

Clearingway requires the following environment variables to start:

* **DISCORD_TOKEN**: You have to create a [Discord bot for Clearingway](https://discord.com/developers/applications). Once you've done so, you can add the bot token here.
* **DISCORD_CHANNEL_ID**: The channel ID on which to listen for the `!clears` message.
* **FFLOGS_CLIENT_ID**: The client ID from [fflogs](https://www.fflogs.com/api/clients/).
* **FFLOGS_CLIENT_SECRET**: The client secret from [fflogs](https://www.fflogs.com/api/clients/).
* **ENCOUNTERS**: A string of `<name>=<encounter_id>`, separated by commas. For example, `P1S=78,P2S=79,P3S=80,P4SP1=81,P4SP2=82`. These are the only encounters the bot will consider relevant.
For each named encounter, the bot will create a `<name>-PF` and `<name>-Cleared` role.

