package geenee

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestGeeneeError_Error(t *testing.T) {
	t.Run("validate error string", func(t *testing.T) {
		want := "heyo error"
		got := GeeneeError("heyo error")
		if want != got.Error() {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}

func Test_NewInterface(t *testing.T) {
	t.Run("validate new", func(t *testing.T) {
		wantName := "tester"
		wantVersion := "0.0.0"
		got := NewInterface(wantName, wantVersion, false)
		if got == nil {
			t.Fatal("want interface, got nil")
		}

		if got.Name != wantName {
			t.Errorf("want %s, got %s", wantName, got.Name)
		}
		if got.Version != wantVersion {
			t.Errorf("want %s, got %s", wantVersion, got.Version)
		}
		if got.Out != os.Stdout {
			t.Errorf("want %v, got %v", os.Stdout, got.Out)
		}
		if got.Err != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.Err)
		}

		if got.RootCommand == nil {
			t.Fatal("want root command, got nil")
		}

		if got.RootCommand.Name != wantName {
			t.Errorf("want %s, got %s", wantName, got.RootCommand.Name)
		}
		if got.RootCommand.Flags.Name() != wantName {
			t.Errorf("want %s, got %s", wantName, got.RootCommand.Flags.Name())
		}
		if got.RootCommand.Out != os.Stdout {
			t.Errorf("want %v, got %v", os.Stdout, got.RootCommand.Out)
		}
		if got.RootCommand.Err != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.RootCommand.Err)
		}
		if got.RootCommand.Usage == nil {
			t.Error("want DefaultCommandUsageFunc, got nil")
		}
	})
}

func TestCommandInterface_SetWriters(t *testing.T) {
	t.Run("validate new", func(t *testing.T) {
		b := bytes.NewBufferString("")
		got := NewInterface("tester", "0.0.0", false)
		got.RootCommand.SubCommands = []*Command{
			NewCommand("heyo", false),
		}

		if got.Out != os.Stdout {
			t.Errorf("want %v, got %v", os.Stdout, got.Out)
		}
		if got.Err != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.Err)
		}

		if got.RootCommand.Out != os.Stdout {
			t.Errorf("want %v, got %v", os.Stdout, got.RootCommand.Out)
		}
		if got.RootCommand.Err != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.RootCommand.Err)
		}

		if got.RootCommand.Flags.Output() != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.RootCommand.Flags.Output())
		}

		if got.RootCommand.SubCommands[0].Out != os.Stdout {
			t.Errorf("want %v, got %v", os.Stdout, got.RootCommand.SubCommands[0].Out)
		}
		if got.RootCommand.SubCommands[0].Err != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.RootCommand.SubCommands[0].Err)
		}

		if got.RootCommand.SubCommands[0].Flags.Output() != os.Stderr {
			t.Errorf("want %v, got %v", os.Stderr, got.RootCommand.SubCommands[0].Flags.Output())
		}

		got.SetWriters(b, b)

		if got.Out != b {
			t.Errorf("want %v, got %v", b, got.Out)
		}
		if got.Err != b {
			t.Errorf("want %v, got %v", b, got.Err)
		}

		if got.RootCommand.Out != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Out)
		}
		if got.RootCommand.Err != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Err)
		}

		if got.RootCommand.Flags.Output() != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Flags.Output())
		}

		if got.RootCommand.SubCommands[0].Out != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Out)
		}
		if got.RootCommand.SubCommands[0].Err != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Err)
		}

		if got.RootCommand.SubCommands[0].Flags.Output() != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Flags.Output())
		}

		got.SetWriters(nil, nil)
		if got.Out != b {
			t.Errorf("want %v, got %v", b, got.Out)
		}
		if got.Err != b {
			t.Errorf("want %v, got %v", b, got.Err)
		}

		if got.RootCommand.Out != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Out)
		}
		if got.RootCommand.Err != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Err)
		}

		if got.RootCommand.Flags.Output() != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.Flags.Output())
		}

		if got.RootCommand.SubCommands[0].Out != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Out)
		}
		if got.RootCommand.SubCommands[0].Err != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Err)
		}

		if got.RootCommand.SubCommands[0].Flags.Output() != b {
			t.Errorf("want %v, got %v", b, got.RootCommand.SubCommands[0].Flags.Output())
		}
	})
}

