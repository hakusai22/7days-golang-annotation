package geecache

import pb "geecache/geecachepb"

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
// 修改 peers.go 中的 PeerGetter 接口，参数使用 geecachepb.pb.go 中的数据类型
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
