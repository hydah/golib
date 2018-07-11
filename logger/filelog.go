package logger

import (
	"fmt"
	"os"
	"time"
)

// FileLogWriter : This log writer sends output to a file
type FileLogWriter struct {
	rec chan *LogRecord
	rot chan bool

	// The opened file
	filename string
	file     *os.File

	// The logging format
	format string

	// File header/trailer
	header, trailer string

	// Rotate at linecount
	maxlines         int
	maxlinesCurlines int

	// Rotate at size
	maxsize        int
	maxsizeCursize int

	// Rotate daily
	daily         bool
	dailyOpendate int

	// Keep old logfiles (.001, .002, etc)
	rotate bool
	suffix bool
}

var fileCloseSignal chan bool

// LogWrite : This is the FileLogWriter's output method
func (w *FileLogWriter) LogWrite(rec *LogRecord) {
	select {
	case w.rec <- rec:
	default:
	}
}

// Close _
func (w *FileLogWriter) Close() {
	close(w.rec)
	t := time.NewTimer(1 * time.Second)
	select {
	case <-fileCloseSignal:
		return
	case <-t.C:
		return
	}
}

// NewFileLogWriter creates a new LogWriter which writes to the given file and
// has rotation enabled if rotate is true.
//
// If rotate is true, any time a new log file is opened, the old one is renamed
// with a .### extension to preserve it.  The various Set* methods can be used
// to configure log rotation based on lines, size, and daily.
//
// The standard log-line format is:
//   [%D %T] [%L] (%S) %M
func NewFileLogWriter(fname string, rotate bool) *FileLogWriter {
	w := &FileLogWriter{
		rec:      make(chan *LogRecord, LogBufferLength),
		rot:      make(chan bool),
		filename: fname,
		format:   "[%D %T] [%L] (%S) %M",
		rotate:   rotate,
		suffix:   true,
	}
	fileCloseSignal = make(chan bool)
	// open the file for the first time
	if err := w.intRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
		return nil
	}

	go func() {
		defer func() {
			if w.file != nil {
				fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
				w.file.Close()
			}
			fileCloseSignal <- true
		}()

		for {
			select {
			case <-w.rot:
				if err := w.intRotate(); err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
					return
				}
			case rec, ok := <-w.rec:
				if !ok {
					return
				}
				now := time.Now()
				if (w.maxlines > 0 && w.maxlinesCurlines >= w.maxlines) ||
					(w.maxsize > 0 && w.maxsizeCursize >= w.maxsize) ||
					(w.daily && now.Day() != w.dailyOpendate) {
					if err := w.intRotate(); err != nil {
						fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
						return
					}
				}

				// Perform the write
				n, err := fmt.Fprint(w.file, FormatLogRecord(w.format, rec))
				if err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
					return
				}

				// Update the counts
				w.maxlinesCurlines++
				w.maxsizeCursize += n
			}
		}
	}()

	return w
}

// Rotate : Request that the logs rotate
func (w *FileLogWriter) Rotate() {
	w.rot <- true
}

// If this is called in a threaded context, it MUST be synchronized
func (w *FileLogWriter) intRotate() error {
	// Close any log file that may be open
	if w.file != nil {
		fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
		w.file.Close()
	}
	var filename string
	prefix := w.filename
	if w.suffix {
		filename = w.getCurrLogName(prefix, time.Now())
	} else {
		if len(prefix) == 0 {
			fmt.Fprintf(os.Stderr, "refuse suffix with log file and log file prefix is empty")
			os.Exit(1)
		} else {
			filename = prefix
		}
	}

	// If we are keeping log files, move it to the next available number
	now := time.Now()
	if w.rotate && (w.daily && now.Day() != w.dailyOpendate) {
		_, err := os.Lstat(filename)
		if err == nil { // file exists
			// Find the next available number
			num := 1
			fname := ""
			for ; err == nil && num <= 999; num++ {
				fname = filename + fmt.Sprintf(".%03d", num)
				_, err = os.Lstat(fname)
			}
			// return error if the last file checked still existed
			if err == nil {
				return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", filename)
			}

			// Rename the file to its newfound home
			err = os.Rename(filename, fname)
			if err != nil {
				return fmt.Errorf("Rotate: %s\n", err)
			}
		}
	}

	// Open the log file
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	w.file = fd
	fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: now}))

	// Set the daily open date to the current date
	w.dailyOpendate = now.Day()

	// initialize rotation values
	w.maxlinesCurlines = 0
	w.maxsizeCursize = 0

	return nil
}

// SetSuffix : Set the logging file suffix,Must be called before the first log
// message is written.
func (w *FileLogWriter) SetSuffix(suffix bool) *FileLogWriter {
	w.suffix = suffix
	return w
}

// SetFormat : Set the logging format (chainable).  Must be called before the first log
// message is written.
func (w *FileLogWriter) SetFormat(format string) *FileLogWriter {
	w.format = format
	return w
}

// SetHeadFoot : Set the logfile header and footer (chainable).  Must be called before the first log
// message is written.  These are formatted similar to the FormatLogRecord (e.g.
// you can use %D and %T in your header/footer for date and time).
func (w *FileLogWriter) SetHeadFoot(head, foot string) *FileLogWriter {
	w.header, w.trailer = head, foot
	if w.maxlinesCurlines == 0 {
		fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: time.Now()}))
	}
	return w
}

// SetRotateLines : Set rotate at linecount (chainable). Must be called before the first log
// message is written.
func (w *FileLogWriter) SetRotateLines(maxlines int) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateLines: %v\n", maxlines)
	w.maxlines = maxlines
	return w
}

// SetRotateSize : Set rotate at size (chainable). Must be called before the first log message
// is written.
func (w *FileLogWriter) SetRotateSize(maxsize int) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateSize: %v\n", maxsize)
	w.maxsize = maxsize
	return w
}

// SetRotateDaily : Set rotate daily (chainable). Must be called before the first log message is written.
func (w *FileLogWriter) SetRotateDaily(daily bool) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateDaily: %v\n", daily)
	w.daily = daily
	return w
}

// SetRotate changes whether or not the old logs are kept. (chainable) Must be
// called before the first log message is written.  If rotate is false, the
// files are overwritten; otherwise, they are rotated to another file before the
// new log is opened.
func (w *FileLogWriter) SetRotate(rotate bool) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotate: %v\n", rotate)
	w.rotate = rotate
	return w
}

// getCurrLogName _
func (w *FileLogWriter) getCurrLogName(prefix string, now time.Time) string {
	return fmt.Sprintf("%s-%d%02d%02d.log", prefix, now.Year(), now.Month(), now.Day())
}
