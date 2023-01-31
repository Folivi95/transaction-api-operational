# Transaction API Operational

This repo contains three microservices related to the transaction api:

- "Luigi": transformer that consumes messages from the raw topics and converts them to the canonical format
- "Mario": api to serve the messages from the Transactions DB
- "Kamek": cronjob to keep only 6 months of data from the Transactions DB

