package main

import (
	"fmt"
	"path/filepath"
	"os"
	"io/ioutil"
	"github.com/fatih/color"
	"strconv"
	"strings"
	"sync"
)

var lines int = 0
var sg sync.WaitGroup

//Plan on shifting to the flag stdlib from os.Args
func handleArgs(args []string) []string{

	if args[1] == "-h" || args[1] == "--help"{
		color.Yellow("acr is a program that finds the number of lines of code in a directory.\nOptions-\n1. -h or --help for basic info and options.\n2. -f or --filter for filtering by a filetype.")
		os.Exit(0)
	}  else if args[1] == "-f" || args[1] == "--filter"{
		if strings.Contains(args[2], ","){
			return strings.Split(args[2],",")
		}  else{
			return []string{args[2]}		
		}
		
	}  

	return []string{""}
	
}

func checkSuffix(path string, args []string) bool{
	for _, v := range args{
		if strings.HasSuffix(path,v){
			return true
		}
	}
	return false
}

// Checks lines one file at a time
func checkLines(path string){

	dat, _ := ioutil.ReadFile(path)
	lines = lines + len(strings.Split(string(dat),"\n"))
	
}

//Main logic
func fileWalk(directory string, args ...[]string){

	scannablePaths := []string{}

	if len(args) !=0{
		filepath.Walk(directory, func(path string, fileinfo os.FileInfo, err error)error{
			if err !=nil{
				color.Red(err.Error())
			}
			if !fileinfo.IsDir() && checkSuffix(path, args[0]){
				scannablePaths = append(scannablePaths,path)
			}
			return nil	
		})
	} else{
		filepath.Walk(directory, func(path string, fileinfo os.FileInfo, err error)error{
			if err !=nil{
				panic(err)
			}
			if !fileinfo.IsDir(){
				scannablePaths = append(scannablePaths,path)
			}
			return nil	
		})
	}

	sg.Add(len(scannablePaths))
	for _, v := range scannablePaths{
		go func(v string){
			defer sg.Done()
			checkLines(v)
		}(v)
	}

	sg.Wait()
	
}

func checkArgs(args []string) []string{
	var resp []string
	if len(os.Args) > 1{
		resp = handleArgs(os.Args)
		} else {
			resp = []string{""}
		}
	return resp	
}

func main(){
	color.Blue("scanning...")
	resp := checkArgs(os.Args)
	d ,_ := os.Getwd()
	if resp[0] == ""{
		fileWalk(d)
	}  else{
		fileWalk(d, resp)
	}
	
	amount := color.CyanString(strconv.Itoa(lines))
	firstphrase := color.GreenString("The directory has ")
	lastphrase := color.GreenString(" lines of code/text")
	fmt.Println(firstphrase+amount+lastphrase)
}
