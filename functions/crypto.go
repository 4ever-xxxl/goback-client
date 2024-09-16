package functions

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

func Encrypt(inBuf *bytes.Buffer, outBuf *bytes.Buffer, key []byte) error {
	// 创建 AES 加密块
	hashedKey := kdf(key)
	block, err := aes.NewCipher(hashedKey)
	if err != nil {
		return err
	}

	// 使用 GCM 模式进行加密
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 生成随机的 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// 加密数据并写入输出缓冲区
	ciphertext := gcm.Seal(nonce, nonce, inBuf.Bytes(), nil)
	_, err = outBuf.Write(ciphertext)
	return err
}

func Decrypt(inBuf *bytes.Buffer, outBuf *bytes.Buffer, key []byte) error {
	hashedKey := kdf(key)
	// 创建 AES 加密块
	block, err := aes.NewCipher(hashedKey)
	if err != nil {
		return err
	}

	// 使用 GCM 模式进行解密
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 从密文中提取 nonce
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := inBuf.Next(nonceSize), inBuf.Bytes()

	// 解密数据并写入输出缓冲区
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}
	_, err = outBuf.Write(plaintext)
	return err
}

func kdf(key []byte) []byte {
	hash := sha256.Sum256(key)
	return hash[:]
}
