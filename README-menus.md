# Overview
This README's purpose is to guide users through setting up Discord GUI menus with Clearingway.  
All things discussed here are only relevant to `config.yaml`.
# Menus
There are currently 4 types of menus:
- `menuMain`
- `menuVerify`
- `menuRemove`
- `menuEncounter`

If the menus module is enabled, there will be 3 menus that will be configured by default, all named by their namesake types: `menuMain`, `menuVerify`, `menuRemove`. 
## Common
All menus will be able to utilize these properties.  
#### Required properties
- `name` (string) - The internal name of the menu. Using the aforementioned default menu names will allow you to overwrite defaults for the title, description, etc. Otherwise, it will create a new menu.  
- `type` (string) - The menu type with 4 possible values
    - `menuMain`
    - `menuVerify`
    - `menuRemove`
    - `menuEncounter`

For creating new menus, the only types to be concerned with are `menuMain` and `menuEncounter`. `menuVerify` and `menuRemove` should only be used to override the default text/style properties.

##### Optional properties
- `title` (string) - The title of the embed
- `description` (string) - The description of the embed
- `imageUrl` (string) - A full-size image to be displayed at the bottom of the embed
- `thumbnailUrl` (string) - A thumbnail embed to be displayed at the top-right of the embed
- `fields` (array of object) - Fields to be displayed within an embed  
  - `name` (string) - Title of field
  - `value` (string) - Description of field
  - `inline` (bool) - Whether field is inline or not

## MenuMain
Menus that are sendable to the channel which are viewable by every user
#### Optional properties
##### `buttons` (array of object)
A list of buttons to be displayed under the embed. Each button object has the following required properties:
- `label` (string) - Button label
- `style` (int) - Button color value, follows `discordgo.ButtonStyle` enums
  - `1` = Primary (Blurple)
  - `2` = Secondary (Grey)
  - `3` = Success (Green)
  - `4` = Danger (Red)
  - `5` = URL (Grey), prob shouldn't be used?
- `menuName` (string) - Name of menu button leads to
- `menuType` (string) - Type of menu referred to by `menuName`


## MenuEncounter
Ephemeral response menus that allow the user to pick role(s) through a dropdown selection. Dropdown can be single option or multi option. 

Roles in this menu can be set up 2 ways, for all ultimate encounters or a per-role basis. The former allows for requiring a cleared role for the respective ultimate. The latter allows to add miscellaneous roles that are connected to Clearingway's internal encounters.

### Common
#### Optional properties
- `multiSelect` (bool) - Whether multiple roles can be selected from this category

### All ultimate encounters
Roles must still be set up separately and connected the respective encounter in the guild's config.
#### Required properties
- `roleType` (array of string) - Category of role (string enum), defined in `role.go`
#### Optional properties
- `requireClear` (bool) - Whether the configured roles require their respective encounter's clear role.

### Per-role basis
#### Required properties
- `roles` (array of object) same object as Clearingway roles. Only `name` is required.
    - Name (string) - Name of role
    - Type (string) - Category of role (string enum), same as above
    - Color (int) - Color's hex code
    - Hoist (bool) - Whether the role is displayed separately in the userlist
    - Mention - Whether the role is mentionable
    - Description - Internal Clearingway description of the role

# Examples
## Overriding title and description for `menuMain`
```yaml
guilds:
- name: guild name
  roles:
    menu: true
  guildId: 1234567891234567891
  channelId: 1234567891234567891
  menu:
  - name: "menuMain"
    type: "menuMain"
    title: "Example"
    description: |-
      This
      is
      a
      multiline
      description
```
For basic styling, we can add a thumbnail, image, and fields to the embed
```yaml
# ...
  menu:
    - name: "menuMain"
      type: "menuMain"
      title: "Example"
      description: "Example"
      imageUrl: https://raw.githubusercontent.com/naurffxiv/assets/refs/heads/main/Discord%20Files/Misc%20Images/  ultimates%20colour%20server%20header.png
      thumbnailUrl: https://github.com/naurffxiv/assets/blob/main/Discord%20Files/Misc%20Images/naur_icon_optimized2.gif?  raw=true
      fields:
        - name: "field title"
          value: "field description"
          # inline can be omitted since the default is false which is probably what you want
```
## Adding buttons to `menuMain`
This example adds a verify character button, a remove roles button, and a basic `menuEncounter` menu called `example` along with its respective button. The `example` menu will contain a dropdown for a single role that is not tied to any encounters.
```yaml
# ...
  menu:
    - name: "menuMain"
      type: "menuMain"
      title: "Example"
      description: "This is a single line description"
      buttons:
      # verify character button
      - label: "Verify Character"
        style: 3  # green button
        menuName: "menuVerify"
        menuType: "menuVerify"
      # remove roles button
      - label: "Remove Roles"
        style: 4 # red button
        menuName: "menuRemove"
        menuType: "menuRemove"
      # example menu button
      - label: "menuEncounter Button"
        style: 1  # blurple button
        menuName: "example"
        menuType: "menuEncounter"
      # example menu
    - name: "example"
      type: "menuEncounter"
      title: "Example submenu"
      description: "Example description"
      roles:
        - name: "Example Role"
```

And if you want to add multiple roles that are not tied to any encounters
```yaml
    # ...
    roles:
      - name: "Example Role 1"
      - name: "Example Role 2"
      - name: "Example Role 3"
    # As many roles as you want
```
If you want to let the user select multiple of these roles on top of that
```yaml
    # ...
    - name: "example"
      type: "menuEncounter"
      title: "Example submenu"
      description: "Example description"
      multiSelect: true
      roles:
        - name: "Example Role 1"
    # And so forth
```
## Creating encounter tied roles
Use these if you want to tie these roles to Clearingway's internal encounters. Usually, you'd want to use this if you want to check the user for the encounter's respective cleared role.

Roles here need to be set up in the encounters section instead of within the menu section.

Below is an example of a reclear submenu with 1 encounter, only allowing for 1 option to be selected at a time
```yaml
guilds:
  - name: guild name
    # ...
    encounters:
      # test encounter 1
      - ids: [1234]
        name: "Test encounter 1"
        roles:
        - name: "Test encounter 1 cleared"
          type: "Cleared"
        - name: "Test encounter 1 reclears"
          type: "Reclear"
    menu:
      # ... menuMain with the respective button ...
      # reclear submenu
      - name: "reclear"
        type: "menuEncounter"
        title: "Reclear Roles"
        description: "Reclear roles"
        roleType:
        - "Reclear"
        requireClear: true
```
Here I add 1 more encounter and another role type, now I have 2 different selectable role types in the submenu. With the addition of `multiSelect`, the user is able to pick any combination of these 4 selectable roles (Reclear/C4X roles for the 2 test encounters)
```yaml
# ...
  encounters:
  # ... test encounter 1 ...
  # test encounter 2
    - ids: [2345]
      name: "Test encounter 2"
      roles:
      - name: "Test encounter 2 cleared"
        type: "Cleared"
      - name: "Test encounter 2 reclears"
        type: "Reclear"
      - name: "Test encounter 2 C4X"
        type: "C4X"
  menu:
    # ... menuMain with the respective button ...
    # reclear/c4x submenu
    - name: "cleared"
      type: "menuEncounter"
      title: "Reclear/C4X Roles"
      description: "Reclear/C4X roles"
      roleType:
        - "Reclear"
        - "C4X"
      requireClear: true
      multiSelect: true
```
