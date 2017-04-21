# Govoice

Govoice creates your invoices only using yaml and django jinja templates.

## Installation

```
go get github.com/valkum/govoice
```

For PDF output please install electron-pdf using
```
npm install -g electron-pdf
```

## Config

Simply create a file named ´Govoice.yaml´ in your working dir, ´~/.govoice´ or in /etc/govoice

### Config Options

_________
#### pdf
***true***
> After the html is exported call electron-pdf to create a pdf with the same filename.

***false***
> Stop after creating the html file.

Default: false

---------
### template

This value is used as a templatefile for pongo2. Just point it to your desired jinja template file.

---------
### out_dir

Output directory based on your current working dir.

---------
### from

A yaml object containing your details.

```
  name: "Bruce Wayne"
  street: "Some Street in Gotham"
  zip: "1337"
  city: "A City"
  mobile: "0123456789"
```

## Usage

```
govoice invoice-1.yaml
```

where `invoice-1.yaml` is a file of the following form
```
date: "20.04.2017"
to:
  name: "Client Inc."
  street: "Clients Street"
  mobile: "Mobile Number"
  zip: "1337"
  city: "Client City"
items:
- id: 1
  description: "Item 1"
  pricePerUnit: 10
  quantity: 100
  total: 1000
- id: 2
  description: "Item 2 (Discounted)"
  pricePerUnit: 10
  quantity: 10
  total: 100
```
