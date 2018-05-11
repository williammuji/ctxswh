package main

import (
	"io"
	"sync"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/williammuji/ctxswh/chatpb"
	"github.com/williammuji/ctxswh/gatewaypb"

	"github.com/garyburd/redigo/redis"
)

type topic struct {
	name string
	pool *redis.Pool
	sync.RWMutex
	sessions map[*session]struct{}
}

func (t *topic) add(sess *session) {
	t.Lock()
	defer t.Unlock()

	msg := t.message(t.getContents())
	sess.writeGameChan(msg)
	t.sessions[sess] = struct{}{}
}

func (t *topic) remove(sess *session) {
	t.Lock()
	defer t.Unlock()

	delete(t.sessions, sess)
}

func (t *topic) getContents() []*gatewaypb.PubResponse_Data {
	c := t.pool.Get()
	defer c.Close()

	values, err := redis.Values(c.Do("ZRANGE", t.name, -10, -1, "WITHSCORES")) //FIXME 10
	if err != nil {
		glog.Error("topic.getContents ZRANGE", t, t.pool, c, values, t.name, err)
		return nil
	}

	length := len(values)
	if length%2 != 0 {
		glog.Error("topic.getContents length%2!=0", t, t.pool, c, values, t.name, length)
		return nil
	}

	pairs := make([]*gatewaypb.PubResponse_Data, length/2)
	for i := range pairs {
		var pair gatewaypb.PubResponse_Data
		values, err = redis.Scan(values, &pair.Content, &pair.PubTime)
		if err != nil {
			glog.Error("topic.getContents redis.Scan error", t, t.pool, c, values, t.name, length, i)
			return nil
		}
		pairs[i] = &pair
	}
	return pairs
}

func (t *topic) pub(data []byte, pubTime time.Time) {
	c := t.pool.Get()
	defer c.Close()

	_, err := c.Do("ZADD", t.name, pubTime.Unix(), data)
	if err != nil {
		glog.Error("topic.pub failed", t, data, pubTime, t.pool, c, err)
		return
	}

	msg := t.message(t.getContents())
	t.Lock()
	for s, _ := range t.sessions {
		s.writeGameChan(msg)
	}
	t.Unlock()
}

func (t *topic) message(ds []*gatewaypb.PubResponse_Data) *chatpb.ChatResponse {
	return &chatpb.ChatResponse{
		MessageResponse: &chatpb.ChatResponse_PubResponse{
			PubResponse: &gatewaypb.PubResponse{
				Topic: t.name,
				Datas: ds,
			},
		},
	}
}

type session struct {
	cs          *chatServer
	gGRPCStream chatpb.ChatService_ServeServer
	wg          sync.WaitGroup
	respC       chan *chatpb.ChatResponse
	topics      map[string]struct{}
}

type chatServer struct {
	pool *redis.Pool
	sync.RWMutex
	topics map[string]*topic
}

func NewChatServer(pool *redis.Pool) (*chatServer, error) {
	return &chatServer{
		pool:   pool,
		topics: make(map[string]*topic),
	}, nil
}

func (cs *chatServer) getTopic(name string) *topic {
	cs.Lock()
	if t, ok := cs.topics[name]; ok {
		cs.Unlock()
		return t
	}

	t := &topic{
		name:     name,
		pool:     cs.pool,
		sessions: make(map[*session]struct{}),
	}
	cs.topics[name] = t
	cs.Unlock()
	return t
}

func (cs *chatServer) pub(name string, content []byte) {
	t := cs.getTopic(name)
	t.pub(content, time.Now())
}

func (cs *chatServer) sub(sess *session, name string) {
	sess.topics[name] = struct{}{}
	t := cs.getTopic(name)
	t.add(sess)
}

func (cs *chatServer) unsub(sess *session, name string) {
	delete(sess.topics, name)
	t := cs.getTopic(name)
	t.remove(sess)
}

func (cs *chatServer) Serve(stream chatpb.ChatService_ServeServer) error {
	sess := &session{
		cs:          cs,
		gGRPCStream: stream,
		respC:       make(chan *chatpb.ChatResponse),
		topics:      make(map[string]struct{}),
	}

	sess.wg.Add(1)
	go func() {
		sess.wg.Wait()
		close(sess.respC)
		glog.Infof("chatServer.Serve sess.wg.Wait chatServer:%#v sess:%#v", cs, sess)
	}()

	go sess.recvLoop()
	return sess.sendLoop()
}

func (sess *session) recvLoop() error {
	defer sess.wg.Done()

	for {
		in, err := sess.gGRPCStream.Recv()
		if err == io.EOF {
			glog.Infof("session.recvLoop sess.gGRPCStream.Recv io.EOF sess:%#v", sess)
			return nil
		}
		if err != nil {
			glog.Infof("session.recvLoop sess.gGRPCStream.Recv error sess:%#v err:%v", sess, err)
			return err
		}

		switch req := in.MessageRequest.(type) {
		case *chatpb.ChatRequest_SubRequest:
			sess.cs.sub(sess, req.SubRequest.GetTopic())
			glog.Infof("session.recvLoop recv chatpb.ChatRequest_SubRequest sess:%#v req:%#v in:%#v sess.gGRPCStream:%#v", sess, req, in, sess.gGRPCStream)
		case *chatpb.ChatRequest_UnsubRequest:
			sess.cs.unsub(sess, req.UnsubRequest.GetTopic())
			glog.Infof("session.recvLoop recv chatpb.ChatRequest_UnsubRequest sess:%#v req:%#v in:%#v sess.gGRPCStream:%#v", sess, req, in, sess.gGRPCStream)
		case *chatpb.ChatRequest_PubRequest:
			sess.cs.pub(req.PubRequest.GetTopic(), req.PubRequest.GetContent())
			glog.Infof("session.recvLoop recv chatpb.ChatRequest_PubRequest sess:%#v req:%#v in:%#v sess.gGRPCStream:%#v", sess, req, in, sess.gGRPCStream)
		default:
			glog.Infof("session.recvLoop recv invalid MessageRequest sess:%#v req:%#v in:%#v sess.gGRPCStream:%#v", sess, req, in, sess.gGRPCStream)
			return grpc.Errorf(codes.InvalidArgument, "invalid MessageRequest: %v", in.MessageRequest)
		}
	}
}

func (sess *session) sendLoop() error {
	for resp := range sess.respC {
		if err := sess.gGRPCStream.Send(resp); err != nil {
			glog.Infof("session.sendLoop sess.gGRPCStream.Send error sess:%#v resp:%#v err:%v", sess, resp, err)
			return err
		}
	}
	glog.Infof("session.sendLoop exit sess:%#v", sess)
	return nil
}

func (sess *session) writeGameChan(resp *chatpb.ChatResponse) error {
	select {
	case sess.respC <- resp:
	case <-sess.gGRPCStream.Context().Done():
		glog.Infof("session.writeGameChan sess.gGRPCStream.Context().Done() sess:%#v resp:%#v", sess, resp)
		return sess.gGRPCStream.Context().Err()
	}
	return nil
}