func TestCommandInterface_Exec(t *testing.T) {
	t.Run("validate interface returns version if version flag present, and doesn't run root command", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "test 0.0.0\n"
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
			},
			Out:             b,
			Err:             b,
			Version:         "test 0.0.0",
			MaxCommandDepth: 3,
		}
		err := subject.Exec([]string{"test", "-version"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate interface root command runs - no flags/args", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "root command ran"
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
			},
			Out:             b,
			Err:             b,
			MaxCommandDepth: 3,
		}
		err := subject.Exec([]string{"test"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate interface root command runs - no flags with args", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "root command ran"
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
			},
			Out:             b,
			Err:             b,
			MaxCommandDepth: 3,
		}
		err := subject.Exec([]string{"test", "notaflag", "notasubcommand"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate interface root command runs - no flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "root command ran"
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
			},
			Out:             b,
			Err:             b,
			MaxCommandDepth: 3,
		}
		err := subject.Exec([]string{"test"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate interface root command runs - flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "root command ran"
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
			},
			Out:             b,
			Err:             b,
			MaxCommandDepth: 3,
		}
		err := subject.Exec([]string{"test", "-flag", "value"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate command runs - no flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "command ran"
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
				SubCommands: []*Command{
					{
						Name: "command",
						Out:  b,
						Err:  b,
						Run: func(command *Command) error {
							command.Out.Write([]byte("command ran"))
							return nil
						},
					},
				},
			},
			Out:             b,
			Err:             b,
			MaxCommandDepth: 3,
		}
		err := subject.Exec([]string{"test", "command"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate command runs - flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "command ran"
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
				SubCommands: []*Command{
					{
						Name: "command",
						Out:  b,
						Err:  b,
						Run: func(command *Command) error {
							command.Out.Write([]byte("command ran"))
							return nil
						},
					},
				},
			},
			Out:             b,
			Err:             b,
			MaxCommandDepth: 3,
		}
		err := subject.Exec([]string{"test", "command", "-flag", "value"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate subcommand runs - no flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "subcommand ran"
		wantPath := "test command subcommand"
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
				SubCommands: []*Command{
					{
						Name: "command",
						Out:  b,
						Err:  b,
						Run: func(command *Command) error {
							command.Out.Write([]byte("command ran"))
							return nil
						},
						SubCommands: []*Command{
							{
								Name:  "subcommand",
								Out:   b,
								Err:   b,
								Flags: flag.NewFlagSet("subcommand", flag.ExitOnError),
								Run: func(command *Command) error {
									command.Out.Write([]byte("subcommand ran"))
									if len(command.Flags.Args()) != 2 {
										t.Errorf("want 2, got %d", len(command.Flags.Args()))
									} else {
										if command.Flags.Arg(0) != "immaarg" {
											t.Errorf("want immaarg, got %s", command.Flags.Arg(0))
										}
										if command.Flags.Arg(1) != "metoo" {
											t.Errorf("want metoo, got %s", command.Flags.Arg(1))
										}
									}
									if command.root {
										t.Errorf("want: false, got %t", command.root)
									}
									if command.path != wantPath {
										t.Errorf("want: %s, got %s", wantPath, command.path)
									}
									return nil
								},
							},
						},
					},
				},
			},
			Out:             b,
			Err:             b,
			MaxCommandDepth: 3,
		}
		err := subject.Exec([]string{"test", "command", "subcommand", "immaarg", "metoo"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate subcommand runs - flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		wantValue := ""
		testCommand := &Command{
			Name:  "subcommand",
			Out:   b,
			Err:   b,
			Flags: flag.NewFlagSet("subcommand", flag.ExitOnError),
			Run: func(command *Command) error {
				command.Out.Write([]byte("subcommand ran"))
				if !command.FlagWasProvided("flag") {
					t.Error("want true, got false")
				}
				if wantValue != "value" {
					t.Errorf("want value, got %s", wantValue)
				}
				return nil
			},
		}
		testCommand.Flags.StringVar(&wantValue, "flag", "", "test flag")
		want := "subcommand ran"
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
				SubCommands: []*Command{
					{
						Name: "command",
						Out:  b,
						Err:  b,
						Run: func(command *Command) error {
							command.Out.Write([]byte("command ran"))
							return nil
						},
						SubCommands: []*Command{
							testCommand,
						},
					},
				},
			},
			Out:             b,
			Err:             b,
			MaxCommandDepth: 3,
		}
		err := subject.Exec([]string{"test", "command", "subcommand", "-flag", "value"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})

	t.Run("validate subcommand runs - flags - full path to binary", func(t *testing.T) {
		b := bytes.NewBufferString("")
		wantValue := ""
		testCommand := &Command{
			Name:  "subcommand",
			Out:   b,
			Err:   b,
			Flags: flag.NewFlagSet("subcommand", flag.ExitOnError),
			Run: func(command *Command) error {
				command.Out.Write([]byte("subcommand ran"))
				if !command.FlagWasProvided("flag") {
					t.Error("want true, got false")
				}
				if wantValue != "value" {
					t.Errorf("want value, got %s", wantValue)
				}
				return nil
			},
		}
		testCommand.Flags.StringVar(&wantValue, "flag", "", "test flag")
		want := "subcommand ran"
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
				SubCommands: []*Command{
					{
						Name: "command",
						Out:  b,
						Err:  b,
						Run: func(command *Command) error {
							command.Out.Write([]byte("command ran"))
							return nil
						},
						SubCommands: []*Command{
							testCommand,
						},
					},
				},
			},
			Out:             b,
			Err:             b,
			MaxCommandDepth: 3,
		}
		err := subject.Exec([]string{"/home/usr/bin/test", "command", "subcommand", "-flag", "value"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
	})
}

