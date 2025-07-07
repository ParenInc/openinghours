# Opening Hours

A Go package for handling recurring weekly opening hours in a standardized format.

## Overview

This package provides functionality to manage and parse opening hours for amenities, businesses, or any entity that operates on a weekly schedule. It implements a custom time format since neither RFC 3339 nor ISO 8601 provides a standard way to represent recurring times within a week.

## Installation
```bash
go get github.com/pareninc/openinghours
```

## Features
- Parse and format weekly opening hours
- Convert between machine-readable and human-readable formats
- Support for multiple opening periods per day
- Handles overnight and multi-day periods
- RFC 3339 compliant weekday numbering (Monday = 1, Sunday = 7)

## Usage
### Basic Example

```go
import "github.com/pareninc/openinghours"

// Parse opening hours from string
hours, err := openinghours.ParseOpeningHours("W2T06:00:00/W2T20:00:00,W5T10:30:00/W5T13:00:00")
if err != nil {
    log.Fatal(err)
}

// Convert to human-readable format
humanReadable, err := openinghours.GetHumanReadableTimes("W3T10:00:00/W3T20:30:00,W5T10:00:00/W5T12:00:00")
if err != nil {
    log.Fatal(err)
}
// humanReadable := map[string][]string{
//  "wednesday": []string{"open: 10:00, close: 20:30"},
//  "friday":    []string{"open: 10:00, close: 12:00"},
// }
```

### Time Format
The package uses a custom string format for representing opening hours:
- Format: `W<day>T<HH>:<MM>:00/W<day>T<HH>:<MM>:00`
- Where:
    - `<day>` is 1-7 (1 = Monday, 7 = Sunday)
    - `<HH>` is hours in 24-hour format (00-24)
    - `<MM>` is minutes (00-59)

- Multiple periods are separated by commas

Example: represents Tuesday from 6:00 AM to 8:00 PM `"W2T06:00:00/W2T20:00:00"`

### Working with Opening Hours
The package provides two main types:
- `OpeningHours`: Contains opening and closing times
- `TimeInWeek`: Represents a specific time within a week

## License
MIT License - see [LICENSE](LICENSE) for details

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.
