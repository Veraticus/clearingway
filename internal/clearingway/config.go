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
	ConfigAchievements        []*ConfigAchievement        `yaml:"achievements"`
	ConfigRoles               *ConfigRoles                `yaml:"roles"`
	ConfigReconfigureRoles    []*ConfigReconfigureRoles   `yaml:"reconfigureRoles"`
	ConfigMenus               []*ConfigMenu               `yaml:"menu"`
}

type ConfigRoles struct {
	RelevantParsing    bool `yaml:"relevantParsing"`
	RelevantFlexing    bool `yaml:"relevantFlexing"`
	RelevantRepetition bool `yaml:"relevantRepetition"`
	Legend             bool `yaml:"legend"`
	UltimateFlexing    bool `yaml:"ultimateFlexing"`
	UltimateRepetition bool `yaml:"ultimateRepetition"`
	Datacenter         bool `yaml:"datacenter"`
	SkipRemoval        bool `yaml:"skipRemoval"`
	NameColor          bool `yaml:"nameColor"`
	Reclear            bool `yaml:"reclear"`
	Menu               bool `yaml:"menu"`
}

type ConfigEncounter struct {
	Ids                   []int         `yaml:"ids"`
	Name                  string        `yaml:"name"`
	Difficulty            string        `yaml:"difficulty"`
	DefaultRoles          bool          `yaml:"defaultRoles"`
	TotalWeaponsAvailable int           `yaml:"totalWeaponsAvailable"`
	The                   string        `yaml:"the"`
	ConfigRoles           []*ConfigRole `yaml:"roles"`
	ConfigProg            []*ConfigRole `yaml:"prog"`
	RequiredKillsToClear  int           `yaml:"requiredKillsToClear"`
}

type ConfigAchievement struct {
	Title       string        `yaml:"name"`
	ConfigRoles []*ConfigRole `yaml:"roles"`
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
	Hoist bool   `yaml:"hoist"`
}

type ConfigReconfigureRoles struct {
	Type          string `yaml:"type"`
	EncounterName string `yaml:"encounterName"`
	From          string `yaml:"from"`
	To            string `yaml:"to"`
	Color         int    `yaml:"color"`
	Skip          bool   `yaml:"skip"`
	DontSkip      bool   `yaml:"dontSkip"`
	Hoist         bool   `yaml:"hoist"`
}

type ConfigMenu struct {
	Name          string          `yaml:"name"`
	Type          string          `yaml:"type"`
	Title         string          `yaml:"title"`
	Description   string          `yaml:"description"`
	ImageUrl      string          `yaml:"imageUrl"`
	ConfigButtons []*ConfigButton `yaml:"buttons"`
	ConfigRoles   []*ConfigRole   `yaml:"roles"`
	RoleType      []string        `yaml:"roleType"`
	MultiSelect   bool            `yaml:"multiSelect"`
	RequireClear  bool            `yaml:"requireClear"`
}

type ConfigButton struct {
	Label    string `yaml:"label"`
	Style    int    `yaml:"style"`
	MenuName string `yaml:"menuName"`
	MenuType string `yaml:"menuType"`
}
