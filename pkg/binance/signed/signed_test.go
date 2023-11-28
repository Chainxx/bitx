package signed

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSign(t *testing.T) {
	key := "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559"
	secret := ""
	
	s, err := Sign(secret, key)
	if err != nil {
		assert.Nil(t, err)
	}
	
	fmt.Println("sign: ", s)
}
