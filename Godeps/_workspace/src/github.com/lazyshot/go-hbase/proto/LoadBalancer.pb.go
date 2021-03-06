// Code generated by protoc-gen-go.
// source: LoadBalancer.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/ghawk1ns/golf/Godeps/_workspace/src/github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = math.Inf

type LoadBalancerState struct {
	BalancerOn       *bool  `protobuf:"varint,1,opt,name=balancer_on" json:"balancer_on,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *LoadBalancerState) Reset()         { *m = LoadBalancerState{} }
func (m *LoadBalancerState) String() string { return proto1.CompactTextString(m) }
func (*LoadBalancerState) ProtoMessage()    {}

func (m *LoadBalancerState) GetBalancerOn() bool {
	if m != nil && m.BalancerOn != nil {
		return *m.BalancerOn
	}
	return false
}

func init() {
}
