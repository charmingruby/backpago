package configs

func LoadConfigs() error {
	if err := LoadEnv(); err != nil {
		return err
	}

	return nil
}
