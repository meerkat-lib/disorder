// Code generated by https://github.com/meerkat-io/disorder; DO NOT EDIT.
package sub

import ()

type Number struct {
	Value int32 `disorder:"value" json:"value"`
}

type NumberWrapper struct {
	Value *Number `disorder:"value" json:"value,omitempty"`
}
