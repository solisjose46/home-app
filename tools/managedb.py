import argparse
import sqlite3
from Crypto.Cipher import Blowfish
from Crypto.Util.Padding import pad
import os

def blowfish_encrypt(password, salt):
    # Ensure the key length is valid for Blowfish (4 to 56 bytes)
    if not (4 <= len(salt) <= 56):
        raise ValueError("Salt must be between 4 and 56 bytes long")

    # Generate a random IV of 8 bytes (Blowfish block size)
    iv = os.urandom(8)
    
    # Create a new Blowfish cipher object
    cipher = Blowfish.new(salt, Blowfish.MODE_CBC, iv)

    # Pad the plaintext password to be a multiple of the block size (8 bytes for Blowfish)
    padded_password = pad(password.encode(), Blowfish.block_size)

    # Encrypt the password
    encrypted_password = iv + cipher.encrypt(padded_password)
    
    return encrypted_password

def insert_user(database_path, username, password, salt):
    encrypted_password = blowfish_encrypt(password, salt)

    # Connect to the SQLite database
    conn = sqlite3.connect(database_path)
    cursor = conn.cursor()

    # Insert the new user
    cursor.execute('''
    INSERT INTO users (username, password) VALUES (?, ?)
    ''', (username, encrypted_password))

    # Commit the changes and close the connection
    conn.commit()
    conn.close()

    print(f"User {username} inserted successfully.")

def main():
    parser = argparse.ArgumentParser(description="Add a new user to the SQLite database with an encrypted password.")
    parser.add_argument('-u', '--username', type=str, required=True, help='Username for the new user')
    parser.add_argument('-p', '--password', type=str, required=True, help='Plain text password')
    parser.add_argument('-s', '--salt', type=str, required=True, help='Salt (must be between 4 and 56 bytes)')
    parser.add_argument('-d', '--database', type=str, required=True, help='Path to the SQLite database')

    args = parser.parse_args()

    username = args.username
    password = args.password
    salt = args.salt.encode()  # Convert salt to bytes
    database_path = args.database

    try:
        insert_user(database_path, username, password, salt)
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    main()
