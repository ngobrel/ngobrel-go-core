package crypto
import (
  "crypto/rand"
  "encoding/hex"
  "fmt"
  "github.com/ridon/ngobrel/crypto/Key"
  "io"
  "testing"
)


/**
 This test X3DH protocol where Bob sends a message to Alice
*/
func TestProto(t *testing.T) {
  random := rand.Reader

  // 1. Bob uploads his public key bundles
  bundleBob, _ := Key.NewBundle(random)
  bundleBob.PopulatePreKeys(random, 100)

  bundleBobPublic := bundleBob.Public // Alice only can access bob's public keys 

  // 2. Alice verifies the SPK
  bundleAlice, _ := Key.NewBundle(random)
  res := bundleBobPublic.Verify();


  if res == false {
    t.Error("SPK is not verified")
  }

  // 3. Alice creates an ephemeral key
  ephKey, err := Key.Generate(random)
  if err != nil {
   t.Error("Ephemeral key is not created") 
  }

  // 4. Alice creates a shared key
  // 5. Alice clears the ephemeral key and keys' content

  sk, preKeyId, err := GetSharedKeySender(random, ephKey, bundleAlice, &bundleBobPublic, "Ridon") 
  if err != nil {
    t.Error(err)
  }

  // 6. Alice creates the associated data
  ad := append(bundleAlice.Public.Identity.Encode()[:], bundleBobPublic.Identity.Encode()[:]...)

  // 7. Alice creates the first message
  var nonce[12]byte
	io.ReadFull(random, nonce[:])

  msgToBeEncrypted := []byte("olala")
  message, err := NewMessage(&bundleAlice.Public.Identity, &ephKey.PublicKey, *preKeyId, nonce[:], *sk, msgToBeEncrypted, ad)

  if err != nil {
    t.Error(err)
  }

  // message is then transfered to transit place

  // 1. Bob fetches Alice keys and all private keys

  bundleAlicePublic := bundleAlice.Public

  // 2. Bob gets the shared key
  skBob, err := GetSharedKeyRecipient(message.EphKey, bundleBob, &bundleAlicePublic, message.PreKeyId, "Ridon")

  if err != nil {
    t.Error(err)
  }

  if hex.EncodeToString(*sk) != hex.EncodeToString(*skBob) {
    t.Error("SK is different")
  }

  // 3. Bob creates the associated data
  adBob := append(bundleAlicePublic.Identity.Encode()[:], bundleBobPublic.Identity.Encode()[:]...)

  if hex.EncodeToString(ad) != hex.EncodeToString(adBob) {
    t.Error("SK is different")
  }

  // 4. Bob decrypts the message
  decrypted, err := message.DecryptMessage(*skBob, adBob)
  if err != nil {
    t.Error(err)
  }

  if len(*decrypted) == 0 {
    t.Error("Decrypted data is zero length")
  }

  if hex.EncodeToString(*decrypted) != hex.EncodeToString(msgToBeEncrypted) {
    t.Error("Can't decrypt")
  }

  fmt.Printf("")
}