func TestCommandInterface_Exec_error(t *testing.T) {
	t.Run("validate interface returns error if no args provided", func(t *testing.T) {
		want := ErrNoArgs
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
			},
			MaxCommandDepth: 3,
		}
		got := subject.Exec(nil)
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if noop state", func(t *testing.T) {
		want := ErrNoOp
		subject := &CommandInterface{
			Name:            "test",
			MaxCommandDepth: 3,
		}
		got := subject.Exec([]string{"test"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if different interface called", func(t *testing.T) {
		want := ErrInvalidInterfaceName
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
			},
			MaxCommandDepth: 3,
		}
		got := subject.Exec([]string{"notvalid"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if invalid command called (noop) - no flag", func(t *testing.T) {
		want := ErrNoOp
		subject := &CommandInterface{
			Name:            "test",
			MaxCommandDepth: 3,
		}
		got := subject.Exec([]string{"test", "nope"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if invalid command called and root not runnable - no flag", func(t *testing.T) {
		want := ErrCommandNotRunnable
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						Run: func(command *Command) error {
							t.Error("unexpected command run")
							return nil
						},
					},
				},
			},
			MaxCommandDepth: 3,
		}
		got := subject.Exec([]string{"test", "nope"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if invalid command called and found command not runnable - no flag", func(t *testing.T) {
		want := ErrCommandNotRunnable
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
					},
				},
			},
			MaxCommandDepth: 3,
		}
		got := subject.Exec([]string{"test", "command", "nope"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if invalid command called", func(t *testing.T) {
		want := ErrCommandNotFound
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						Run: func(command *Command) error {
							t.Error("unexpected command run")
							return nil
						},
					},
				},
			},
			MaxCommandDepth: 3,
		}
		got := subject.Exec([]string{"test", "nope", "-flag", "value"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if invalid subcommand called", func(t *testing.T) {
		want := ErrCommandNotFound
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						Run: func(command *Command) error {
							t.Error("unexpected command run")
							return nil
						},
						SubCommands: []*Command{
							{
								Name: "subcommand",
								Run: func(command *Command) error {
									t.Error("unexpected subcommand run")
									return nil
								},
							},
						},
					},
				},
			},
			MaxCommandDepth: 3,
		}
		got := subject.Exec([]string{"test", "command", "nope", "-flag", "value"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if valid subcommand called on an invalid command", func(t *testing.T) {
		want := ErrCommandNotFound
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						Run: func(command *Command) error {
							t.Error("unexpected command run")
							return nil
						},
						SubCommands: []*Command{
							{
								Name: "subcommand",
								Run: func(command *Command) error {
									t.Error("unexpected subcommand run")
									return nil
								},
							},
						},
					},
				},
			},
			MaxCommandDepth: 3,
		}
		got := subject.Exec([]string{"test", "nope", "subcommand", "-flag", "value"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if command depth invalid - flags", func(t *testing.T) {
		want := ErrCommandDepthInvalid
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						Run: func(command *Command) error {
							t.Error("unexpected command run")
							return nil
						},
						SubCommands: []*Command{
							{
								Name: "subcommand",
								Run: func(command *Command) error {
									t.Error("unexpected subcommand run")
									return nil
								},
							},
						},
					},
				},
			},
			MaxCommandDepth: 3,
		}
		got := subject.Exec([]string{"test", "command", "subcommand", "subsubcommand", "-flag", "value"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("validate interface returns error if command depth invalid - no flags", func(t *testing.T) {
		want := ErrCommandDepthInvalid
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						Run: func(command *Command) error {
							t.Error("unexpected command run")
							return nil
						},
						SubCommands: []*Command{
							{
								Name: "subcommand",
								Run: func(command *Command) error {
									t.Error("unexpected subcommand run")
									return nil
								},
								SubCommands: []*Command{
									{
										Name: "subsubcommand",
									},
								},
							},
						},
					},
				},
			},
			MaxCommandDepth: 3,
		}
		got := subject.Exec([]string{"test", "command", "subcommand", "subsubcommand", "arg1", "arg2"})
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})
}

