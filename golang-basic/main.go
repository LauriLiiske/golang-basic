package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// delete result file contents before the start of every test
	err := os.Truncate(os.Args[2], 0)
	if err != nil {
		log.Fatal(err)
	}

	// read file and save each line to a slice
	fileName := os.Args[1]
	contentSlice := readLine(fileName)

	for _, line := range contentSlice { // loop through the line_slice for line
		// fix all problems with the strings
		fixedPunct := fixPunctuation(line)             // fix the punctuation of the string, return string
		fixedUpperCommas := fixUpperCommas(fixedPunct) // fix upper commas
		fixedAtoAn := convertAtoAn(fixedUpperCommas)   // fix a to an conversion
		fixedContent := makeIntoSlice(fixedAtoAn)      // make a slice out of the string, return string
		finalFixedContent := Connector(fixedContent)   // return slice of fixed strings
		finalModSlice := Modify(finalFixedContent)     // return the final fixed slice of strings with modifications completed
		result := strings.Join(finalModSlice, " ")     // connect all the strings from result slice
		fixedPunct = fixPunctuation(result)            // fix punctuation again. reason "hello (cap,2) : world" -> "hello : world"
		// write to file
		writeFile(fixedPunct)
	}
}

func Modify(s []string) []string { // modify up low cap etc cases
	// declare all matching for the multiples
	multipleUp := regexp.MustCompile(`\(up,(\d)\)`)
	multipleLow := regexp.MustCompile(`\(low,(\d)\)`)
	multipleCap := regexp.MustCompile(`\(cap,(\d)\)`)

	linesSlice := []string{}
	var index int = 0

	for i := 0; i < len(s); i++ {
		word := s[i]
		if word == "(up)" { // uppercase
			linesSlice[index-1] = strings.ToUpper(s[i-1])
		} else if word == "(low)" { // lowercase
			linesSlice[index-1] = strings.ToLower(s[i-1])
		} else if word == "(cap)" { // capitalize
			m := strings.ToLower(s[i-1])
			linesSlice[index-1] = strings.Title(m)
		} else if multipleUp.MatchString(word) { // MULTIPLE TO UP
			number := string(s[i][4])            // take the number from the word to see how many times we have to go back
			n, _ := strconv.Atoi(string(number)) // convert number from byte to string to int (this is how many times we have to go back and change to upper)
			for a := 1; a <= n; a++ {
				linesSlice[index-a] = strings.ToUpper(s[i-a]) // change as many from the front as big is the n to UPPERCASE
			}
		} else if multipleLow.MatchString(word) { // MULTIPLE TO LOWER
			number := string(s[i][5])            // take the number from the word to see how many times we have to go back
			n, _ := strconv.Atoi(string(number)) // convert number from byte to string to int (this is how many times we have to go back and change to upper)
			for a := 1; a <= n; a++ {
				linesSlice[index-a] = strings.ToLower(s[i-a]) // change as many from the front as big is the n to lowercase
			}
		} else if multipleCap.MatchString(word) { // MULTIPLE TO CAP
			number := string(s[i][5])            // take the number from the word to see how many times we have to go back
			n, _ := strconv.Atoi(string(number)) // convert number from byte to string to int (this is how many times we have to go back and change to upper)
			for a := 1; a <= n; a++ {
				linesSlice[index-a] = strings.Title(s[i-a]) // change as many from the front as big is the n to Capitalize
			}
		} else if word == "(hex)" {
			linesSlice[index-1] = hexToDec(s[i-1])
		} else if word == "(bin)" {
			linesSlice[index-1] = binToDec(s[i-1])
		} else {
			linesSlice = append(linesSlice, s[i])
			index++
		}
	}
	return linesSlice
}

func readLine(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil { // deal with error
		fmt.Println(err)
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	linesSlice := []string{} // lines slice

	for scanner.Scan() {
		linesSlice = append(linesSlice, scanner.Text()) // read line by line from scanner
	}
	if err := scanner.Err(); err != nil { // deal with error
		log.Fatal(err)
	}
	return linesSlice
}

func writeFile(result string) {
	f, err := os.OpenFile(os.Args[2], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.WriteString(result + "\n"); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func makeIntoSlice(content string) []string { // make slice from the content string
	result := strings.Split(content, " ") // split the string on whitespaces
	return result
}

func hexToDec(hexadecimal_num string) string {
	decimalNum, err := strconv.ParseInt(hexadecimal_num, 16, 64) // use the parseInt() func for conversion
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(decimalNum))
}

func binToDec(bin_num string) string {
	decimalNum, err := strconv.ParseInt(bin_num, 2, 64) // use the parseInt() func for conversion
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(decimalNum))
}

func Connector(s []string) []string { // for example connect the "(cap," and  "2) "
	newSlice := []string{}
	tempString := ""
	matchSingle, _ := regexp.Compile(`\(cap\)|\(low\)|\(up\)`) // match singles
	matchMulti, _ := regexp.Compile(`\(cap,|\(low,|\(up,`)     // match first part of multiples assume that the second part is next

	for i := 0; i < len(s); i++ {
		if matchSingle.MatchString(s[i]) {
			newSlice = append(newSlice, s[i])
		} else if matchMulti.MatchString(s[i]) {
			tempString = s[i] + s[i+1]
			newSlice = append(newSlice, tempString)
			tempString = ""
			i++
		} else {
			newSlice = append(newSlice, s[i])
		}
	}
	return newSlice
}

func fixPunctuation(content string) string {
	allMultiSpaces := regexp.MustCompile(`\s+([\.\,\!\?\:\;"])`)   // choose all spaces and then puncts
	resultString := allMultiSpaces.ReplaceAllString(content, `$1`) // leave only the subgroup 1, jäta alles ainult grupp 1 ehk siis eemalda tühikud

	matches := regexp.MustCompile(`(\))([,.;:!?])`) // find all cases of ), ). ): etc and add a space in between
	regexed := matches.ReplaceAllString(resultString, `${1} ${2}`)

	find_without_spaces := regexp.MustCompile(`([\.\,\!\?\:\;"])(\w)`)        // choose all puncts and then characters (means they will need a space added to them)
	resultString = find_without_spaces.ReplaceAllString(regexed, `${1} ${2}`) // use replaceAll to put a space between group 1 and 2
	return resultString
}

func fixUpperCommas(content string) string {
	matchGroups := regexp.MustCompile(`(')(\s*)(.+?)(\s*)(')`)            // find and select all the subgroups
	resultString := matchGroups.ReplaceAllString(content, `${1}${3}${5}`) // leave out the subgroups 2 and 4 which are the unnecessary spaces.
	return resultString
}

func convertAtoAn(content string) string { // convert a to an where necessary
	matchGroups := regexp.MustCompile(`(A|a)( +)([AEIOUHaeiouh])`)      // find all cases where a or A is not correct
	resultString := matchGroups.ReplaceAllString(content, `${1}n ${3}`) // add n to the first capture group followed by a space and then the vowel that was matched
	return resultString
}
