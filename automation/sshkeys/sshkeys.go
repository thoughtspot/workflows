package sshkeys

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"

	"automation/logger"

	"golang.org/x/crypto/ssh"
)

type SSHKeys struct {
	PublicKey  []byte
	PrivateKey []byte
}

func GenerateED25519Keys() SSHKeys {
	l := logger.New()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		l.Fatal(err)
	}

	privateKeyBlock, err := ssh.MarshalPrivateKey(privateKey, "")
	if err != nil {
		l.Fatal(err)
	}
	opensshPrivateKey := pem.EncodeToMemory(privateKeyBlock)

	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		l.Fatal(err)
	}

	opensshPublicKey := ssh.MarshalAuthorizedKey(sshPublicKey)

	return SSHKeys{
		PublicKey:  opensshPublicKey,
		PrivateKey: opensshPrivateKey,
	}
}
