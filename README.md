<h1 align="center">InvestMate</h1>
<p align="center">
    <a href="https://github.com/rios0rios0/investmate/releases/latest">
        <img src="https://img.shields.io/github/release/rios0rios0/investmate.svg?style=for-the-badge&logo=github" alt="Latest Release"/></a>
    <a href="https://github.com/rios0rios0/investmate/blob/main/LICENSE">
        <img src="https://img.shields.io/github/license/rios0rios0/investmate.svg?style=for-the-badge&logo=github" alt="License"/></a>
</p>

A Go application that scrapes and processes ETF (Exchange-Traded Fund) data, including dividend cash amounts, average closing prices, and dividend yields over a specified number of years. The data is then displayed in a formatted table.

## Features

- Scrapes dividend cash amounts for specified ETFs
- Fetches average closing prices for specified ETFs
- Calculates and displays dividend yields
- Displays data in a formatted table with color-coded dividend yields

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/rios0rios0/investmate.git
   cd investmate
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

## Usage

Run the application:

```sh
go run main.go
```

The application will scrape data for the specified ETFs and display it in a formatted table in the console.

## Configuration

- **Years to Fetch:**
  You can configure the number of years to fetch data for by changing the `YearsToFetch` constant in the code:
  ```go
  const (
      YearsToFetch = 5 // Number of years to fetch data for
  )
  ```

- **ETF Names:**
  You can specify the ETFs to scrape by modifying the `etfNames` slice in the `main` function:
  ```go
  etfNames := []string{"HYGW", "RIET", "SDIV", "SVOL", "XYLD"}
  ```

## Code Structure

- **ETF Struct:**
  Represents an ETF and its data.
  ```go
  type ETF struct {
      Name                       string
      AmountDividendsPerYear     map[string]float64
      AverageClosingPricePerYear map[string]float64
      DividendYieldPerYear       map[string]float64
  }
  ```

- **Functions:**
    - `ShowDividendsPerYear`: Formats yearly dividend sums for display
    - `AverageDividends`: Calculates average dividends
    - `ShowClosingPricesPerYear`: Formats average closing prices for display
    - `AverageClosingPrices`: Calculates average closing prices
    - `ShowDividendYieldPerYear`: Calculates and formats dividend yields
    - `AverageDividendYield`: Calculates average dividend yield
    - `processETF`: Populates ETF data
    - `crawlingDividendsPerYear`: Scrapes dividend data
    - `fetchAverageClosingPricesPerYear`: Fetches average closing prices
    - `getColors`: Returns colors for the dividend yield row

## Dependencies

- `gocolly/colly` - Scraping framework for Go
- `olekukonko/tablewriter` - Library for rendering ASCII tables in Go
- `sirupsen/logrus` - Structured logger for Go

## Contributing

Contributions are welcome. Please open an issue or submit a pull request for improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
