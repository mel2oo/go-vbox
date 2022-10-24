package vbox

import (
	"bufio"
	"strconv"
	"strings"
)

// MachineState stores the last retrieved VM state.
type MachineState string

const (
	// Running is a MachineState value.
	Running = MachineState("running")
	// Poweroff is a MachineState value.
	Poweroff = MachineState("poweroff")
	// Paused is a MachineState value.
	Paused = MachineState("paused")
	// Saved is a MachineState value.
	Saved = MachineState("saved")
	// Aborted is a MachineState value.
	Aborted = MachineState("aborted")
	// Stopping is a MachineState value.
	Stopping = MachineState("stopping")
	// Guru meditation is a MachineState value.
	Gurumeditation = MachineState("gurumeditation")
)

type Machine struct {
	*Manage `json:"-"`

	Name     string       `json:"name"`
	UUID     string       `json:"uuid"`
	OSType   string       `json:"ostype"`
	Firmware string       `json:"firmware"`
	CfgFile  string       `json:"cfgfile"`
	CPUs     uint         `json:"cpus"`
	Memory   uint         `json:"memory"`
	VRAM     uint         `json:"vram"`
	State    MachineState `json:"state"`
}

// cmd: vboxmanage list vms
func (m *Manage) ListMachines() ([]*Machine, error) {
	stdout, err := m.cmd.RunOutput(m.bin, "list", "vms")
	if err != nil {
		return nil, err
	}

	vmss := make([]*Machine, 0)
	scan := bufio.NewScanner(strings.NewReader(stdout))
	for scan.Scan() {
		res := reVMNameUUID.FindStringSubmatch(scan.Text())
		if res == nil {
			continue
		}

		vm, err := m.GetMachine(res[1])
		if err != nil {
			continue
		}

		vmss = append(vmss, vm)
	}

	return vmss, nil
}

// cmd: vboxmanage showvminfo < uuid | vmname > [--details] [--machinereadable] [--password-id] [--password]
func (m *Manage) GetMachine(id string) (*Machine, error) {
	stdout, err := m.cmd.RunOutput(m.bin, "showvminfo", id, "--machinereadable")
	if err != nil {
		return nil, err
	}

	prop := make(map[string]string, 0)
	scan := bufio.NewScanner(strings.NewReader(stdout))
	for scan.Scan() {
		res := reVMInfoLine.FindStringSubmatch(scan.Text())
		if res == nil {
			continue
		}

		var key, val string
		for i := 1; i < len(res); i++ {
			if len(res[i]) > 0 {
				if len(key) == 0 {
					key = res[i]
					continue
				}

				if len(val) == 0 {
					val = res[i]
					break
				}
			}
		}

		prop[key] = val
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}

	vm := Machine{
		Manage: m,

		Name:     prop["name"],
		UUID:     prop["UUID"],
		OSType:   prop["ostype"],
		Firmware: prop["firmware"],
		CfgFile:  prop["CfgFile"],
		State:    MachineState(prop["VMState"]),
	}

	cpus, err := strconv.Atoi(prop["cpus"])
	if err == nil {
		vm.CPUs = uint(cpus)
	}

	memory, err := strconv.Atoi(prop["memory"])
	if err == nil {
		vm.Memory = uint(memory)
	}

	vram, err := strconv.Atoi(prop["vram"])
	if err == nil {
		vm.VRAM = uint(vram)
	}

	return &vm, nil
}

// VBoxManage startvm < uuid | vmname ...> [--putenv=name[=value]] [--type= [ gui | headless | sdl | separate ]] --password file --password-id password identifier
func (m *Machine) Start() error {
	switch m.State {
	case Running:
		return nil
	case Paused:
		return m.Resume()
	case Poweroff, Saved, Aborted:
		return m.cmd.Run(m.bin, "startvm", m.UUID, "--type", "headless")
	case Stopping, Gurumeditation:
		if err := m.cmd.Run(m.bin, "startvm", m.UUID, "--type", "emergencystop"); err != nil {
			return err
		}
		return m.cmd.Run(m.bin, "startvm", m.UUID, "--type", "headless")
	default:
		return ErrMachineState
	}
}

// VBoxManage controlvm < uuid | vmname > pause
func (m *Machine) Pause() error {
	switch m.State {
	case Poweroff, Paused, Saved, Aborted:
		return nil
	case Stopping, Gurumeditation:
		return ErrMachineState
	default:
		return m.cmd.Run(m.bin, "controlvm", m.UUID, "pause")
	}
}

// VBoxManage controlvm < uuid | vmname > resume
func (m *Machine) Resume() error {
	switch m.State {
	case Paused:
		return m.cmd.Run(m.bin, "controlvm", m.UUID, "resume")
	default:
		return ErrMachineState
	}
}

// VBoxManage controlvm < uuid | vmname > reset
func (m *Machine) Reset() error {
	return m.cmd.Run("controlvm", m.UUID, "reset")
}

// VBoxManage controlvm < uuid | vmname > poweroff
func (m *Machine) Poweroff() error {
	switch m.State {
	case Poweroff, Aborted, Saved:
		return nil
	case Stopping, Gurumeditation:
		return m.cmd.Run(m.bin, "startvm", m.UUID, "--type", "emergencystop")
	default:
		return m.cmd.Run(m.bin, "controlvm", m.UUID, "poweroff")
	}
}

// VBoxManage controlvm < uuid | vmname > savestate
func (m *Machine) Save() error {
	return m.cmd.Run(m.bin, "controlvm", m.UUID, "savestate")
}

// VBoxManage controlvm < uuid | vmname > acpipowerbutton
// VBoxManage controlvm < uuid | vmname > acpisleepbutton
// VBoxManage controlvm < uuid | vmname > reboot
// VBoxManage controlvm < uuid | vmname > shutdown [--force]

// VBoxManage snapshot <uuid|vmname>
// VBoxManage snapshot <uuid|vmname> take <snapshot-name> [--description=description] [--live] [--uniquename Number,Timestamp,Space,Force]
// VBoxManage snapshot <uuid|vmname> delete <snapshot-name>
// VBoxManage snapshot <uuid|vmname> restore <snapshot-name>
// VBoxManage snapshot <uuid|vmname> restorecurrent
// VBoxManage snapshot <uuid|vmname> edit < snapshot-name | --current > [--description=description] [--name=new-name]
// VBoxManage snapshot <uuid|vmname> list [[--details] | [--machinereadable]]
// VBoxManage snapshot <uuid|vmname> showvminfo <snapshot-name>
