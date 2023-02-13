// Package deployer for grid deployer
package deployer

import (
	"fmt"
	"log"
	"math/big"
	"net"
	"regexp"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/subi"
	proxy "github.com/threefoldtech/grid_proxy_server/pkg/client"
	"github.com/threefoldtech/substrate-client"
)

// validateAccount checks the mnemonics is associated with an account with key type ed25519
func validateAccount(sub subi.SubstrateExt, identity substrate.Identity, mnemonics string) error {
	_, err := sub.GetAccount(identity)
	if err != nil && !errors.Is(err, substrate.ErrAccountNotFound) {
		return errors.Wrap(err, "failed to get account with the given mnemonics")
	}

	if err != nil { // Account not found
		funcs := map[string]func(string) (substrate.Identity, error){"ed25519": substrate.NewIdentityFromEd25519Phrase, "sr25519": substrate.NewIdentityFromSr25519Phrase}
		for keyType, f := range funcs {
			ident, err2 := f(mnemonics)
			if err2 != nil { // shouldn't happen, return original error
				log.Printf("couldn't convert the mnemonics to %s key: %s", keyType, err2.Error())
				return err
			}
			_, err2 = sub.GetAccount(ident)
			if err2 == nil { // found an identity with key type other than the provided
				return fmt.Errorf("found an account with %s key type and the same mnemonics, make sure you provided the correct key type", keyType)
			}
		}
		// didn't find an account with any key type
		return err
	}
	return nil
}

func validateTwinYggdrasil(sub subi.SubstrateExt, twinID uint32) error {
	yggIP, err := sub.GetTwinIP(twinID)
	if err != nil {
		return errors.Wrapf(err, "could not get twin %d from substrate", twinID)
	}

	ip := net.ParseIP(yggIP)
	listenIP := yggIP
	if ip != nil && ip.To4() == nil {
		// if it's ipv6 surround it with brackets
		// otherwise, keep as is (can be ipv4 or a domain (probably will fail later but we don't care))
		listenIP = fmt.Sprintf("[%s]", listenIP)
	}

	s, err := net.Listen("tcp", fmt.Sprintf("%s:0", listenIP))
	if err != nil {
		return errors.Wrapf(err, "couldn't listen on port. make sure the twin id is associated with a valid yggdrasil ip, twin id: %d, ygg ip: %s, err", twinID, yggIP)
	}
	defer s.Close()

	port := s.Addr().(*net.TCPAddr).Port
	arrived := false
	go func() {
		c, err := s.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			return
		}
		arrived = true
		c.Close()
	}()

	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", listenIP, port))
	if err != nil {
		return errors.Wrapf(err, "failed to connect to ip. make sure the twin id is associated with a valid yggdrasil ip, twin id: %d, ygg ip: %s, err", twinID, yggIP)
	}
	c.Close()

	if !arrived {
		return errors.Wrapf(err, "sent request but didn't arrive to me. make sure the twin id is associated with a valid yggdrasil ip, twin id: %d, ygg ip: %s, err", twinID, yggIP)
	}

	return nil
}

func validateRMBProxyServer(gridProxyClient proxy.Client) error {
	return gridProxyClient.Ping()
}

func validateMnemonics(mnemonics string) error {
	if len(mnemonics) == 0 {
		return errors.New("mnemonics required")
	}

	alphaOnly := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !alphaOnly.MatchString(mnemonics) {
		return errors.New("mnemonics can only be composed of a non-alphanumeric character or a whitespace")
	}

	return nil
}

func validateSubstrateURL(url string) error {
	if len(url) == 0 {
		return errors.New("substrate url is required")
	}

	alphaOnly := regexp.MustCompile(`^wss:\/\/[a-z0-9]+\.[a-z0-9]\/?([^\s<>\#%"\,\{\}\\|\\\^\[\]]+)?$`)
	if !alphaOnly.MatchString(url) {
		return errors.New("substrate url is not valid")
	}

	return nil
}

func validateProxyURL(url string) error {
	if len(url) == 0 {
		return errors.New("proxy url is required")
	}

	alphaOnly := regexp.MustCompile(`^https:\/\/[a-z0-9]+\.[a-z0-9]\/?([^\s<>\#%"\,\{\}\\|\\\^\[\]]+)?$`)
	if !alphaOnly.MatchString(url) {
		return errors.New("proxy url is not valid")
	}

	return nil
}

func validateAccountBalanceForExtrinsics(sub subi.SubstrateExt, identity substrate.Identity) error {
	balance, err := sub.GetBalance(identity)
	if err != nil && !errors.Is(err, substrate.ErrAccountNotFound) {
		return errors.Wrap(err, "failed to get account with the given mnemonics")
	}

	log.Printf("balance %d\n", balance.Free)
	if balance.Free.Cmp(big.NewInt(20000)) == -1 {
		return fmt.Errorf("account contains %s, min fee is 20000", balance.Free)
	}

	return nil
}
