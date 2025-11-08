package driver

import (
	"fmt"
	"strings"

	"github.com/jimtang2/bc-go/pkg/pb/v1"
	"google.golang.org/protobuf/proto"
)

type Ticker struct {
	*pb.Ticker
}

func (t *Ticker) Fmt() {
	t.Ticker.Pair = strings.ReplaceAll(t.Ticker.Pair, "/", "")
	t.Ticker.Pair = strings.ReplaceAll(t.Ticker.Pair, "-", "")
	t.Ticker.Pair = strings.ReplaceAll(t.Ticker.Pair, ":", "")
	if i := strings.Index(t.Ticker.Pair, "USD"); i > -1 {
		t.Ticker.Pair = t.Ticker.Pair[:i] + ":" + t.Ticker.Pair[i:]
	}
}

// implements kafka.KMessage
func (t *Ticker) IsValid() bool {
	return t != nil && t.Ticker != nil && len(t.Ticker.Pair) > 0 && len(t.Ticker.Exchange) > 0 && t.Ticker.Bid > 0 && t.Ticker.Ask > 0
}

// implements kafka.KMessage
func (t *Ticker) Key() string {
	return fmt.Sprintf("%v:%v", t.Ticker.Exchange, t.Ticker.Pair)
}

// implements kafka.KMessage
func (t *Ticker) Bytes() []byte {
	b, _ := proto.Marshal(t.Ticker)
	return b
}
