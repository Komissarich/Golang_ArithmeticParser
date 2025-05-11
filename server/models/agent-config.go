package models

type AgentConfig struct {
	TimeAdditionMS       int `yaml:"TIME_ADDITION_MS"`
	TimeSubtractionMS    int `yaml:"TIME_SUBTRACTION_MS"`
	TimeMultiplicationMS int `yaml:"TIME_MULTIPLICATIONS_MS"`
	TimeDivisionMS       int `yaml:"TIME_DIVISIONS_MS"`
	ComputingPower       int `yaml:"COMPUTING_POWER"`
}
