package main

import(
  "bufio"
  "golang.org/x/crypto/ed25519"
  "flag"
  "fmt"
  "github.com/dignifiedquire/go-basex"
  bip39 "github.com/tyler-smith/go-bip39"
  "encoding/hex"
  "io/ioutil"
  "os"
  "strconv"
)

//Alphabet used for bs58 encoding
var BASE_58_CHARSET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func main() {
    //Flag to recover a key.json from existing mnemonic
   recoverPtr := flag.Bool("r", false, "recover from mnemonic")
   flag.Parse()
   var mnemonic string
   if(*recoverPtr) {
     mnemonic = recoverMnemonic()
   } else {
     mnemonic = ""
   }
     key := generateKey(mnemonic)
     writeKeyToFile(key)
 }

 //Generates an ed25519 keypair from mnemonic, if no mnemonic is given as argument, creates a new one and prints it to cli.
 func generateKey(mnemonic string) ed25519.PrivateKey {
   //if no mnemonic is set, create a new one
   if(mnemonic == "") {
     entropy, _ := bip39.NewEntropy(128)
     mnemonic, _ = bip39.NewMnemonic(entropy)
     fmt.Println("#############################################################################")
     fmt.Println(mnemonic)
     fmt.Println("#############################################################################")
   }
   entropy, err := bip39.EntropyFromMnemonic(mnemonic)
   if err != nil {
     fmt.Println("Mnenomic not valid", err)
   }
   key := ed25519.NewKeyFromSeed([]byte(hex.EncodeToString(entropy)))
   return key
 }

 //Prints out the bs58 encoded address of an ed25519 keypair
 func getAddress(key ed25519.PrivateKey) string {
   var pubBytes [32]byte
 	 copy(pubBytes[:], key[32:])
 	 bs58 := basex.NewAlphabet(BASE_58_CHARSET)
   enc := bs58.Encode(pubBytes[:])
   return enc
 }

 //Simple request to input mnemonic into cli
 func getMnemonicFromUser() string {
   var mnemonic string
   fmt.Print("Enter mnemonic to recover from: ")
   reader := bufio.NewReader(os.Stdin)
   mnemonic, _ = reader.ReadString('\n')
   return mnemonic
 }

 //Gets mnemonic from user input and verifies if it is valid. If not, asks user again and prints error
 func recoverMnemonic() string {
   for true {
     mnemonic := getMnemonicFromUser()
     _, err := bip39.EntropyFromMnemonic(mnemonic)
     if err != nil {
       fmt.Println("Error:",err)
     }
     if err == nil {
       return mnemonic
     }
   }
   return ""
 }

 //Takes ed25519 keypair and writes it to a file called <bs58_address>.json. Format is compatible with Solana CLI.
 func writeKeyToFile(key ed25519.PrivateKey) {
   intArr := make([]int, len(key))
   for i,j := range key {
     intArr[i] = int(j)
   }
   var strKey string
   strKey = "["
   for i, n := range intArr {
    str:= strconv.Itoa(n)
    strKey = strKey + str
    if(i < len(intArr) - 1) {
      strKey = strKey +  ","
    }
   }
   strKey += "]"
   address := getAddress(key)
   err := ioutil.WriteFile(address + ".json", []byte(strKey), 0644)
   if(err != nil) {
     fmt.Println("Error writing to file:", err)
   } else {
     fmt.Println("Imported Account", address)
   }
 }
