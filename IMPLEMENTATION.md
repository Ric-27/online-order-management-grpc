# Implementation in go for a backend using gRPC that describes an online store
## Ricardo RICO URIBE
## Tasks
### I - Order service
The current system exposes a product service that provides a fixed list of products through an RPC.

The goal here is to improve the order service in order to create orders based on the products:

1. Complete the proto of the `product` service to implement a RPC to retrieve a single product by its ID.
#### **R/**
The service is ProductOfId(id string) &Product. As expected it returns the full product information based on the id. 

2. Complete the proto of the `order` service and implement a RPC to create a new order. Order must have the following fields:
    - customer firstname
    - customer lastname
    - line items (products & quantities)
    - shipping address (destination)
        - address line (45 Rue des Petites Ecuries)
        - postal code (75010)
        - city (Paris)
        - country (FR)
#### **R/**
The service is CreateOrder(fName string, lName string, productList map[string]int, address Address) &Order
the suplementary types are Address(street string, postal string, city string, country string) and for productList the string key is a *valid* product id and the int is the quantity.

I specify valid because if an Order would not be created if a requested product is not available/doesn't exists.

The order is assigned an id, with the format [[:num:]]{4}-[[:date:]]{8}-[[:hour:]]{6}-[[:alnum:]]{}
the first set of values is a 4 random number to ensure uniqueness, the next two groups represent the date and time of the order ex.(november 21, 2009 at 09h57m40s, this will be transformed into 20091121-095740) the last group is the customer lastname in all caps.
### II - Address validation
As discussed previously, you should propose and implement a solution to validate the shipping address of an order before creating it:

- If the address contains some slight errors and the correct data can be identified with certainty by the system, the address will be automatically fixed and the order is created. Some examples:
    - 45 Rue des Pet**is** Ecuries → 45 Rue des Pet**ites** Ecuries
    - 1 Square Emile Z**i**la → 1 Square Emile Z**o**la
    - Par**i** → Par**is**
    - Aubervi**l**iers → Aubervi**ll**iers

-  Otherwise, if some parts of the address cannot be recognised and the system fails to validate it, the order is not created and a response with an error code is returned.

#### **R/**
The address correction is done with the address api, it automatically transforms an address to the closest match.

Because of the concept of the closest match, the api will always return an anwers, it would only give an error if the request is empty. This would allow poorly written address to be valid by the api. Thats why I used the score value (a value between 0 and 1 that represents the confidence the address api has in the result), if the score is less than 0.3, I consider that the address was written incorrectly and an order is not created

To the api we do not send the country information because is a french api, and we ignore this value from the user and save directly "france" as country, we could eliminate the country field from the service declaration
#
## General Info
The file *server_test.go* includes the unitary test for both services, it checks if a known action should throw an error or not.
## Possible Improvements
- The unitary test:
    - check if the address is corrected as we want (it is done but not tested)
- OrderCreation:
    - modify an existing order (core function needed)
    - erase invalid products and then create the order with the valid products (maybe not useful)