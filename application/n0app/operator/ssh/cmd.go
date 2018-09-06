package ssh

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
)

type SSH struct {
	addr    string
	config  *ssh.ClientConfig
	conn    *ssh.Client
	session *ssh.Session
}

func ConnectSSH(ip string, port int, user string, key []byte) (*SSH, error) {
	s := &SSH{
		addr: ip + ":" + strconv.Itoa(port),
	}

	k, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key")
	}

	s.config = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(k),
		},
	}

	if s.conn, err = ssh.Dial("tcp", s.addr, s.config); err != nil {
		return nil, fmt.Errorf("Failed to dial ssh, network:tcp, addr:'%s'", s.addr)
	}

	return s, nil
}

func (s *SSH) Close() error {
	defer s.conn.Close()

	return nil
}

func (s SSH) Start(cmd string) (io.Reader, io.Reader, error) {
	// TODO: backoffの実装
	session, err := s.conn.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create new session, err:'%s'", err.Error())
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed open stdout pipe, err:'%s'", err.Error())
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed open stderr pipe, err:'%s'", err.Error())
	}

	if err := session.Start(cmd); err != nil {
		return nil, nil, fmt.Errorf("Failed to start command, err:'%s'", err.Error())
	}

	return stdout, stderr, nil
}

// func (s SSH) Wait()

func (s SSH) RunInteractive(cmd string) (int, error) {
	session, err := s.conn.NewSession()
	if err != nil {
		return 0, fmt.Errorf("Failed to create new session, err:'%s'", err.Error())
	}

	// 依存が強力になってしまうので引数でとったほうがいいかもしれない
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Run(cmd); err != nil {
		if ee, ok := err.(*ssh.ExitError); ok {
			return ee.ExitStatus(), nil
		}

		return 0, err
	}

	return 0, nil
}
