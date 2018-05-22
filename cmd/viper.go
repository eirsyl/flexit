package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func uppercaseName(name string) string {
	s := camelcase.Split(name)
	snake := strings.Join(s, "_")
	return strings.ToUpper(snake)
}

// StringConfig adds a string flag to a cli
func StringConfig(cmd *cobra.Command, name, short, value, description string) {
	cmd.PersistentFlags().StringP(name, short, value, description)
	err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
	if err != nil {
		panic(err)
	}
	err = viper.BindEnv(name, uppercaseName(name))
	if err != nil {
		panic(err)
	}
}

// IntConfig adds a string flag to a cli
func IntConfig(cmd *cobra.Command, name, short string, value int, description string) {
	cmd.PersistentFlags().IntP(name, short, value, description)
	err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
	if err != nil {
		panic(err)
	}
	err = viper.BindEnv(name, uppercaseName(name))
	if err != nil {
		panic(err)
	}
}

// BoolConfig adds a bool flag to a cli
func BoolConfig(cmd *cobra.Command, name, short string, value bool, description string) {
	cmd.PersistentFlags().BoolP(name, short, value, description)
	err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
	if err != nil {
		panic(err)
	}
	err = viper.BindEnv(name, uppercaseName(name))
	if err != nil {
		panic(err)
	}
}

type FlagChecker func() error

func CheckFlags(checkers ...FlagChecker) {
	var fails []string
	for _, checker := range checkers {
		if err := checker(); err != nil {
			fails = append(fails, err.Error())
		}
	}
	if len(fails) > 0 {
		fmt.Println(strings.Join(fails, "\n"))
		os.Exit(1)
	}
}

// RequireString returns an error if the given setting is not a string
func RequireString(flag string) FlagChecker {
	return func() error {
		v := viper.GetString(flag)
		if v == "" {
			return fmt.Errorf("flag %s can not be an empty string", flag)
		}
		return nil
	}
}
