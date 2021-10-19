from Cryptodome.Cipher import PKCS1_OAEP
from Cryptodome.PublicKey import RSA
import base64
import binascii
import hashlib
import os

# The key is distirbuted with Google Play Services. This one is from version
# 7.3.29.
B64_KEY_7_3_29 = (
    b"AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3"
    b"iJIZdodyhKZQrNWp5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pK"
    b"RI16kB0YppeGx5qIQ5QjKzsR8ETQbKLNWgRY0QRNVz34kMJR3P/LgHax/"
    b"6rmf5AAAAAwEAAQ=="
)

def bytes_to_int(bytes_seq) -> int:
    return int.from_bytes(bytes_seq, "big")

def construct_signature(email, password, key) -> bytes:
   signature = bytearray(b"\x00")
   mod = int_to_bytes(key.n)
   struct = b"\x00\x00\x00\x80" + mod + b"\x00\x00\x00\x03" + int_to_bytes(key.e)
   signature.extend(hashlib.sha1(struct).digest()[:4])
   cipher = PKCS1_OAEP.new(key)
   encrypted_login = cipher.encrypt((email + "\x00" + password).encode("utf-8"))
   signature.extend(encrypted_login)
   return base64.urlsafe_b64encode(signature)

def int_to_bytes(num, pad_multiple = 1) -> bytes:
    """Packs the num into a byte string 0 padded to a multiple of pad_multiple
    bytes in size. 0 means no padding whatsoever, so that packing 0 result
    in an empty string. The resulting byte string is the big-endian two's
    complement representation of the passed in long."""
    if num == 0:
        return b"\0" * pad_multiple
    if num < 0:
        raise ValueError("Can only convert non-negative numbers.")
    value = hex(num)[2:]
    value = value.rstrip("L")
    if len(value) & 1:
        value = "0" + value
    result = binascii.unhexlify(value)
    if pad_multiple not in [0, 1]:
        filled_so_far = len(result) % pad_multiple
        if filled_so_far != 0:
            result = b"\0" * (pad_multiple - filled_so_far) + result
    return result

def key_from_b64(b64_key):
    binKey = base64.b64decode(b64_key)
    i = bytes_to_int(binKey[:4])
    mod = bytes_to_int(binKey[4 : 4 + i])
    j = bytes_to_int(binKey[i + 4 : i + 4 + 4])
    exponent = bytes_to_int(binKey[i + 8 : i + 8 + j])
    key = RSA.construct((mod, exponent))
    return key

encryptedPass = construct_signature(
   'srpen6@gmail.com', os.environ['PASS'], key_from_b64(B64_KEY_7_3_29)
)

print(encryptedPass)
