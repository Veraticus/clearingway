package clearingway

type Config struct {
	ConfigGuilds []*ConfigGuild `yaml:"guilds"`
}

type ConfigGuild struct {
	Name             string             `yaml:"name"`
	GuildId          string             `yaml:"guildId"`
	ChannelId        string             `yaml:"channelId"`
	ConfigEncounters []*ConfigEncounter `yaml:"encounters"`
	ConfigRoles      *ConfigRoles       `yaml:"roles"`
}

type ConfigRoles struct {
	RelevantParsing bool `yaml:"relevantParsing"`
	RelevantFlexing bool `yaml:"relevantFlexing"`
	Legend          bool `yaml:"legend"`
	UltimateFlexing bool `yaml:"ultimateFlexing"`
	World           bool `yaml:"world"`
}

type ConfigEncounter struct {
	Ids          []int         `yaml:"ids"`
	Name         string        `yaml:"name"`
	Difficulty   string        `yaml:"difficulty"`
	DefaultRoles bool          `yaml:"defaultRoles"`
	ConfigRoles  []*ConfigRole `yaml:"roles"`
}

type ConfigRole struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Color int    `yaml:"color"`
}
