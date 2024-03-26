package config

import (
	"errors"
	"os"
)

const (
	EnvTesting    = Environment("testing")
	EnvProduction = Environment("production")
)

type Environment string

func (e *Environment) String() string {
	return string(*e)
}

func (e *Environment) Production() Environment {
	return EnvProduction
}

func (e *Environment) Testing() Environment {
	return EnvTesting
}

func (e *Environment) Invalid() bool {
	return *e != EnvTesting && *e != EnvProduction
}

func NewGlobalEnvironment() (Environment, error) {
	environment, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		return "", errors.New("system environment: ENVIRONMENT not found")
	}

	env := Environment(environment)
	if env != EnvTesting && env != EnvProduction {
		return "", errors.New("invalid environment")
	}
	return env, nil
}
