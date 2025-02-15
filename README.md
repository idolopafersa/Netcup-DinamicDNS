# Netcup-DinamicDNS

A Go-based tool to update Netcup DNS A records dynamically, ensuring they match your current public IP.

## Prerequisites
- Go installed on your system
- A Netcup account with API access

## Setup Instructions

### 1. Configure the `.env` File
Create a `.env` file in the project directory and add your Netcup API credentials and domain information. You can find your API key and password in the [Netcup Customer Control Panel (CCP)](https://www.customercontrolpanel.de/).

```
customernumber=YOUR_CUSTOMER_NUMBER
apikey=YOUR_API_KEY
apipassword=YOUR_API_PASSWORD
domain=YOUR_DOMAIN
```

### 2. Compile the Program
Run the following command to build the executable:

```
go build -o netcup-updater main/main.go
```

### 3. Execute the Program
Run the compiled binary manually using:

```
./netcup-updater
```

To automate execution every 5 minutes on Linux, add a cron job:

1. Open the crontab editor:
   ```
   crontab -e
   ```
2. Add the following line at the end:
   ```
   */5 * * * * /path/to/netcup-updater
   ```
   Replace `/path/to/netcup-updater` with the actual path to the compiled binary.

This ensures that your DNS records are updated automatically every 5 minutes.

## License
This project is licensed under the MIT License.

## Contributing
Feel free to submit pull requests or open issues to improve this tool!

