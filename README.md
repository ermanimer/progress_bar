# progress_bar
Go ASCII Progress Bar

![Go](https://github.com/ermanimer/progress_bar/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/ermanimer/progress_bar)](https://goreportcard.com/report/github.com/ermanimer/progress_bar)

## Features
progress_bar creates a single customizable progress bar for terminal.

## Installation
```bash
go get -u github.com/ermanimer/progress_bar
```

## Functions
#### DefaultProgressBar
Creates a progress bar with default parameters and given total value.

|Parameter |Data Type|
|:---------|:--------|
|totalValue|float64  |

Default parameters:

|Default Parameter       |Value                                                                            |
|:-----------------------|:--------------------------------------------------------------------------------|
|Default Schema          |[{bar}][{percent}][{current}/{total}][Elapsed: {elapsed}s Remaining:{remaining}s]|
|Default Filled Character|#                                                                                |
|Default Blank Character |.                                                                                |
|Default Length          |50                                                                               |

#### NewProgressBar
Creates a progress bar with default parameters and given total value.

|Parameter      |Data Type|
|:--------------|:--------|
|output         |io.Writer|
|schema         |string   |
|filledCharacter|float64  |
|blankCharacter |float64  |
|length         |float64  |
|totalValue     |float64  |

Schema Variables:

|Schema Variable|Value                                                                            |
|:--------------|:----------------------------|
|{bar}          |Bar of progress bar          |
|{percent}      |Percentege of progress bar   |
|{current}      |Current value of progress bar|
|{total}        |Total value of progress bar  |
|{elapsed}      |Elapsed duration             |
|{remaining}    |Estimated remaining duration |

## Methods
#### Start
Starts progress bar.

#### Stop
Stops progress bar.

#### Update
Updates progress bar with given value and stops progress bar is total value is reached.

|Parameter|Data Type|
|:--------|:--------|
|value    |float64  |

## Usage
Default Progress Bar:

```go
package main

import (
	"fmt"
	"time"

	"github.com/ermanimer/progress_bar"
)

func main() {
	//create new progress bar
	pb := progress_bar.DefaultProgressBar(100)
	//start
	err := pb.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//update
	for value := 1; value <= 100; value++ {
		time.Sleep(20 * time.Millisecond)
		err := pb.Update(float64(value))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}
}
```

Terminal Output:
![Default Terminal Output](/images/default.gif)

New Progress Bar:

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ermanimer/progress_bar"
)

func main() {
	//create parameters
	output := os.Stdout
	schema := "({bar}) ({percent}) ({current} of {total} completed)"
	filledCharacter := "="
	blankCharacter := "-"
	var length float64 = 60
	var totalValue float64 = 80
	//create new progress bar
	pb := progress_bar.NewProgressBar(output, schema, filledCharacter, blankCharacter, length, totalValue)
	//start
	err := pb.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//update
	for value := 1; value <= 80; value++ {
		time.Sleep(20 * time.Millisecond)
		err := pb.Update(float64(value))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}
}
```

Terminal Output:
![New Terminal Output](/images/new.gif)

New Progress Bar Colored With [color](https://github.com/ermanimer/color) On Linux:

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ermanimer/color/v2"
	"github.com/ermanimer/progress_bar"
)

func main() {
	//create parameters
	output := os.Stdout
	//create color functions
	orange := (&color.Color{Foreground: 172}).SprintFunction()
	grey := (&color.Color{Foreground: 246}).SprintFunction()
	//create colored schema
	bar := orange("{bar}")
	percent := grey("{percent}")
	schema := fmt.Sprintf("%s %s", bar, percent)
	filledCharacter := "▆"
	blankCharacter := " "
	var length float64 = 50
	var totalValue float64 = 80
	//create new progress bar
	pb := progress_bar.NewProgressBar(output, schema, filledCharacter, blankCharacter, length, totalValue)
	//start
	err := pb.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//update
	for value := 1; value <= 80; value++ {
		time.Sleep(20 * time.Millisecond)
		err := pb.Update(float64(value))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}
}
```

Terminal Output:
![New Terminal Output](/images/colored.gif)

## References
- [ANSI Escape Codes, Wikipedia](https://en.wikipedia.org/wiki/ANSI_escape_code)
- [Build Your Own Command Line With ANSI Escape Codes, Haoyi's Programming Blog](https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html)