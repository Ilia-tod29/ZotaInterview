# ZotaInteriew

ZotaInterview is a back-end api project, which provides endpoints for calling the ZOTA API.
This project is and interview task that aims to show skills for calling third-party APIs.

## Endpoints
- POST(“/deposit”) - Creates a new deposit
  - Required parameters: "merchantOrderDesc", "orderAmount", "orderCurrency", "customerEmail", "customerFirstName", 
    "customerLastName", "customerAddress", "customerCountryCode", "customerCity", "customerZipCode", 
    "customerPhone, "customerIP", "redirectUrl" "callbackUrl"
  - For more information about the parameters, please refer to: https://doc.zota.com/deposit/1.0/#deposit-request
- GET(“/status”) - Checks the status of an existing deposit
  - Required query parameters: "merchantOrderID", "orderID"
  - For more information about the parameters, please refer to: https://doc.zota.com/deposit/1.0/#order-status-request
- GET(“/payment-return”) - Simply displays the result after finished deposit. After a deposit is finished the user is being redirected to this endpoint.

## How to run

- Prerequisites:
  - go1.21.1

- Run the following command(make sure you are in the folder of project) in order:
```bash
make server
```

## Test

- Run the unit tests in each directory using:
```bash
go test -v
```

## Development

- To make env changes ether do them in the [app.env file](https://github.com/Ilia-tod29/ZotaInterview/blob/main/app.env)
  or export the declared in the file variables in your local machine