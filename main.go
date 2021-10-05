package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("require 1 or more args")
	}
	err := remoteEdit(os.Args[1])
	if err != nil {
		log.Fatalf("failed to remote edit %v", err)
	}
}

var rx = regexp.MustCompile(`([^:]*)(?::(\d+))?`)

func remoteEdit(s string) error {
	m := rx.FindStringSubmatch(s)
	if m == nil {
		return errors.New("syntax problem on argument")
	}

	path := m[1]
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("not regular file: %s", path)
	}

	var num int
	if len(m) >= 3 && m[2] != "" {
		n, err := strconv.ParseInt(m[2], 10, 64)
		if err != nil {
			return fmt.Errorf("syntax problem on line num: %w", err)
		}
		num = int(n)
	}

	args := make([]string, 0, 3)
	args = append(args, "--remote")
	if num > 0 {
		args = append(args, "+"+strconv.Itoa(num))
	}
	args = append(args, path)
	return exec.Command("gvim", args...).Run()
}
