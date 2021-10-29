package encrypt

import (
   "testing"
)

const emailPass = "AFcb4KQetybcphDo9XGwoqkEW2iZfYGNrufsM9ABBv9Us461Iyc6Xb9sKAw7dHgBrsOCBV_K7Cv0LqxO2r6E_LBrpSctl8JfWk526oR2F-ZC0lK6S-J9tR31Geilt-yoxgkwNuiLncFKC-O9K_me3K5qo0DxzsWHpdraiU8jV8GvolmizA=="

func TestEncrypt(t *testing.T) {
   sig, err := signature("email", "password")
   if err != nil {
      t.Fatal(err)
   }
   if sig != emailPass {
      t.Fatal(sig)
   }
}