func TestCommandInterface_Execute(t *testing.T) {
	t.Run("validate command runs - no flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "command ran"
		wantCommand := &Command{
			Name: "command",
			Out:  b,
			Err:  b,
			Run: func(command *Command) error {
				command.Out.Write([]byte("command ran"))
				return nil
			},
		}
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
				SubCommands: []*Command{
					wantCommand,
				},
			},
			Out:             b,
			Err:             b,
			MaxCommandDepth: 3,
		}
		gotCommand, err := subject.Execute([]string{"test", "command"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
		if !reflect.DeepEqual(gotCommand, wantCommand) {
			t.Errorf("want %v, got %v", wantCommand, gotCommand)
		}
	})

	t.Run("validate command runs - flags", func(t *testing.T) {
		b := bytes.NewBufferString("")
		want := "command ran"
		wantCommand := &Command{
			Name: "command",
			Out:  b,
			Err:  b,
			Run: func(command *Command) error {
				command.Out.Write([]byte("command ran"))
				return nil
			},
		}
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				Out:  b,
				Err:  b,
				Run: func(command *Command) error {
					command.Out.Write([]byte("root command ran"))
					return nil
				},
				SubCommands: []*Command{
					wantCommand,
				},
			},
			Out:             b,
			Err:             b,
			MaxCommandDepth: 3,
		}
		gotCommand, err := subject.Execute([]string{"test", "command", "-flag", "value"})
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		got, err := ioutil.ReadAll(b)
		if err != nil {
			t.Errorf("[err] want nil, got %s", err)
		}
		if string(got) != want {
			t.Errorf("want: %s, got %s", want, string(got))
		}
		if !reflect.DeepEqual(gotCommand, wantCommand) {
			t.Errorf("want %v, got %v", wantCommand, gotCommand)
		}
	})
}

