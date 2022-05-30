package dns

import (
	"encoding/binary"
	"io"
	"net"
	"net/netip"
	"netsim/internal/config"
	"netsim/internal/tunnel"
	"netsim/pkg/pool"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type AnswerHandler func(domain string, ip netip.Addr) error

type Config struct {
	DnsServers []string
	Timeout    int
	MaxTTL     int
	MinTTL     int
	CacheSize  int
}

type Client struct {
	forwarder   *tunnel.Forwarder
	cache       *LruCache
	config      *Config
	upStream    *UPStream
	upStreamMap map[string]*UPStream
	handlers    []AnswerHandler
}

func NewClient(fc *config.Forwarder, cc *Config) (*Client, error) {
	forwarder, err := tunnel.NewForwarder(
		fc.Mode,
		time.Duration(fc.DTimeout)*time.Second,
		time.Duration(fc.RTimeout)*time.Second,
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		forwarder:   forwarder,
		cache:       NewLruCache(cc.CacheSize),
		config:      cc,
		upStream:    NewUPStream(cc.DnsServers),
		upStreamMap: make(map[string]*UPStream),
	}, nil
}

func (c *Client) Exchange(reqBytes []byte, clientAddr string, preferTCP bool) ([]byte, error) {
	req, err := UnmarshalMessage(reqBytes)
	if err != nil {
		return nil, err
	}

	if req.Question.QTYPE == QTypeAAAA {
		respBytes := valCopy(reqBytes)
		respBytes[2] |= uint8(ResponseMsg) << 7
		return respBytes, nil
	}

	if req.Question.QTYPE == QTypeA || req.Question.QTYPE == QTypeAAAA {
		if v, expired := c.cache.Get(qKey(req.Question)); len(v) > 2 {
			v = valCopy(v)
			binary.BigEndian.PutUint16(v[:2], req.ID)

			if expired { // update cache
				go func(qname string, reqBytes []byte, preferTCP bool) {
					defer pool.PutBuffer(reqBytes)
					if dnsServer, network, dialerAddr, respBytes, err := c.exchange(qname, reqBytes, preferTCP); err == nil {
						c.handleAnswer(respBytes, "cache", dnsServer, network, dialerAddr)
					}
				}(req.Question.QNAME, valCopy(reqBytes), preferTCP)
			}
			return v, nil
		}
	}

	dnsServer, network, dialerAddr, respBytes, err := c.exchange(req.Question.QNAME, reqBytes, preferTCP)
	if err != nil {
		return nil, err
	}

	if req.Question.QTYPE != QTypeA && req.Question.QTYPE != QTypeAAAA {
		logrus.Infof("[dns] %s <-> %s(%s) via %s, type: %d, %s",
			clientAddr, dnsServer, network, dialerAddr, req.Question.QTYPE, req.Question.QNAME)
		return respBytes, nil
	}

	err = c.handleAnswer(respBytes, clientAddr, dnsServer, network, dialerAddr)
	return respBytes, err
}

func (c *Client) handleAnswer(respBytes []byte, clientAddr, dnsServer, network, dialerAddr string) error {
	resp, err := UnmarshalMessage(respBytes)
	if err != nil {
		return err
	}

	ips, ttl := c.extractAnswer(resp)
	if ttl > c.config.MaxTTL {
		ttl = c.config.MaxTTL
	} else if ttl < c.config.MinTTL {
		ttl = c.config.MinTTL
	}

	if ttl <= 0 { // we got a null result
		ttl = 1800
	}

	c.cache.Set(qKey(resp.Question), valCopy(respBytes), ttl)
	logrus.Infof("[dns] %s <-> %s(%s) via %s, %s/%d: %s, ttl: %ds",
		clientAddr, dnsServer, network, dialerAddr, resp.Question.QNAME, resp.Question.QTYPE, strings.Join(ips, ","), ttl)

	return nil
}

func (c *Client) extractAnswer(resp *Message) ([]string, int) {
	var ips []string
	ttl := c.config.MinTTL
	for _, answer := range resp.Answers {
		if answer.TYPE == QTypeA || answer.TYPE == QTypeAAAA {
			if answer.IP.IsValid() && !answer.IP.IsUnspecified() {
				for _, h := range c.handlers {
					h(resp.Question.QNAME, answer.IP)
				}
				ips = append(ips, answer.IP.String())
			}
			if answer.TTL != 0 {
				ttl = int(answer.TTL)
			}
		}
	}

	return ips, ttl
}

func (c *Client) exchange(qname string, reqBytes []byte, preferTCP bool) (
	server, network, dialerAddr string, respBytes []byte, err error) {

	// use tcp to connect upstream server default
	network = "tcp"
	dialer := c.forwarder.Handler(qname + ":0")

	// if we are resolving a domain which uses a forwarder `REJECT`, then use `DIRECT` instead
	// so we can resolve it correctly.
	// TODO: dialer.Addr() == "REJECT", tricky
	if dialer.Addr() == "REJECT" {
		dialer = c.proxy.NextDialer("direct:0")
	}

	// If client uses udp and no forwarders specified, use udp
	// TODO: dialer.Addr() == "DIRECT", tricky
	if !preferTCP && dialer.Addr() == "DIRECT" {
		network = "udp"
	}

	ups := c.UpStream(qname)
	server = ups.Server()
	for i := 0; i < ups.Len(); i++ {
		var rc net.Conn
		rc, err = dialer.Dial(network, server)
		if err != nil {
			newServer := ups.SwitchIf(server)
			logrus.Infof("[dns] error in resolving %s, failed to connect to server %v via %s: %v, next server: %s",
				qname, server, dialer.Addr(), err, newServer)
			server = newServer
			continue
		}
		defer rc.Close()

		// TODO: support timeout setting for different upstream server
		if c.config.Timeout > 0 {
			rc.SetDeadline(time.Now().Add(time.Duration(c.config.Timeout) * time.Second))
		}

		switch network {
		case "tcp":
			respBytes, err = c.exchangeTCP(rc, reqBytes)
		case "udp":
			respBytes, err = c.exchangeUDP(rc, reqBytes)
		}

		if err == nil {
			break
		}

		newServer := ups.SwitchIf(server)
		logrus.Infof("[dns] error in resolving %s, failed to exchange with server %v via %s: %v, next server: %s",
			qname, server, dialer.Addr(), err, newServer)

		server = newServer
	}

	// if all dns upstreams failed, then maybe the forwarder is not available.
	if err != nil {
		c.proxy.Record(dialer, false)
	}

	return server, network, dialer.Addr(), respBytes, err
}

func (c *Client) exchangeTCP(rc net.Conn, reqBytes []byte) ([]byte, error) {
	lenBuf := pool.GetBuffer(2)
	defer pool.PutBuffer(lenBuf)

	binary.BigEndian.PutUint16(lenBuf, uint16(len(reqBytes)))
	if _, err := (&net.Buffers{lenBuf, reqBytes}).WriteTo(rc); err != nil {
		return nil, err
	}

	var respLen uint16
	if err := binary.Read(rc, binary.BigEndian, &respLen); err != nil {
		return nil, err
	}

	respBytes := pool.GetBuffer(int(respLen))
	_, err := io.ReadFull(rc, respBytes)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}

func (c *Client) exchangeUDP(rc net.Conn, reqBytes []byte) ([]byte, error) {
	if _, err := rc.Write(reqBytes); err != nil {
		return nil, err
	}

	respBytes := pool.GetBuffer(UDPMaxLen)
	n, err := rc.Read(respBytes)
	if err != nil {
		return nil, err
	}

	return respBytes[:n], nil
}

func (c *Client) SetServers(domain string, servers []string) {
	c.upStreamMap[strings.ToLower(domain)] = NewUPStream(servers)
}

func (c *Client) UpStream(domain string) *UPStream {
	domain = strings.ToLower(domain)
	for i := len(domain); i != -1; {
		i = strings.LastIndexByte(domain[:i], '.')
		if upstream, ok := c.upStreamMap[domain[i+1:]]; ok {
			return upstream
		}
	}
	return c.upStream
}

func (c *Client) AddHandler(h AnswerHandler) {
	c.handlers = append(c.handlers, h)
}

func MakeResponse(domain, ip string, ttl uint32) (*Message, error) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil, err
	}

	var qtype, rdlen uint16 = QTypeA, net.IPv4len
	if addr.Is6() {
		qtype, rdlen = QTypeAAAA, net.IPv6len
	}

	m := NewMessage(0, ResponseMsg)
	m.SetQuestion(NewQuestion(qtype, domain))
	rr := &RR{NAME: domain, TYPE: qtype, CLASS: ClassINET,
		TTL: ttl, RDLENGTH: rdlen, RDATA: addr.AsSlice()}
	m.AddAnswer(rr)

	return m, nil
}

func qKey(q *Question) string {
	return q.QNAME + "/" + strconv.FormatUint(uint64(q.QTYPE), 10)
}

func valCopy(v []byte) (b []byte) {
	if v != nil {
		b = pool.GetBuffer(len(v))
		copy(b, v)
	}
	return
}
