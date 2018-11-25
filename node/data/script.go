package data

import (
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/version"
)

type Scripts struct {
	Name map[string]Versions
}

type Versions struct {
	Version map[string]Script
}

type Script struct {
	Script []byte
}

func init() {
	serial.Register(Scripts{})
	serial.Register(Versions{})
	serial.Register(Script{})
}

func NewScripts() *Scripts {
	return &Scripts{
		Name: make(map[string]Versions, 0),
	}
}

func (scripts *Scripts) Set(name string, version version.Version, script Script) {
	var versions Versions
	var ok bool

	if versions, ok = scripts.Name[name]; !ok {
		scripts.Name[name] = *NewVersions()
		versions = scripts.Name[name]
	}
	versions.Version[version.String()] = script
}

func NewVersions() *Versions {
	return &Versions{
		Version: make(map[string]Script, 0),
	}
}
