package clearingway

type Config struct {
	ConfigGuilds []*ConfigGuild `yaml:"guilds"`
}

type ConfigGuild struct {
	Name                   string                    `yaml:"name"`
	GuildId                string                    `yaml:"guildId"`
	ChannelId              string                    `yaml:"channelId"`
	ConfigDatacenters      []*ConfigDatacenter       `yaml:"datacenters"`
	ConfigEncounters       []*ConfigEncounter        `yaml:"encounters"`
	ConfigRoles            *ConfigRoles              `yaml:"roles"`
	ConfigReconfigureRoles []*ConfigReconfigureRoles `yaml:"reconfigureRoles"`
}

type ConfigRoles struct {
	RelevantParsing bool `yaml:"relevantParsing"`
	RelevantFlexing bool `yaml:"relevantFlexing"`
	Legend          bool `yaml:"legend"`
	UltimateFlexing bool `yaml:"ultimateFlexing"`
	Datacenter      bool `yaml:"datacenter"`
}

type ConfigEncounter struct {
	Ids          []int         `yaml:"ids"`
	Name         string        `yaml:"name"`
	Difficulty   string        `yaml:"difficulty"`
	DefaultRoles bool          `yaml:"defaultRoles"`
	ConfigRoles  []*ConfigRole `yaml:"roles"`
}

type ConfigRole struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Color       int    `yaml:"color"`
	Description string `yaml:"description"`
}

type ConfigDatacenter struct {
	Name       string `yaml:"name"`
	Datacenter string `yaml:"datacenter"`
	Color      int    `yaml:"color"`
}

type ConfigReconfigureRoles struct {
	From  string `yaml:"from"`
	To    string `yaml:"to"`
	Color int    `yaml:"color"`
}