func TestCommandInterface_searchPathForCommand(t *testing.T) {
	t.Run("validate command is found", func(t *testing.T) {
		want := &Command{Name: "command2"}
		wantPos := 0
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
					},
					{
						Name: "command2",
					},
					{
						Name: "command3",
					},
				},
			},
		}
		got, found, position := subject.searchPathForCommand([]string{"command2"}, false)
		if !found {
			t.Errorf("want true, got %t", found)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
		if position != wantPos {
			t.Errorf("want %d, got %d", wantPos, position)
		}
	})

	t.Run("validate command is found - deep", func(t *testing.T) {
		want := &Command{Name: "dang", Description: "It's me you're looking for."}
		wantPos := 3
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
						SubCommands: []*Command{
							{
								Name: "subcommand",
							},
						},
					},
					{
						Name: "command2",
						SubCommands: []*Command{
							{
								Name: "subcommand2",
								SubCommands: []*Command{
									{
										Name: "subsubcommand3",
										SubCommands: []*Command{
											{
												Name:        "dang",
												Description: "It's me you're looking for.",
											},
										},
									},
								},
							},
						},
					},
					{
						Name: "command3",
					},
				},
			},
		}
		got, found, position := subject.searchPathForCommand([]string{"command2", "subcommand2", "subsubcommand3", "dang"}, false)
		if !found {
			t.Errorf("want true, got %t", found)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
		if position != wantPos {
			t.Errorf("want %d, got %d", wantPos, position)
		}
	})

	t.Run("validate command is found with alias", func(t *testing.T) {
		want := &Command{Name: "command2", Aliases: []string{"heyo"}}
		wantPos := 0
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name:    "command",
						Aliases: []string{"nope"},
					},
					{
						Name:    "command2",
						Aliases: []string{"heyo"},
					},
					{
						Name: "command3",
					},
				},
			},
		}
		got, found, position := subject.searchPathForCommand([]string{"heyo"}, false)
		if !found {
			t.Errorf("want true, got %t", found)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
		if position != wantPos {
			t.Errorf("want %d, got %d", wantPos, position)
		}
	})

	t.Run("validate command is found - partial", func(t *testing.T) {
		testRun := func(command *Command) error {
			if len(command.Flags.Args()) != 1 {
				t.Errorf("want 1, got %d", len(command.Flags.Args()))
			} else {
				if command.Flags.Arg(0) != "nope" {
					t.Errorf("want nope, got %s", command.Flags.Arg(0))
				}
			}
			return nil
		}
		want := &Command{Name: "subcommand2", Run: testRun, SubCommands: []*Command{{Name: "unused"}}}
		wantPos := 1
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
					},
					{
						Name: "command2",
						SubCommands: []*Command{
							want,
						},
					},
					{
						Name: "command3",
					},
				},
			},
		}
		got, found, position := subject.searchPathForCommand([]string{"command2", "subcommand2", "nope"}, true)
		if !found {
			t.Errorf("want true, got %t", found)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
		if position != wantPos {
			t.Errorf("want %d, got %d", wantPos, position)
		}
	})

	t.Run("validate command is not found", func(t *testing.T) {
		wantPos := -1
		subject := &CommandInterface{
			Name: "test",
			RootCommand: &Command{
				Name: "test",
				SubCommands: []*Command{
					{
						Name: "command",
					},
					{
						Name: "command2",
					},
					{
						Name: "command3",
					},
				},
			},
		}
		got, found, position := subject.searchPathForCommand([]string{"command", "nope"}, false)
		if found {
			t.Errorf("want false, got %t", found)
		}
		if got != nil {
			t.Errorf("want nil, got %v", got)
		}
		if position != wantPos {
			t.Errorf("want %d, got %d", wantPos, position)
		}
	})
}

func TestCommandInterface_hasFlags(t *testing.T) {
	t.Run("validate flag is found in correct position", func(t *testing.T) {
		want := 3
		args := []string{"interface", "command", "subcommand", "-flag", "value", "heyo"}
		subject := &CommandInterface{}
		got, hasFlag := subject.hasFlags(args)
		if !hasFlag {
			t.Errorf("want true, got %t", hasFlag)
		}
		if got != want {
			t.Errorf("want %d, got %d", want, got)
		}
	})

	t.Run("validate flag is found in correct position (--)", func(t *testing.T) {
		want := 2
		args := []string{"interface", "command", "--flag", "value", "heyo"}
		subject := &CommandInterface{}
		got, hasFlag := subject.hasFlags(args)
		if !hasFlag {
			t.Errorf("want true, got %t", hasFlag)
		}
		if got != want {
			t.Errorf("want %d, got %d", want, got)
		}
	})

	t.Run("validate flag is not found", func(t *testing.T) {
		want := -1
		args := []string{"interface", "command", "subcommand", "notaflag", "value"}
		subject := &CommandInterface{}
		got, hasFlag := subject.hasFlags(args)
		if hasFlag {
			t.Errorf("want false, got %t", hasFlag)
		}
		if got != want {
			t.Errorf("want %d, got %d", want, got)
		}
	})
}
