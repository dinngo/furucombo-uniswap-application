# FURUCOMBO Affected Addresses for Uniswap Retroactive Airdrop

The script exports 57 [affected addresses](affected_addresses.json) that use Uniswap through FURUCOMBO until 2020/9/1 12:00 am UTC, and has omitted addresses that were already in the initial retroactive airdrop by [Uniswap API](https://gentle-frost-9e74.uniswap.workers.dev/1/<USER_ADDRESS>).

## 1. Collect transaction hash by Dune Analytics
  Firstly create all [txs](txs.txt) from FURUCOMBO to Uniswap, the list is derived from https://explore.duneanalytics.com/public/dashboards/Cw10tmbxQ09KemdZdbzkW9rFMWTsBdE3CWW27jIF

## 2. Run
  Secondly collect senders of txs and filter out addresses by Uniswap API
  ```sh
  ENDPOINT=https://mainnet.infura.io/v3/<YOUR_INFURA_PROJECT_ID> go run main.go
  ```