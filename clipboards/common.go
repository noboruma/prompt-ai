package clipboards

import (
	"errors"

	"github.com/atotto/clipboard"
)

const (
	max_copy = 8
)

type CopyClipboards struct {
	clipboards [max_copy]string
	copy_id    int
}

func NewCopyClipboards() CopyClipboards {
	return CopyClipboards{
		copy_id: 0,
	}
}

func (cc *CopyClipboards) Append(content string) int {
	tmp := cc.copy_id
	cc.clipboards[cc.copy_id] = content
	cc.copy_id += 1
	cc.copy_id %= max_copy
	return tmp
}

func (cc CopyClipboards) Fetch(index int) error {
	if index < 0 || index > max_copy {
		return errors.New("Out of range")
	}
	return clipboard.WriteAll(cc.clipboards[index])
}
