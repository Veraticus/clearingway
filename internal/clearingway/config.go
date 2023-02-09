package clearingway

type Config struct {
	ConfigGuilds []*ConfigGuild `yaml:"guilds"`
}

type ConfigGuild struct {
	Name                      string                      `yaml:"name"`
	GuildId                   string                      `yaml:"guildId"`
	ChannelId                 string                      `yaml:"channelId"`
	ConfigPhysicalDatacenters []*ConfigPhysicalDatacenter `yaml:"physicalDatacenters"`
	ConfigEncounters          []*ConfigEncounter          `yaml:"encounters"`
	ConfigRoles               *ConfigRoles                `yaml:"roles"`
	ConfigReconfigureRoles    []*ConfigReconfigureRoles   `yaml:"reconfigureRoles"`
}

type ConfigRoles struct {
	RelevantParsing bool `yaml:"relevantParsing"`
	RelevantFlexing bool `yaml:"relevantFlexing"`
	Legend          bool `yaml:"legend"`
	UltimateFlexing bool `yaml:"ultimateFlexing"`
	Datacenter      bool `yaml:"datacenter"`
	SkipRemoval     bool `yaml:"skipRemoval"`
}

type ConfigEncounter struct {
	Ids          []int         `yaml:"ids"`
	Name         string        `yaml:"name"`
	Difficulty   string        `yaml:"difficulty"`
	DefaultRoles bool          `yaml:"defaultRoles"`
	ConfigRoles  []*ConfigRole `yaml:"roles"`
	ConfigProg   []*ConfigRole `yaml:"prog"`
}

type ConfigRole struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Color       int    `yaml:"color"`
	Hoist       bool   `yaml:"hoist"`
	Mention     bool   `yaml:"mention"`
	Description string `yaml:"description"`
}

type ConfigPhysicalDatacenter struct {
	Name               string                     `yaml:"name"`
	LogicalDatacenters []*ConfigLogicalDatacenter `yaml:"logicalDatacenters"`
}

type ConfigLogicalDatacenter struct {
	From  string `yaml:"from"`
	To    string `yaml:"to"`
	Color int    `yaml:"color"`
}

type ConfigReconfigureRoles struct {
	From     string `yaml:"from"`
	To       string `yaml:"to"`
	Color    int    `yaml:"color"`
	Skip     bool   `yaml:"skip"`
	DontSkip bool   `yaml:"dontSkip"`
}
