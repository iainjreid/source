// Copyright 2026 Iain J. Reid
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ssh

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/anmitsu/go-shlex"
	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"golang.org/x/crypto/ssh"
)

type IdentityLoader struct {
	storer storer.Storer
}

func NewIdentityLoader(storer storer.Storer) Loader {
	return &IdentityLoader{
		storer: storer,
	}
}

func (i *IdentityLoader) Load(ep *transport.Endpoint) (storer.Storer, error) {
	return i.storer, nil
}

type LoggedReadWriter struct {
	internal io.ReadWriter
}

func (l LoggedReadWriter) Read(data []byte) (int, error) {
	i, e := l.internal.Read(data)
	// println("read  ", string(data[:]))
	return i, e
}

func (l LoggedReadWriter) Write(data []byte) (int, error) {
	i, e := l.internal.Write(data)
	// println("write ", string(data[:]))
	return i, e
}

var config = &ssh.ServerConfig{
	NoClientAuth: true,
}

func Init(pemString string) error {
	signer, err := ssh.ParsePrivateKey([]byte(pemString))

	if err != nil {
		return err
	}

	config.AddHostKey(signer)

	return nil
}

func NewServer(storage storer.Storer, port int) error {
	loader := NewIdentityLoader(storage)
	svr := NewSSHServer(loader)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	defer lis.Close()

	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}

		go func(conn net.Conn) {
			defer conn.Close()

			sshConn, chanc, reqc, err := ssh.NewServerConn(conn, config)
			if err != nil {
				log.Println(err)
				return
			}
			defer sshConn.Close()
			go ssh.DiscardRequests(reqc)
			for chanr := range chanc {
				switch chanr.ChannelType() {
				case "session":
					ch, reqc, err := chanr.Accept()
					if err != nil {
						log.Println(err)
						return
					}
					handleSSHSession(svr, ch, reqc)
				default:
					log.Printf("unhandled channel: %s", chanr.ChannelType())
				}
			}
		}(conn)
	}

	return nil
}

func handleSSHSession(svr transport.Transport, ch ssh.Channel, reqc <-chan *ssh.Request) {
	defer ch.Close()

	var exitCode uint32
	defer func() {
		b := ssh.Marshal(struct{ Value uint32 }{exitCode})
		ch.SendRequest("exit-status", false, b)
	}()

	envs := make(map[string]string)
	for req := range reqc {
		switch req.Type {
		case "env":
			payload := struct{ Key, Value string }{}
			ssh.Unmarshal(req.Payload, &payload)
			envs[payload.Key] = payload.Value
			req.Reply(true, nil)
		case "exec":
			payload := struct{ Value string }{}
			ssh.Unmarshal(req.Payload, &payload)
			args, err := shlex.Split(payload.Value, true)
			if err != nil {
				log.Println("lex args", err)
				exitCode = 1
				return
			}
			log.Printf("args: #%v", args)

			cmd := args[0]
			switch cmd {
			case "git-upload-pack": // read
				if gp := envs["GIT_PROTOCOL"]; gp != "version=2" {
					log.Println("unhandled GIT_PROTOCOL", gp)
					exitCode = 1
					return
				}
				err = handleUploadPack(svr, ch)
				if err != nil {
					log.Println(err)
					exitCode = 1
					return
				}

				req.Reply(true, nil)
				return
			case "git-receive-pack": // write
				err = handleReceivePack(svr, ch)
				if err != nil {
					log.Println(err)
					exitCode = 1
					return
				}

				req.Reply(true, nil)
				return
			default:
				log.Printf("unhandled cmd: %s", cmd)
				req.Reply(false, nil)
				exitCode = 1
				return
			}
		case "auth-agent-req@openssh.com":
			if req.WantReply {
				req.Reply(true, nil)
			}
		default:
			log.Printf("unhandled req type: %s", req.Type)
			req.Reply(false, nil)
			exitCode = 1
			return
		}
	}
}

func handleReceivePack(svr transport.Transport, ch ssh.Channel) error {
	ctx := context.Background()

	chwrap := LoggedReadWriter{internal: ch}

	ep, err := transport.NewEndpoint("/")
	if err != nil {
		return fmt.Errorf("create transport endpoint: %w", err)
	}

	recievePackSession, err := svr.NewReceivePackSession(ep, nil)
	if err != nil {
		return fmt.Errorf("create receive-pack session: %w", err)
	}

	advertisedReferences, err := recievePackSession.AdvertisedReferencesContext(ctx)
	if err != nil {
		return fmt.Errorf("get advertised references: %w", err)
	}

	err = advertisedReferences.Encode(chwrap)
	if err != nil {
		return fmt.Errorf("encode advertised references: %w", err)
	}

	referenceUpdateRequest := packp.NewReferenceUpdateRequest()
	err = referenceUpdateRequest.Decode(chwrap)
	if err != nil {
		if err == packp.ErrFlushPacketRecieved {
			return pktline.NewEncoder(chwrap).Flush()
		} else {
			return fmt.Errorf("decode reference-update request: %w", err)
		}
	}

	res, err := recievePackSession.ReceivePack(ctx, referenceUpdateRequest)
	if err != nil {
		return fmt.Errorf("create receive-pack response: %w", err)
	}
	err = res.Encode(chwrap)
	if err != nil {
		return fmt.Errorf("encode receive-pack response: %w", err)
	}

	return nil
}

func handleUploadPack(svr transport.Transport, ch ssh.Channel) error {
	ctx := context.Background()

	chwrap := LoggedReadWriter{internal: ch}

	ep, err := transport.NewEndpoint("/")
	if err != nil {
		return fmt.Errorf("create transport endpoint: %w", err)
	}
	sess, err := svr.NewUploadPackSession(ep, nil)
	if err != nil {
		return fmt.Errorf("create upload-pack session: %w", err)
	}

	ar, err := sess.AdvertisedReferencesContext(ctx)
	if err != nil {
		return fmt.Errorf("get advertised references: %w", err)
	}

	err = ar.Encode(chwrap)
	if err != nil {
		return fmt.Errorf("encode advertised references: %w", err)
	}

	upr := packp.NewUploadPackRequest()
	err = upr.Decode(chwrap)
	if err != nil {
		return fmt.Errorf("decode upload-pack request: %w", err)
	}

	res, err := sess.UploadPack(ctx, upr)
	if err != nil {
		return fmt.Errorf("create upload-pack response: %w", err)
	}
	err = res.Encode(ch)
	if err != nil {
		return fmt.Errorf("encode upload-pack response: %w", err)
	}

	return nil
}
