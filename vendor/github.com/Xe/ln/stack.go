package ln

import (
	"os"
	"runtime"
	"strings"
)

type frame struct {
	filename string
	function string
	lineno   int
}

// skips 2 frames, since Caller returns the current frame, and we need
// the caller's caller.
func callersFrame() frame {
	var out frame
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return out
	}
	srcLoc := strings.LastIndex(file, "/src/")
	if srcLoc >= 0 {
		file = file[srcLoc+5:]
	}
	out.filename = file
	out.function = functionName(pc)
	out.lineno = line

	return out
}

func functionName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "???"
	}
	name := fn.Name()
	beg := strings.LastIndex(name, string(os.PathSeparator))
	return name[beg+1:]
	//	end := strings.LastIndex(name, string(os.PathSeparator))
	//	return name[end+1 : len(name)]
}
