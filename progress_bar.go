package progress_bar

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

//default schema
const (
	defaultSchema = "[%s][%s][%s/%s][Elapsed: %ss Remaining: %ss]"
)

//schema variables
const (
	svBar               = "{bar}"
	svPercent           = "{percent}"
	svCurrentValue      = "{current}"
	svTotalValue        = "{total}"
	svElapsedDuration   = "{elapsed}"
	svRemainingDuration = "{remaining}"
)

//default parameters
const (
	defaultFilledCharacter = "#"
	defaultBlankCharacter  = "."
	defaultLength          = 50
)

//escape codes
const (
	ecClearLine      = "\u001b[2K"
	ecMovecursorLeft = "\u001b[%dD"
)

type ProgressBar struct {
	Output            io.Writer
	Schema            string
	FilledCharacter   string
	BlankCharacter    string
	Length            float64
	CurrentValue      float64
	TotalValue        float64
	ElapsedDuration   float64
	RemainingDuration float64
	isStarted         bool
	startingTime      time.Time
	offset            int
}

func DefaultProgressBar(totalValue float64) *ProgressBar {
	return &ProgressBar{
		Output:          os.Stdout,
		Schema:          fmt.Sprintf(defaultSchema, svBar, svPercent, svCurrentValue, svTotalValue, svElapsedDuration, svRemainingDuration),
		FilledCharacter: defaultFilledCharacter,
		BlankCharacter:  defaultBlankCharacter,
		Length:          defaultLength,
		TotalValue:      totalValue,
	}
}

func NewProgressBar(output io.Writer, schema string, filledCharacter string, blankCharacter string, length float64, totalValue float64) *ProgressBar {
	return &ProgressBar{
		Output:          output,
		Schema:          schema,
		FilledCharacter: filledCharacter,
		BlankCharacter:  blankCharacter,
		Length:          length,
		TotalValue:      totalValue,
	}
}

func (pb *ProgressBar) Start() error {
	if pb.isStarted {
		return errors.New("progress bar is already started")
	}
	pb.isStarted = true
	pb.startingTime = time.Now()
	pb.print()
	return nil
}

func (pb *ProgressBar) Stop() error {
	if !pb.isStarted {
		return errors.New("progress bar is not started")
	}
	pb.CurrentValue = 0
	pb.isStarted = false
	pb.offset = 0
	fmt.Fprintln(pb.Output, "")
	return nil
}

func (pb *ProgressBar) Update(value float64) error {
	if !pb.isStarted {
		return errors.New("prograss bar is not started")
	}
	if value > pb.TotalValue {
		return errors.New("value is greater then total value")
	}
	pb.CurrentValue = value
	pb.print()
	return nil
}

func (pb *ProgressBar) print() {
	//create bar
	filledCharacterCount := int(pb.Length * pb.CurrentValue / pb.TotalValue)
	blankCharacterCount := int(pb.Length) - filledCharacterCount
	filledCharacters := strings.Repeat(pb.FilledCharacter, filledCharacterCount)
	blankCharacters := strings.Repeat(pb.BlankCharacter, blankCharacterCount)
	bar := fmt.Sprintf("%s%s", filledCharacters, blankCharacters)
	//calculate percentage
	percent := pb.CurrentValue / pb.TotalValue * 100
	//calculate elapsed and remaining durations
	pb.ElapsedDuration = time.Since(pb.startingTime).Seconds()
	remainingValue := pb.TotalValue - pb.CurrentValue
	pb.RemainingDuration = remainingValue * pb.ElapsedDuration / pb.CurrentValue
	//create progress bar
	progressBar := strings.Replace(pb.Schema, svBar, bar, 1)
	progressBar = strings.Replace(progressBar, svPercent, fmt.Sprintf("%%%.1f", percent), 1)
	progressBar = strings.Replace(progressBar, svCurrentValue, fmt.Sprintf("%.1f", pb.CurrentValue), 1)
	progressBar = strings.Replace(progressBar, svTotalValue, fmt.Sprintf("%.1f", pb.TotalValue), 1)
	progressBar = strings.Replace(progressBar, svElapsedDuration, fmt.Sprintf("%.1f", pb.ElapsedDuration), 1)
	progressBar = strings.Replace(progressBar, svRemainingDuration, fmt.Sprintf("%.1f", pb.RemainingDuration), 1)
	//clear line and offset cursor
	if pb.offset > 0 {
		fmt.Fprint(pb.Output, ecClearLine)
		fmt.Fprintf(pb.Output, ecMovecursorLeft, pb.offset)
	}
	//print progress bar
	fmt.Fprint(pb.Output, progressBar)
	pb.offset = len(progressBar)
	//stop
	if pb.CurrentValue == pb.TotalValue {
		pb.Stop()
	}
}
