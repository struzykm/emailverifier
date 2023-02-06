package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"runtime"
	"strings"
)

var pl = fmt.Println
var pf = fmt.Printf
var pr = fmt.Print
var lg = log.Printf

// Terminal colors
const (
	colorRed        = "\033[31m"
	colorDarkBlue   = "\033[34m"
	colorReset      = "\033[0m"
	colorDarkYellow = "\033[33m"
	colorGreen      = "\033[92m"
	colorBlue       = "\033[94m"
	colorMagenta    = "\033[95m"
	colorCyan       = "\033[96m"
)

func main() {

	if runtime.GOOS == "windows" {
		clrScr()
	}

	scan := bufio.NewScanner(os.Stdin)
	pf("%sEnter a invalid e-mail or any tekst/char to end program\n", colorGreen)
	pf("%sEnter e-mail :%s", colorRed, colorDarkYellow)
	for scan.Scan() {
		if email, err := isEmail(scan.Text()); err == nil {
			checkDomain(email)
			//break
		} else {
			lg("%sError: %v%s\n", colorDarkBlue, colorReset, err)
			break
		}
	}
	if err := scan.Err(); err != nil {
		log.Fatalf("Error : %v", err)
	}
}

func checkDomain(email string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	components := strings.Split(email, "@")
	domain := components[1]

	if mxRecords, err := net.LookupMX(domain); err != nil {
		lg("%sError: %v%s\n", colorDarkBlue, colorReset, err)
	} else {
		if len(mxRecords) > 0 {
			hasMX = true
		}
	}

	if txtRecords, err := net.LookupTXT(domain); err != nil {
		lg("%sError: %v%s\n", colorDarkBlue, colorReset, err)
	} else {
		for _, record := range txtRecords {
			if strings.HasPrefix(record, "v=spf1") {
				hasSPF = true
				spfRecord = record
				break
			}
		}
	}

	if dmarcRecords, err := net.LookupTXT("_dmarc." + domain); err != nil {
		lg("%sError: %v%s\n", colorDarkBlue, colorReset, err)
	} else {
		for _, record := range dmarcRecords {
			if strings.HasPrefix(record, "v=DMARC1") {
				hasDMARC = true
				dmarcRecord = record
				break
			}
		}
	}

	pf("%v", colorReset)
	pf("Domain : %s\n", domain)
	pf("hasMX : %v\n", hasMX)
	pf("hasSPF : %v\n", hasSPF)
	pf("spfRecord : %s\n", spfRecord)
	pf("hasDMARC : %v\n", hasDMARC)
	pf("dmarcRecord : %s\n\n", dmarcRecord)
	pf("%sNext e-mail :%s", colorRed, colorDarkYellow)
}

func isEmail(s string) (string, error) {
	r, _ := regexp.Compile(`[\w._%+-]{1,20}@[\w.-]{2,20}.[A-Za-z]{2,3}`)
	if r.MatchString(s) {
		return s, nil
	} else {
		return "", fmt.Errorf("Not a valid e-mail")
	}
}

func clrScr() {
	fmt.Print("\033[H\033[2J")
	pf("%s", colorReset)
}
