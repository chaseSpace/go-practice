package main

import "regexp"

func main() {
	re := regexp.MustCompile(`0{3,5}|1{3,5}|2{3,5}|3{3,5}|4{3,5}|5{3,5}|6{3,5}|7{3,5}|8{3,5}|9{3,5}|^520|^1314|.520|.1314`)

	for _, s := range []string{`00000123`, `100000654`, `54000000`} {
		println(re.MatchString(s))
	}

}
