gomodule := 'github.com/ircurry/dfh'
package := 'default'
cmds := ```
    dirs="$(find ./cmd/ -mindepth 1 -maxdepth 1 -type d)"
    for dir in "$dirs"; do basename "$dir"; done
	```

# run 'just gbuild dfh'
[group('build')]
default: (gobuild)

# format 'cmd' files using go fmt
[group('format')]
goformatcmd cmd: 
	go fmt {{ gomodule }}/cmd/{{ cmd }}
alias gfcmd := goformatcmd

# format all go command files using go fmt
[group('format')]
goformat: 
	for x in {{ cmds }}; do go fmt {{ gomodule }}/cmd/$x; done
alias gf := goformat

# test 'cmd' files using go test
[group('test')]
gotestpkg pkg:
	go test {{ pkg }}
alias gtpkg := gotestpkg

# test all go command files using go test
[group('test')]
gotest:
	go test ./...
alias gt := gotest

# build 'cmd' using go build
[group('build')]
gobuildcmd cmd:
	go build {{ gomodule }}/cmd/{{ cmd }}
alias gbcmd := gobuildcmd

# build all go binaries in cmd/ directory
[group('build')]
gobuild:
	for x in {{ cmds }}; do go build {{ gomodule }}/cmd/$x; done
alias gb := gobuild

# build 'package' using nix
[group('build')]
build:
	nix build ".#"{{ package }}
alias b := build
