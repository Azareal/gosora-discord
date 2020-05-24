# Discord Push for Gosora

A plugin for pushing to Discord channels when a new reply or topic is created.

# Requirements

A recent build of [Gosora](https://github.com/Azareal/Gosora).

# Installation

You can install this plugin simply by throwing it in Gosora's /extend/ folder before building / running it.

# Options

You can configure this plugin via your [config.json file](https://github.com/Azareal/Gosora/blob/master/docs/configuration.md).

DiscordWebhook - The [webhook URL](https://support.discordapp.com/hc/en-us/articles/228383668-Intro-to-Webhooks) for the Discord channel you're trying to push events to. This is required for this plugin to function.

DiscordForums - A comma separated list of IDs of the forum (or forums) you want this plugin to listen to. Defaults to all of them.

DiscordEvents - The events you want this plugin to listen to. You can set this to `topics` to only have it listen to new topics being created and to `replies` to only listen to new replies being created. For both, you can just leave this blank. Default: Both.

# Future Features

I'd like to add more Discord integration features like syncing Discord ranks with forum ranks and maybe Discord accounts with user accounts.