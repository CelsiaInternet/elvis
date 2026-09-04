package envar

type ConfigStore struct {
	projectId  string
	packegName string
	stage      string
}

func NewConfigStore(projectId string, packegName string, stage string) *ConfigStore {
	return &ConfigStore{
		projectId:  projectId,
		packegName: packegName,
		stage:      stage,
	}
}

func (c *ConfigStore) Get(name string, _default string) string {
	return ""
}

func (c *ConfigStore) SetConfig(name string, value string) {

}
