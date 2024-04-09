# Golang-invoice
an invoice program in Golang that take .json file and output an invoice PDF

Have the json data with the format:

1. Top-Level FIelds:
'heading' : A descriptive string represending the name or title of the invoice
'email' : provide the email address from where the invoice is being sent
'From address' : provide the address from where the invoice is being sent
'bill to' ( Object ): provide the 'name','address','date' of the recipient

2. Additional Fields:
'project' : specify the project's name in string
'invoice_number' : provide the invoice number as an integer value
'billing_items' ( Array / Slice ) :
- description : A descriptive string that describe the item
- quantity : An int provide the quantity of the billing item
- unit_price : A string that provide the currency & price of the item
- cost : A string that provide the total cost of the billing's items ( multiple unit_price to quantity )
'signature' : A string that provide the signature associated with the invoice's sender

Sample of a correct Json input:
{
  "heading": "DWARVES FOUNDATION",
  "invoice_from_email": "huy@dwarvesv.com",
  "invoice_from_address": "1234 Main Street Anytown, State ZIP",
  "bill_to": {
    "name": "Mr. John",
    "address": "4321 First Street Anytown, State ZIP",
    "date": "12/28/17"
  },
  "project": "Project Name",
  "invoice_number": 10,
  "billing_items": [
    {
      "decription": "item 1",
      "quantity": 10,
      "unit_price": "$9",
      "cost": "$90"
    },
    {
      "decription": "item 2",
      "quantity": 5,
      "unit_price": "$10",
      "cost": "$50"
    }
  ],
  "signature": "Huy Giang"
}

Note when writing Json input:
- Ensure that each field is enclosed in double quotes (")
- Use curly brace {} to define "object".
- Separate fields and values with a colon :
- Separate items in an array with a comma ;
- Use square brackets [] to define arrays


