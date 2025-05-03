#!/usr/bin/env python3
"""
Script to generate a self-signed SSL certificate for development purposes.
For production, you should use a proper certificate from a trusted CA.
"""
import os
from OpenSSL import crypto

def generate_self_signed_cert(cert_file="ssl/cert.pem", key_file="ssl/key.pem"):
    """Generate a self-signed certificate and key pair."""
    # Create directory if it doesn't exist
    os.makedirs(os.path.dirname(cert_file), exist_ok=True)
    
    # Create a key pair
    k = crypto.PKey()
    k.generate_key(crypto.TYPE_RSA, 2048)
    
    # Create a self-signed cert
    cert = crypto.X509()
    cert.get_subject().C = "US"
    cert.get_subject().ST = "California"
    cert.get_subject().L = "Berkeley"
    cert.get_subject().O = "Berkeley Technology and Justice Lab"
    cert.get_subject().OU = "Motion Index"
    cert.get_subject().CN = "localhost"
    cert.set_serial_number(1000)
    cert.gmtime_adj_notBefore(0)
    cert.gmtime_adj_notAfter(10*365*24*60*60)  # 10 years
    cert.set_issuer(cert.get_subject())
    cert.set_pubkey(k)
    cert.sign(k, 'sha256')
    
    # Write certificate and key to files
    with open(cert_file, "wb") as f:
        f.write(crypto.dump_certificate(crypto.FILETYPE_PEM, cert))
    
    with open(key_file, "wb") as f:
        f.write(crypto.dump_privatekey(crypto.FILETYPE_PEM, k))
    
    print(f"Certificate generated: {cert_file}")
    print(f"Private key generated: {key_file}")

if __name__ == "__main__":
    generate_self_signed_cert()
