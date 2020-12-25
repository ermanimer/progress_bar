package progress_bar

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestDefaultProgressBar(t *testing.T) {
	//create output
	output := new(bytes.Buffer)
	//create default progress bar
	pb := DefaultProgressBar(50)
	pb.Output = output
	//start
	expectedOutput := "[..................................................][%0.0][0.0/50.0][Elapsed: 0.0s Remaining: +Infs]"
	expectedByteCount := len(expectedOutput) //100
	err := start(pb, output, expectedByteCount, expectedOutput)
	if err != nil {
		t.Error(err.Error())
	}
	//update
	expectedOutput = "\u001b[2K\u001b[100D[#########################.........................][%50.0][25.0/50.0][Elapsed: 1.0s Remaining: 1.0s]"
	expectedByteCount = len(expectedOutput) //111
	err = update(pb, 25, output, expectedByteCount, expectedOutput)
	if err != nil {
		t.Error(err.Error())
	}
	//stop
	expectedOutput = "\u001b[2K\u001b[101D[##################################################][%100.0][50.0/50.0][Elapsed: 2.0s Remaining: 0.0s]\n" //101: 111 - 10(escape codes)
	expectedByteCount = len(expectedOutput)
	err = update(pb, 50, output, expectedByteCount, expectedOutput)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestNewProgressBar(t *testing.T) {
	//create output
	output := new(bytes.Buffer)
	//create custom parameters
	schema := fmt.Sprintf("(%s)(%s)(%s of %s)(E: %ss R: %ss)", svBar, svPercent, svCurrentValue, svTotalValue, svElapsedDuration, svRemainingDuration)
	filledCharacter := "="
	blankCharacter := "-"
	var length float64 = 60
	var value float64 = 80
	pb := NewProgressBar(output, schema, filledCharacter, blankCharacter, length, value)
	//start

	expectedOutput := "(------------------------------------------------------------)(%0.0)(0.0 of 80.0)(E: 0.0s R: +Infs)"
	expectedByteCount := len(expectedOutput) //99
	err := start(pb, output, expectedByteCount, expectedOutput)
	if err != nil {
		t.Error(err.Error())
	}
	//update
	expectedOutput = "\u001b[2K\u001b[99D(==============================------------------------------)(%50.0)(40.0 of 80.0)(E: 1.0s R: 1.0s)"
	expectedByteCount = len(expectedOutput) //109
	err = update(pb, 40, output, expectedByteCount, expectedOutput)
	if err != nil {
		t.Error(err.Error())
	}
	//stop
	expectedOutput = "\u001b[2K\u001b[100D(============================================================)(%100.0)(80.0 of 80.0)(E: 2.0s R: 0.0s)\n" //100: 109 - 9(escape codes)
	expectedByteCount = len(expectedOutput)
	err = update(pb, 80, output, expectedByteCount, expectedOutput)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestStartError(t *testing.T) {
	//create output
	output := new(bytes.Buffer)
	//create default progress bar
	pb := DefaultProgressBar(100)
	pb.Output = output
	//start
	err := pb.Start()
	if err != nil {
		t.Error("starting progress bar failed")
	}
	//start again
	err = pb.Start()
	if err == nil {
		t.Error("catching \"progress bar is already started\" error failed")
	}
}

func TestStopError(t *testing.T) {
	//create output
	output := new(bytes.Buffer)
	//create default progress bar
	pb := DefaultProgressBar(100)
	pb.Output = output
	//stop again
	err := pb.Stop()
	if err == nil {
		t.Error("catching \"progress bar is not started\" error failed")
	}
}

func TestUpdateErrors(t *testing.T) {
	//create output
	output := new(bytes.Buffer)
	//create default progress bar
	pb := DefaultProgressBar(100)
	pb.Output = output
	//update
	err := pb.Update(50)
	if err == nil {
		t.Error("catching \"progress bar is noty started\" error failed")
	}
	//start
	err = pb.Start()
	if err != nil {
		t.Error("starting progress bar failed")
	}
	//update with a value which is greater then total value
	err = pb.Update(101)
	if err == nil {
		t.Error("catching \"value is greater then total value")
	}
}

func start(pb *ProgressBar, output *bytes.Buffer, expectedByteCount int, expectedOutput string) error {
	err := pb.Start()
	if err != nil {
		return err
	}
	bs := make([]byte, expectedByteCount)
	n, err := output.Read(bs)
	if err != nil {
		return errors.New("reading output failed")
	}
	if n != expectedByteCount {
		return errors.New("byte count doesn't match")
	}
	if string(bs) != expectedOutput {
		return errors.New("output doesn't match")
	}
	return nil
}

func update(pb *ProgressBar, value float64, output *bytes.Buffer, expectedByteCount int, expectedOutput string) error {
	time.Sleep(1 * time.Second)
	err := pb.Update(value)
	if err != nil {
		return err
	}
	bs := make([]byte, expectedByteCount)
	n, err := output.Read(bs)
	if err != nil {
		return errors.New("reading output failed")
	}
	if n != expectedByteCount {
		return errors.New("byte count doesn't match")
	}
	if string(bs) != expectedOutput {
		return errors.New("output doesn't match")
	}
	return nil
}
