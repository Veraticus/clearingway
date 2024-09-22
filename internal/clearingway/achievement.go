package clearingway

type Achievement struct {
	Title string `yaml:"name"`
	Type  string
	Roles map[RoleType]*Role
}

type Achievements struct {
	Achievements []*Achievement
}

func (a *Achievements) Roles() *Roles {
	roles := &Roles{Roles: []*Role{}}

	for _, achievement := range a.Achievements {
		for _, role := range achievement.Roles {
			roles.Roles = append(roles.Roles, role)
		}
	}
	return roles
}

func (a *Achievement) Init(c *ConfigAchievement) {
	a.Title = c.Title
	a.Roles = map[RoleType]*Role{}

	for _, configRole := range c.ConfigRoles {
		roleType := RoleType(configRole.Type)
		role, ok := a.Roles[roleType]
		if !ok {
			role = &Role{Type: roleType}
			a.Roles[roleType] = role
		}
		if len(configRole.Name) != 0 {
			role.Name = configRole.Name
		}
		if len(configRole.Description) != 0 {
			role.Description = configRole.Description
		}
		if configRole.Color != 0 {
			role.Color = configRole.Color
		}
	}
}
