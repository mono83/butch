package butch

import (
	"io"
	"os/exec"
	"path/filepath"
)

// ShellOptions contains CLI exec options
type ShellOptions struct {
	Path    string
	Command string
	Args    []string
	Done    func(worker Worker, p *Pool)
}

// AddShellWorker registers shell worker
func (p *Pool) AddShellWorker(opt ShellOptions) {
	w := &Worker{
		Name: filepath.Base(opt.Path) + " " + opt.Command,
		Run: func(p *Pool, w Worker, stdout, stderr io.Writer) {
			// Building command
			cmd := exec.Command(opt.Command, opt.Args...)

			// Setting path
			if len(opt.Path) > 0 {
				cmd.Dir = opt.Path
			}

			// Setting stdout & stderr
			cmd.Stdout = stdout
			cmd.Stderr = stderr

			// Running
			cmd.Run()
		},
	}

	p.AddWorker(w)
}
