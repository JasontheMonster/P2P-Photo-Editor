package photoEditor

import (
	"math"
	"log"
	"time"
)

type Tag struct {
	id		string
	ver_num	int
}

func (this *Tag) smallerTag(other Tag) bool {
	if this.ver_num < other.ver_num {
		return true
	}
	else if this.ver_num > other.ver_num {
		return false
	}
	else {
		return this.id < other.id
	}
}

type TagVal struct {
	Tag_var	Tag
	val 	[]byte
}