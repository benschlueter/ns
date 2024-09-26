package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "ns":
		ns()
	default:
		panic("pass me an argument please")
	}
}

var userID = 1000

func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command("/proc/self/exe", append([]string{"ns"}, os.Args[2:]...)...)
	// cmd := exec.Command("/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:                 syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
		Unshareflags:               syscall.CLONE_NEWNS,
		GidMappingsEnableSetgroups: true,
		Setsid:                     true,
		Setctty:                    true,
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}

func ns() {
	fmt.Printf("Running in namespace %v as %d\n", os.Args[2:], os.Getpid())
	// Requires root privileges
	if err := syscall.Chroot("/sandbox"); err != nil {
		fmt.Println("chroot error: ", err)
	}
	// set the working directory inside container
	if err := syscall.Chdir("/"); err != nil {
		fmt.Println("chdir error: ", err)
	}
	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		fmt.Println("mount error: ", err)
	}
	fmt.Printf("User ID: %d\n", os.Getuid())
	fmt.Printf("Effective UID: %d\n", syscall.Geteuid())
	fmt.Printf("Group ID: %d\n", os.Getgid())
	if err := syscall.Setgid(userID); err != nil {
		fmt.Println("setgid error: ", err)
		return
	}
	if err := syscall.Setgroups([]int{userID}); err != nil {
		fmt.Println("setgroups error: ", err)
		return
	}
	if err := syscall.Setuid(userID); err != nil {
		fmt.Println("setuid error: ", err)
		return
	}
	fmt.Printf("User ID: %d\n", os.Getuid())
	fmt.Printf("Effective UID: %d\n", syscall.Geteuid())
	fmt.Printf("Group ID: %d\n", os.Getgid())
	if err := syscall.Exec(os.Args[2], os.Args[2:], os.Environ()); err != nil {
		fmt.Println("exec error: ", err)
	}

	// syscall.Unmount("/proc", 0)
}
