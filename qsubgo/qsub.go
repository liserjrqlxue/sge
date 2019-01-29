package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	numProc = flag.Int(
		"p",
		0,
		"-l p=n, will be overwrited by -l",
	)
	virtualFree = flag.String(
		"vf",
		"",
		"-l vf=m, will be overwrited by -l",
	)
	project = flag.String(
		"P",
		"",
		"-P project",
	)
	cwd = flag.Bool(
		"cwd",
		false,
		"-cwd",
	)
	queue = flag.String(
		"q",
		"",
		"-q queue.q",
	)
	hardResourceList = flag.String(
		"l",
		"",
		"-l p=n,vf=m",
	)
	bingding = flag.String(
		"bingding",
		"",
		"-bingding linear:n, n will be overwited by -p",
	)
)

func main() {
	flag.Parse()
	fmt.Println(flag.Args())

	var args []string
	if *cwd {
		args = append(args, "-cwd")
	}
	if *hardResourceList != "" {
		hrList := commaSplit(*hardResourceList)
		hrHash := str2map(hrList, "=")
		if hrHash["vf"] != "" {
			*virtualFree = hrHash["vf"]
		}
		if hrHash["p"] != "" {
			n, err := strconv.Atoi(hrHash["p"])
			if err != nil {
				fmt.Println("-l ", *hardResourceList, " has not give proper num_proc value")
				flag.Usage()
				os.Exit(1)
			} else {
				*numProc = n
			}
		} else if *numProc <= 0 {
			fmt.Println("-l ", *hardResourceList, " or -p has not give proper num_proc value")
			flag.Usage()
			os.Exit(1)
		}
	} else if *virtualFree != "" && *numProc > 0 {
		*hardResourceList = "vf=" + *virtualFree + ",p=" + strconv.Itoa(*numProc)
	} else {
		fmt.Println("no -vf -p or -l to set num_proc and virtual_free")
		flag.Usage()
		os.Exit(1)
	}
	args = append(args, "-l", *hardResourceList)

	if *bingding != "" {
		bd := commaSplit(*bingding)
		bdMap := str2map(bd, ":")
		bdMap["linear"] = strconv.Itoa(*numProc)
		var newBd []string
		for k, v := range bdMap {
			newBd = append(newBd, k+":"+v)
		}
		*bingding = strings.Join(newBd, ",")
	} else {
		*bingding = "linear:" + strconv.Itoa(*numProc)
	}
	args = append(args, "-binding", *bingding)

	if *project != "" {
		args = append(args, "-P", *project)
	} else {
		fmt.Println("no -P project to set project")
		flag.Usage()
		os.Exit(1)
	}

	if *queue != "" {
		args = append(args, "-q", *queue)
	}

	args = append(args, flag.Args()...)

	fmt.Printf("run cmd:\n qsub %s\n", args)
	runCmd("qsub", args...)
}

func commaSplit(str string) []string {
	return strings.Split(str, ",")
}

func str2map(strs []string, sep string) map[string]string {
	var hash = make(map[string]string)
	for _, str := range strs {
		kv := strings.SplitN(str, sep, 2)
		hash[kv[0]] = kv[1]
	}
	return hash
}

func runCmd(name string, args ...string) {
	c := exec.Command(name, args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err := c.Run()
	fmt.Println(err)
}
