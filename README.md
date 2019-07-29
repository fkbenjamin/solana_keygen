# solana_keygen
Generates keypair.json for Solana CLI in Go (using mnemonics)

## Usage
If you want to generate a new keypair, simply run  

````
go run main.go
````

It will output a new, randomly generated mnemonic and create a keypair.json named after the address of the new key.

### Recover

If you want to recover a Keypair from an existing mnemonic, simply run

````
go run main.go -r
````

It will ask you to enter your mnemonic and create a keypair.json named after the address of the recovered key.
