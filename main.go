package main

import (
	"database/sql"
	"flag"
	"fmt"
	"reflect"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

// Options Struct
type Options struct {
	mode        string
	sourceSrv   string
	targetSrv   string
	outputWidth int
	sourceUser  string
	sourcePass  string
	targetUser  string
	targetPass  string
	user        string
	pass        string
}

// DiffValues struct
type DiffValues struct {
	sValue string
	tValue string
}

// DiffUsers struct
type DiffUsers struct {
	vSusers []string
	vTusers []string
}

// DiffValuesType type
type DiffValuesType map[string]*DiffValues

// DiffUsersType type
type DiffUsersType map[string]*DiffUsers

// VariabMaxLen length
var VariabMaxLen int

// NOTFOUND variables
const NOTFOUND = "NOT_FOUND_THE_VARIABLE"
const VARIABLES = "Variables"
const USER = "User"
const SSERVER = "source"
const TSERVER = "target"

// OutputWidth variables
var OutputWidth int

// Source string
var Source string

// Target string
var Target string

// Mode string
var Mode string

// Global Mutex
var (
	m = sync.Mutex{}
)

func main() {
	modelist := map[string]int{VARIABLES: 1, USER: 2}
	//typeServers := []string{"source", "target"}
	typeServers := map[string]*sql.DB{"source": nil, "target": nil}
	Opt := Options{}
	flag.StringVar(&Opt.mode, "m", "", fmt.Sprintf("Mode [%s,%s]", VARIABLES, USER))
	flag.StringVar(&Opt.sourceSrv, "s", "", "SourceServer[ip:port]")
	flag.StringVar(&Opt.targetSrv, "t", "", "TargetServer[ip:port]")
	flag.StringVar(&Opt.sourceUser, "su", "", "Source User")
	flag.StringVar(&Opt.sourcePass, "sp", "", "Source Password")
	flag.StringVar(&Opt.targetUser, "tu", "", "Target User")
	flag.StringVar(&Opt.targetPass, "tp", "", "Target Password")
	flag.StringVar(&Opt.user, "u", "", "Source and Target User")
	flag.StringVar(&Opt.pass, "p", "", "Source and Target Password")
	flag.IntVar(&Opt.outputWidth, "o", 50, "outputWidth")
	flag.Parse()

	Source = Opt.sourceSrv
	Target = Opt.targetSrv
	sUser := ""
	tUser := ""
	sPass := ""
	tPass := ""
	OutputWidth = Opt.outputWidth
	Mode = Opt.mode

	if _, ok := modelist[Mode]; !ok {
		fmt.Println(fmt.Sprintf("Choose -m mode [%s,%s]", VARIABLES, USER))
		return
	}

	if Source == "" || Target == "" {
		fmt.Println(fmt.Sprintf("Set source server[-s] and target server[-t]"))
		return
	}

	if Opt.user != "" {
		sUser = Opt.user
		tUser = Opt.user
	}
	if Opt.pass != "" {
		sPass = Opt.pass
		tPass = Opt.pass
	}
	if Opt.sourceUser != "" {
		sUser = Opt.sourceUser
	}
	if Opt.targetUser != "" {
		tUser = Opt.targetUser
	}
	if Opt.sourcePass != "" {
		sPass = Opt.sourcePass
	}
	if Opt.targetPass != "" {
		tPass = Opt.targetPass
	}
	if sUser == "" || tUser == "" || sPass == "" || tPass == "" {
		fmt.Println(fmt.Sprintf("Set User[-u] and Password [-p]"))
		return
	}

	typeServers[SSERVER] = dbconn(Source, sUser, sPass)
	defer typeServers[SSERVER].Close()
	typeServers[TSERVER] = dbconn(Target, tUser, tPass)
	defer typeServers[TSERVER].Close()
	wg := &sync.WaitGroup{}
	VariabMaxLen = len(Mode)

	if Mode == VARIABLES {
		varList := make(DiffValuesType, 0)
		for key, db := range typeServers {
			wg.Add(1)
			go func(key string, db *sql.DB) {
				selectVariables(db, varList, key)
				wg.Done()
			}(key, db)
		}
		wg.Wait()
		diffValSourceTarget(varList)
	}
	if Mode == USER {
		varList := make(DiffUsersType, 0)
		for key, db := range typeServers {
			wg.Add(1)
			go func(key string, db *sql.DB) {
				selectUser(db, varList, key)
				wg.Done()
			}(key, db)
		}
		wg.Wait()
		diffUserSourceTarget(varList)
	}
}

func diffUserSourceTarget(vType DiffUsersType) {
	maxLen(vType)
	formater := createFormat()
	initFormat(Mode, formater)
	for key, val := range vType {
		if !chkStrSliceEqual(val.vSusers, val.vTusers) {
			str := returnMoreThanElements(val.vSusers, val.vTusers)
			for i, _ := range str {
				s, t := returnStringIfHave(val.vSusers, val.vTusers, i)
				ss := splitValByNumber(s, OutputWidth)
				tt := splitValByNumber(t, OutputWidth)
				sstr := returnMoreThanElements(ss, tt)
				for ii, _ := range sstr {
					vn := padS(VariabMaxLen, " ")
					if i == 0 && ii == 0 {
						vn = key
					}
					sss, ttt := returnStringIfHave(ss, tt, ii)
					fmt.Println(fmt.Sprintf(formater, vn, sss, ttt))
				}
			}
			fmt.Println(fmt.Sprintf(formater, padS(VariabMaxLen, "-"), padS(OutputWidth, "-"), padS(OutputWidth, "-")))
		}
	}
}

func chkStrSliceEqual(s []string, t []string) bool {
	if len(s) != len(t) {
		return false
	}
	for _, sv := range s {
		rtn := 0
		for _, tv := range t {
			if sv == tv {
				rtn = 1
				break
			}
		}
		if rtn == 0 {
			return false
		}
	}
	return true
}
func diffValSourceTarget(vType DiffValuesType) {
	maxLen(vType)
	formater := createFormat()
	initFormat(Mode, formater)
	for key, vVal := range vType {
		if vVal.sValue != vVal.tValue {
			vSstrs := splitValByNumber(vVal.sValue, OutputWidth)
			vTstrs := splitValByNumber(vVal.tValue, OutputWidth)
			if len(vSstrs) == 1 && len(vTstrs) == 1 {
				fmt.Println(fmt.Sprintf(formater, key, vVal.sValue, vVal.tValue))
				fmt.Println(fmt.Sprintf(formater, padS(VariabMaxLen, "-"), padS(OutputWidth, "-"), padS(OutputWidth, "-")))
				continue
			}
			str := returnMoreThanElements(vSstrs, vTstrs)
			for idx, _ := range str {
				vn := padS(VariabMaxLen, " ")
				if idx == 0 {
					vn = key
				}
				s, t := returnStringIfHave(vSstrs, vTstrs, idx)
				fmt.Println(fmt.Sprintf(formater, vn, s, t))
			}
			fmt.Println(fmt.Sprintf(formater, padS(VariabMaxLen, "-"), padS(OutputWidth, "-"), padS(OutputWidth, "-")))
		}
	}
}

func createFormat() string {
	formater := "%-" + fmt.Sprintf("%d", VariabMaxLen) + "s"
	formater += " %-" + fmt.Sprintf("%d", OutputWidth) + "s %-" + fmt.Sprintf("%d", OutputWidth) + "s"
	return formater
}

func initFormat(s string, formater string) {
	fmt.Println(fmt.Sprintf(formater, s, limitString(Source, OutputWidth), limitString(Target, OutputWidth)))
	fmt.Println(fmt.Sprintf(formater, padS(VariabMaxLen, "-"), padS(OutputWidth, "-"), padS(OutputWidth, "-")))
}

func returnStringIfHave(s []string, t []string, idx int) (string, string) {
	var sBuf, tBuf string
	if len(s) > idx {
		sBuf = s[idx]
	} else {
		sBuf = ""
	}
	if len(t) > idx {
		tBuf = t[idx]
	} else {
		tBuf = ""
	}
	return sBuf, tBuf
}

func returnMoreThanElements(a []string, b []string) []string {
	buf := a
	if len(b) > len(a) {
		buf = b
	}
	return buf
}

func maxLen(v interface{}) {
	switch v.(type) {
	case DiffValuesType:
		for key, vVal := range v.(DiffValuesType) {
			if vVal.sValue != vVal.tValue {
				if VariabMaxLen < len(key) {
					VariabMaxLen = len(key)
				}
			}
		}
	case DiffUsersType:
		for key, val := range v.(DiffUsersType) {
			if !reflect.DeepEqual(val.vSusers, val.vTusers) {
				if VariabMaxLen < len(key) {
					VariabMaxLen = len(key)
				}
			}
		}
	}
}

func splitValByNumber(str string, splitno int) []string {
	strSlice := make([]string, 0)
	loopLen := len(str) / splitno
	// 一行で済むためreturnする
	if loopLen == 0 {
		strSlice = append(strSlice, str)
		return strSlice
	}
	start := 0
	end := splitno
	for i := 0; i <= loopLen; i++ {
		if splitno >= len(str[start:]) {
			end = len(str)
		}
		strSlice = append(strSlice, str[start:end])
		start += splitno
		end += splitno
	}
	return strSlice
}

func selectVariables(db *sql.DB, varList DiffValuesType, serv string) {
	var rows *sql.Rows
	var err error
	rows, err = db.Query("SHOW GLOBAL VARIABLES")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	var vName, val string
	for rows.Next() {
		if err := rows.Scan(&vName, &val); err != nil {
			panic(err.Error())
		}
		m.Lock()
		if v, ok := varList[vName]; ok {
			if serv == SSERVER {
				v.sValue = val
			} else {
				v.tValue = val
			}
		} else {
			if serv == SSERVER {
				varList[vName] = &DiffValues{sValue: val, tValue: NOTFOUND}
			} else {
				varList[vName] = &DiffValues{sValue: NOTFOUND, tValue: val}
			}
		}
		m.Unlock()
	}
}

func selectGrant(db *sql.DB, usrIP string) []string {
	rows, err := db.Query("show grants for " + usrIP)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	var grant string
	grants := make([]string, 0)
	for rows.Next() {
		if err := rows.Scan(&grant); err != nil {
			panic(err.Error())
		}
		passwordIndex := strings.Index(grant, "IDENTIFIED BY PASSWORD")
		grantOptionIndex := strings.Index(grant, "WITH GRANT OPTION")
		if passwordIndex != -1 {
			grant = grant[:passwordIndex-1]
			if grantOptionIndex != -1 {
				grant += " WITH GRANT OPTION"
			}
		}
		grants = append(grants, grant)

	}
	return grants

}
func selectUser(db *sql.DB, varList DiffUsersType, serv string) {
	rows, err := db.Query("select user , host from mysql.user order by 1,2")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	var user, host string
	for rows.Next() {
		if err := rows.Scan(&user, &host); err != nil {
			panic(err.Error())
		}
		usrIP := "'" + user + "'@'" + host + "'"
		grants := selectGrant(db, usrIP)
		m.Lock()
		if v, ok := varList[usrIP]; ok {
			if serv == SSERVER {
				v.vSusers = grants
			} else {
				v.vTusers = grants
			}
		} else {
			if serv == SSERVER {
				varList[usrIP] = &DiffUsers{vSusers: grants}
			} else {
				varList[usrIP] = &DiffUsers{vTusers: grants}
			}
		}
		m.Unlock()
	}
}

func dbconn(host string, dbuser string, dbpasswd string) *sql.DB {
	db, err := sql.Open("mysql", dbuser+":"+dbpasswd+"@tcp("+host+")/mysql")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func padS(no int, char string) string {
	rtn := ""
	for i := 0; i < no; i++ {
		rtn += char
	}
	return rtn
}

func limitString(v string, maxlen int) string {
	if len(v) > maxlen {
		v = v[:maxlen-1] + "."
	}
	return v
}
