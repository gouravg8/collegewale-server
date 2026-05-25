package common

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sqids/sqids-go"
)

type MaskedId string

func Mask(id uint, at time.Time) MaskedId {
	if id == 0 {
		return ""
	}
	if scramble {
		sid, _ := s.Encode([]uint64{uint64(id), uint64(at.Unix())})
		return MaskedId(sid)
	}
	return MaskedId(fmt.Sprintf("%d", id))
}

func Unmask(sid MaskedId) uint {

	if scramble {
		//this is to avoid panic and crash
		if len(string(sid)) < 2 {
			return 0
		}
		ids := s.Decode(string(sid))
		if len(ids) == 0 {
			return 0
		}
		return uint(ids[0])
	}
	id, _ := strconv.Atoi(string(sid))
	return uint(id)
}

var s, _ = sqids.New(sqids.Options{
	MinLength: 10,
})
