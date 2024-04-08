package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/jung-kurt/gofpdf"
)

type Invoice struct {
	Heading            string        `json:"heading"`
	InvoiceFromEmail   string        `json:"invoice_from_email"`
	InvoiceFromAddress string        `json:"invoice_from_address"`
	BillingItems       []BillingItem `json:"billing_items"`
	BillTo             BillTo        `json:"bill_to"`
	Project            string        `json:"project"`
	InvoiceNumber      int           `json:"invoice_number"`
	Signature          string        `json:"signature"`
}

type BillTo struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Date    string `json:"date"`
}

type BillingItem struct {
	Description string `json:"decription"`
	Quantity    int    `json:"quantity"`
	UnitPrice   string `json:"unit_price"`
	Cost        string `json:"cost"`
}

func main() {

	handleSignals()

	// Parse the data from JSON file
	invoice, err := parseInvoice("data.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate PDF from the parsed data
	err = InvoicePDF(invoice, "invoice.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}

}

func parseInvoice(file string) (Invoice, error) {
	//open Json
	jsonFile, err := os.Open(file)
	// if Open return error
	if err != nil {
		fmt.Println(("Error opening JSON file:"), err)
	} else {
		fmt.Println("File opened successfully")
	}
	defer jsonFile.Close() // close the Json file

	// read the Json file
	byteRead, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return Invoice{}, err
	}

	// Create an invoice variable
	var invoice Invoice

	// Unmarshal the Json file to check if its in a correct format
	err = json.Unmarshal(byteRead, &invoice)
	if err != nil {
		fmt.Println("Wrong Json format", err) // Return the error if the Json file is not in a correct format
		return Invoice{}, err                 // return an empty instance of Invoice struct, along with the error
	}
	return invoice, nil // return the invoice and nill for no error
}

func InvoicePDF(invoice Invoice, file string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// 0. Create a line atop the heading

	pdf.SetDrawColor(52, 124, 162) // Set L color to cyan

	_, pageHeight := pdf.GetPageSize() // Get the page width

	pdf.SetLineWidth(1) // Set the line width to 1

	pdf.Line(0, 7, pageHeight, 7) // Draw a line at the top of the page

	pdf.SetLineWidth(0.2) //Reset the line width for table creation later

	// 1. Add Heading

	pdf.SetFont("Arial", "B", 20) // Set the heading / title for the document

	pdf.SetTextColor(52, 124, 162) // Set fill color to cyan

	pdf.CellFormat(190, 10, invoice.Heading, "", 0, "L", false, 0, "") // Add the heading from json file's data

	// Reset text color to black and fill color to white
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(50)
	pdf.SetFillColor(255, 255, 255)

	// 2. Add Seller Info

	//Get the current Y position of Seller so we can position Buyer to be at the same line as Seller
	y := pdf.GetY()

	// Set font & text color to cyan
	pdf.SetTextColor(52, 124, 162)
	pdf.SetFont("Arial", "", 20)
	pdf.CellFormat(90, 10, "INVOICE", "", 0, "L", false, 0, "")

	pdf.SetTextColor(0, 0, 0) // Reset text color to black
	pdf.Ln(15)

	// Set font & text color to cyan
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(52, 124, 162)
	pdf.CellFormat(90, 10, invoice.InvoiceFromEmail, "", 0, "L", false, 0, "")

	pdf.SetTextColor(0, 0, 0) // Reset text color to black
	pdf.Ln(10)
	pdf.MultiCell(40, 8, invoice.InvoiceFromAddress, "", "L", false) // Add Seller's address

	// 3. Add Buyer Info

	pdf.SetXY(60, y) // Format BUYER INFO at the same line as Seller Info
	pdf.CellFormat(90, 10, "Bill To:"+invoice.BillTo.Name, "", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Additional Buyer Info
	pdf.CellFormat(120, 10, invoice.BillTo.Address, "", 0, "R", false, 0, "")
	pdf.Ln(10)
	pdf.CellFormat(78, 10, "Date:"+invoice.BillTo.Date, "", 0, "R", false, 0, "")
	pdf.Ln(20)
	pdf.CellFormat(92, 10, "Project: "+invoice.Project, "", 0, "R", false, 0, "")
	pdf.Ln(10)
	pdf.CellFormat(87, 10, "Invoice Number: "+strconv.Itoa(invoice.InvoiceNumber), "", 0, "R", false, 0, "")
	pdf.Ln(10)

	//4. Add table headers
	headers := []string{"Description", "Quantity", "Unit Price", "Cost"}

	pdf.SetFillColor(52, 124, 162)  // Set fill color to cyan
	pdf.SetTextColor(255, 255, 255) // Set text color to white

	pdf.SetX(59)                    // Move header to the right
	pdf.SetDrawColor(192, 192, 192) // Set the cell border to grey
	// Add headers to the PDF
	for _, header := range headers {
		pdf.CellFormat(38, 10, header, "1", 0, "L", true, 0, "")

	}
	pdf.Ln(10)

	// Reset fill color to white (default) and text color to black for the rest of the document
	pdf.SetFillColor(255, 255, 255)
	pdf.SetTextColor(0, 0, 0)

	// Set draw color to grey
	pdf.SetDrawColor(192, 192, 192)

	// Initialize total cost
	totalCost := 0

	//Add a row for each item
	for _, item := range invoice.BillingItems {
		pdf.SetX(59)
		pdf.CellFormat(38, 10, item.Description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(38, 10, strconv.Itoa(item.Quantity), "1", 0, "R", false, 0, "")
		pdf.CellFormat(38, 10, item.UnitPrice, "1", 0, "R", false, 0, "")
		pdf.CellFormat(38, 10, item.Cost, "1", 0, "R", false, 0, "")
		pdf.Ln(10)

		// Remove dollar sign from item cost
		itemCostStr := strings.Replace(item.Cost, "$", "", -1)

		// Convert item cost to integer and add to total cost
		itemCost, err := strconv.Atoi(itemCostStr)
		if err != nil {
			// Handle error
			fmt.Println("Error converting item cost to integer:", err)
			continue
		}
		totalCost += itemCost

	}

	pdf.SetX(59)

	// Set draw color to black
	pdf.SetDrawColor(0, 0, 0)

	// Add TOTAL section
	pdf.CellFormat(114, 10, "TOTAL", "1", 0, "R", false, 0, "")
	pdf.CellFormat(38, 10, "$"+strconv.Itoa(totalCost), "1", 0, "R", false, 0, "")
	pdf.Ln(20)

	pdf.SetX(59)
	//5. Thank you message & Signature
	pdf.MultiCell(150, 8, `Thank you for your business. Its a pleasure to work with you on your project. 
	Your next order will ship in 30 days`, "", "L", false)
	pdf.SetX(59)
	pdf.CellFormat(159, 10, "Sincerely yours,", "", 0, "L", false, 0, "")
	pdf.Ln(10)
	pdf.SetX(59)
	pdf.CellFormat(50, 10, "Signature: "+invoice.Signature, "", 0, "L", false, 0, "")

	// Add a row for each item
	err := pdf.OutputFileAndClose(file)
	if err != nil {
		return err
	}
	return nil
}

func handleSignals() {

	sigs := make(chan os.Signal, 1) // Receive signal

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // receive notification of the specified signals

	// Start a goroutine that will do something when a signal is received.
	go func() {
		sig := <-sigs // assign signal to variable sig
		fmt.Println(sig)
		os.Exit(0)
	}()
}
