package ishell

import (
	"strings"
        "os"
	"io/ioutil"
	"github.com/flynn-archive/go-shlex"
)

type iCompleter struct {
	cmd      *Cmd
	disabled func() bool
}

func (ic iCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	if ic.disabled != nil && ic.disabled() {
		return nil, len(line)
	}
	var words []string
	if w, err := shlex.Split(string(line)); err == nil {
		words = w
	} else {
		// fall back
		words = strings.Fields(string(line))
	}

	var cWords []string
	prefix := ""
	if len(words) > 0 && pos > 0 && line[pos-1] != ' ' {
		prefix = words[len(words)-1]
		cWords = ic.getWords(words[:len(words)-1])
	} else {
		cWords = ic.getWords(words)
	}

	var suggestions [][]rune
	for _, w := range cWords {
		if strings.HasPrefix(w, prefix) {
			suggestions = append(suggestions, []rune(strings.TrimPrefix(w, prefix)))
		}
	}
	if len(suggestions) == 1 && prefix != "" && string(suggestions[0]) == "" {
		suggestions = [][]rune{[]rune(" ")}
	}
	return suggestions, len(prefix)
}

func (ic iCompleter) getWords(w []string) (s []string) {
	cmd, args := ic.cmd.FindCmd(w)
	if cmd == nil {
		if len(args)==0{
			cmd, args = ic.cmd, w
		}else{
			return listFileNames()
		}
	}
	if cmd.Completer != nil {
		return cmd.Completer(args)
	}
	for k := range cmd.children {
		s = append(s, k)
	}
	return
}
func listFileNames()[]string{
	var names []string
	dir,_ := os.Getwd()
	files, _ := ioutil.ReadDir(dir)
	for _,file := range files {
		names = append(names, file.Name())
	}
	return names
}
