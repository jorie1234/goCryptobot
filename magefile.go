// +build mage

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

// Compile and lint the cli
func Build() error {
	mg.Deps(Lint)

	return sh.RunV("go", "build", "-ldflags="+ldflags(), "./...")
}

// Nit the hell outta my code
func Lint() error {
	//mg.Deps(EnsureGoLint)

	return sh.RunV("golangci-lint.exe", "run", "./...")
}

func ldflags() string {
	timestamp := time.Now().Format(time.RFC822)
	hash := hash()
	tag := tag()
	if tag == "" {
		tag = "dev"
	}
	ver := IncVersion()
	return fmt.Sprintf(`-X "main.version=%s" `+
		`-X "main.commitHash=%s" `+
		`-X "main.date=%s"`, ver, hash, timestamp)
}

func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}

func tag() string {
	s, _ := sh.Output("git", "describe", "--tags")
	version, err := semver.Make(strings.ReplaceAll(s, "v", ""))
	if err != nil {
		fmt.Print(err)
	}
	version.IncrementPatch()

	return version.String()
}
func IncVersion() string {
	v := semver.MustParse(GetVersion())
	v.IncrementPatch()
	err := ioutil.WriteFile("VERSION", []byte(v.String()), 0600)
	if err != nil {
		fmt.Print(err)
	}
	return v.String()
}

func GetVersion() string {
	b, err := ioutil.ReadFile("VERSION")
	if err != nil {
		ioutil.WriteFile("VERSION", []byte("0.0.1"), 0600)
		return "0.0.1"
	}
	return string(b)
}

func Release() (err error) {
	if os.Getenv("TAG") == "" {
		return errors.New("TAG environment variable is required")
	}
	if err := sh.RunV("git", "tag", "-a", "$TAG"); err != nil {
		return err
	}
	if err := sh.RunV("git", "push", "origin", "$TAG"); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			sh.RunV("git", "tag", "--delete", "$TAG")
			sh.RunV("git", "push", "--delete", "origin", "$TAG")
		}
	}()
	//return retool("goreleaser")
	return nil
}
