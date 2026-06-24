package tuix

//--------USE ANYWHERE-
// tuix.Debug("Render Dashboard")
// tuix.Debug("Focused:", focused)
// tuix.Debug("Current Route:", route)
//----------
import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	f, _ := os.OpenFile(
		"tuix.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)

	logger = log.New(f, "", log.LstdFlags)
}

func Debug(v ...any) {
	logger.Println(v...)
}
