package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chuckpreslar/emission"
	"github.com/sourcegraph/jsonrpc2"
	"github.com/uscott/go-api-deribit/inout"
	syncgrp "github.com/uscott/go-syncgrp"
	"github.com/uscott/go-tools/errs"
	"github.com/uscott/go-tools/tm"
	"github.com/uscott/go-tools/tmath"
	"nhooyr.io/websocket"
)

// Test and Production URLs
const (
	ProdBaseURL    = "wss://www.deribit.com/ws/api/v2/"
	TestBaseURL    = "wss://test.deribit.com/ws/api/v2/"
	exchTmStmpUnit = time.Millisecond
)

// MaxTries is the max number of reconnect attempts
const MaxTries = 10
const prvt string = "private"

var (
	// ErrAuthRequired is an error value corresponding to authorization
	// being required for a request
	ErrAuthRequired    = errors.New("AUTHENTICATION IS REQUIRED")
	matchEngineRequest = []string{
		"buy", "sell", "edit", "cancel", "close_position",
		"verify_block_trade", "execute_block_trade",
	}
)

var (
	ceil  = math.Ceil
	clamp = tmath.Clamp
	imin  = tmath.Imin
	max   = math.Max
	min   = math.Min
	trunc = math.Trunc
)

// Event is wrapper of received event
type Event struct {
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

// Configuration contains data for creating
// a client
type Configuration struct {
	Ctx                context.Context
	Address            string  `json:"addr"`
	AutoReconnect      bool    `json:"autoReconnect"`
	AutoRefillMatch    float64 `json:"auto_refill_match"`
	AutoRefillNonmatch float64 `json:"auto_refill_nonmatch"`
	Currency           string  `json:"currency"`
	DebugMode          bool    `json:"debugMode"`
	Key                string  `json:"api_key"`
	Production         bool    `json:"production"`
	Secret             string  `json:"secret_key"`
	UseLogFile         bool    `json:"use_log_file"`
}

// DfltCnfg returns a default Configuration
func DfltCnfg() *Configuration {
	return &Configuration{
		Address:            TestBaseURL,
		AutoReconnect:      true,
		AutoRefillMatch:    0.8,
		AutoRefillNonmatch: 0.8,
		Currency:           BTC,
		DebugMode:          true,
		Key:                os.Getenv("DERIBIT_TEST_MAIN_KEY"),
		Production:         false,
		Secret:             os.Getenv("DERIBIT_TEST_MAIN_SECRET"),
		UseLogFile:         false,
	}
}

// Client is the base client for connecting to the exchange
type Client struct {
	auth struct {
		token   string
		refresh string
	}
	autoRefill       rqstCntData
	conn             *websocket.Conn
	emitter          *emission.Emitter
	heartCancel      chan struct{}
	isConnected      bool
	rpcConn          *jsonrpc2.Conn
	rqstCnt          rqstCntData
	rqstTmr          rqstTmrData
	subscriptions    []string
	subscriptionsMap map[string]byte
	Acct             inout.AcctSummaryOut
	Config           *Configuration
	Logger           *log.Logger
	SG               *syncgrp.SyncGrp
	StartTime        time.Time
	Sub              *Subordinate
}

func (c *Client) NewMinimal(cfg *Configuration) (err error) {
	if cfg == nil {
		return errs.ErrNilPtr
	}
	if cfg.Ctx == nil {
		cfg.Ctx = context.Background()
	}
	*c = Client{}
	c.Config = cfg
	c.emitter = emission.NewEmitter()
	c.SG = syncgrp.New()
	c.Sub = NewSubordinate()
	c.subscriptionsMap = make(map[string]byte)
	c.StartTime = tm.UTC()
	return nil
}

func NewMinimal(cfg *Configuration) (*Client, error) {
	if cfg == nil {
		return nil, errs.ErrNilPtr
	}
	var c *Client = new(Client)
	if err := c.NewMinimal(cfg); err != nil {
		return c, err
	}
	return c, nil
}

// New returns pointer to new Client
func New(cfg *Configuration) (*Client, error) {
	if cfg == nil {
		return nil, errs.ErrNilPtr
	}
	var (
		c   *Client = new(Client)
		err error
	)
	if err = c.New(cfg); err != nil {
		return c, err
	}
	return c, nil
}

func (c *Client) New(cfg *Configuration) (err error) {
	if err = c.NewMinimal(cfg); err != nil {
		return err
	}
	if err = c.CreateLogger(); err != nil {
		log.Fatalln(err.Error())
		return err
	}
	if err = c.Start(); err != nil {
		c.Logger.Fatalln(err.Error())
		return err
	}
	c.rqstTmr = rqstTmrData{t0: c.StartTime, t1: c.StartTime, dt: 0}
	if err = c.GetAccountSummary(c.Config.Currency, true, &c.Acct); err != nil {
		go c.Logger.Println(err.Error())
		return err
	}
	var ub float64
	ub = clamp(
		float64(c.Acct.Lmts.MatchingEngine)*cfg.AutoRefillMatch,
		0,
		float64(c.Acct.Lmts.MatchingEngine))
	c.autoRefill.mch = int(math.Floor(ub))
	ub = clamp(
		float64(c.Acct.Lmts.NonMatchingEngine)*cfg.AutoRefillNonmatch,
		0,
		float64(c.Acct.Lmts.NonMatchingEngine))
	c.autoRefill.non = int(math.Floor(ub))
	c.resetRqstTmr()
	return nil
}

func (c *Client) Connect() (*websocket.Conn, *http.Response, error) {
	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()
	conn, resp, err := websocket.Dial(ctx, c.Config.Address, &websocket.DialOptions{})
	if err == nil {
		conn.SetReadLimit(32768 * 64)
	}
	return conn, resp, err
}

func (c *Client) CreateLogger() error {
	var (
		dir, logFilePath, testprod string
		err                        error
		logFile                    *os.File
	)
	if c.Config.Production {
		dir = "log-prod/"
		testprod = "prod"
	} else {
		dir = "log-test/"
		testprod = "test"
	}
	_ = os.Mkdir(dir, os.ModeDir)
	_ = os.Chmod(dir, 0754)
	if c.Config.UseLogFile {
		stamp := tm.Format2(tm.UTC())
		s := fmt.Sprintf("%v%v-%v-%v.%v", dir, "api-log", testprod, stamp, "log")
		logFilePath = strings.ReplaceAll(s, " ", "-")
		logFilePath = strings.ReplaceAll(logFilePath, ":", "")
		logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	} else {
		logFile = os.Stdout
	}
	c.Logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile|log.LUTC|log.Lmsgprefix)
	return nil
}

func (c *Client) decrementRqstCnt(nsecs int) {
	if nsecs > 0 {
		c.SG.Lock()
		c.rqstCnt.mch = imax(0, c.rqstCnt.mch-nsecs*c.Acct.Lmts.MatchingEngine)
		c.rqstCnt.non = imax(0, c.rqstCnt.non-nsecs*c.Acct.Lmts.NonMatchingEngine)
		c.SG.Unlock()
	}
}

func (c *Client) heartbeat() {
	t := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-t.C:
			c.Test()
		case <-c.heartCancel:
			return
		}
	}
}

func (c *Client) Reconnect() {
	notify := c.rpcConn.DisconnectNotify()
	<-notify
	c.setIsConnected(false)
	c.Logger.Println("disconnect, reconnect...")
	close(c.heartCancel)
	time.Sleep(4 * time.Second)
	if err := c.Start(); err != nil {
		go c.Logger.Println(err.Error())
	}
}

func (c *Client) resetRqstTmr() {
	c.SG.Lock()
	t0 := c.rqstTmr.t1
	t1 := tm.UTC()
	c.rqstTmr.t0, c.rqstTmr.t1, c.rqstTmr.dt = t0, t1, t1.Sub(t0)
	c.SG.Unlock()
}

// setIsConnected sets state for isConnected
func (c *Client) setIsConnected(state bool) {
	c.SG.RWLock()
	c.isConnected = state
	c.SG.RWUnlock()
}

func (c *Client) Start() (err error) {
	c.setIsConnected(false)
	c.subscriptionsMap = make(map[string]byte)
	c.conn, c.rpcConn = nil, nil
	c.heartCancel = make(chan struct{})
	for i := 0; i < MaxTries; i++ {
		conn, _, err := c.Connect()
		if err != nil {
			c.Logger.Println(err.Error())
			tm := time.Duration(i+1) * 5 * time.Second
			c.Logger.Printf("Sleeping %v\n", tm)
			time.Sleep(tm)
			continue
		}
		c.conn = conn
		break
	}
	if c.conn == nil {
		return errs.ErrNotConnected
	}
	c.rpcConn = jsonrpc2.NewConn(
		context.Background(), NewObjectStream(c.conn), c)
	c.setIsConnected(true)
	// auth
	if c.Config.Key != "" && c.Config.Secret != "" {
		if err = c.Auth(c.Config.Key, c.Config.Secret); err != nil {
			return err
		}
	}
	// subscribe
	if err = c.subscribe(c.subscriptions); err != nil {
		return err
	}
	_, err = c.SetHeartbeat(&inout.Heartbeat{Interval: 30})
	if err != nil {
		return err
	}
	if c.Config.AutoReconnect {
		go c.Reconnect()
	}
	go c.heartbeat()
	return nil
}

func (c *Client) subscribe(channels []string) (e error) {
	var (
		pblcChannels []string
		prvtChannels []string
	)
	c.SG.Lock()
	for _, v := range c.subscriptions {
		if _, ok := c.subscriptionsMap[v]; ok {
			continue
		}
		if strings.HasPrefix(v, "user.") {
			prvtChannels = append(prvtChannels, v)
		} else {
			pblcChannels = append(pblcChannels, v)
		}
	}
	c.SG.Unlock()
	if len(pblcChannels) > 0 {
		_, e = c.SubPblc(pblcChannels)
		if e != nil {
			return e
		}
		c.SG.Lock()
		for _, v := range pblcChannels {
			c.subscriptionsMap[v] = 0
		}
		c.SG.Unlock()
	}
	if len(prvtChannels) > 0 {
		_, e = c.SubPrvt(prvtChannels)
		if e != nil {
			return e
		}
		c.SG.Lock()
		for _, v := range prvtChannels {
			c.subscriptionsMap[v] = 0
		}
		c.SG.Unlock()
	}
	return nil
}

func (c *Client) updtRqstTmr() {
	c.SG.Lock()
	c.rqstTmr.t1 = tm.UTC()
	c.rqstTmr.dt = c.rqstTmr.t1.Sub(c.rqstTmr.t0)
	c.SG.Unlock()
}

// AutoRefillRqsts automatically refills rate limit if request counts
// are above certain threshold
func (c *Client) AutoRefillRqsts() {
	c.RefillRqstsCndtnl(c.autoRefill.mch, c.autoRefill.non)
}

// Call issues JSONRPC v2 calls
func (c *Client) Call(method string, params interface{}, result interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			c.Logger.Println(err.Error())
		}
	}()
	if !c.IsConnected() {
		return errs.ErrNotConnected
	}
	if params == nil {
		params = emptyParams
	}
	if token, ok := params.(privateParams); ok {
		if c.auth.token == "" {
			return ErrAuthRequired
		}
		token.setToken(c.auth.token)
	}
	c.SG.Lock()
	ml, pl, engine := len(method), len(prvt), false
	if ml >= pl && method[:pl] == prvt {
		rmdr := method[pl+1:]
		rl := len(rmdr)
		for _, s := range matchEngineRequest {
			if sl := len(s); rl >= sl && rmdr[:sl] == s {
				engine = true
				break
			}
		}
	}
	if engine {
		c.rqstCnt.mch++
	} else {
		c.rqstCnt.non++
	}
	c.SG.Unlock()
	return c.rpcConn.Call(c.Config.Ctx, method, params, result)
}

// ConvertExchStmp converts an exchange time stamp
// to a client-side time.Time
func (c *Client) ConvertExchStmp(ts int64) time.Time {
	ts *= int64(exchTmStmpUnit) / int64(time.Nanosecond)
	return time.Unix(ts/int64(time.Second), ts%int64(time.Second)).UTC()
}

// ExchangeTime returns the exchange time as
// a client-side time.Time
func (c *Client) ExchangeTime() (time.Time, error) {
	var (
		ms  int64
		err error
	)
	if ms, err = c.GetTime(); err != nil {
		return time.Time{}, err
	}
	return c.ConvertExchStmp(ms), nil
}

// Handle implements jsonrpc2.Handler
func (c *Client) Handle(
	ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {

	if req.Method == "subscription" { // update events
		if req.Params != nil && len(*req.Params) > 0 {
			var event Event
			if err := json.Unmarshal(*req.Params, &event); err != nil {
				go c.Logger.Println(err.Error())
				return
			}
			_, err := c.subscriptionsProcess(&event)
			if err != nil {
				go c.Logger.Println(err.Error())
			}
		}
	}
}

// IsConnected returns the WebSocket connection state
func (c *Client) IsConnected() bool {
	c.SG.RLock()
	defer c.SG.RUnlock()
	return c.isConnected
}

// IsProduction returns whether the client is connected
// to the production server
func (c *Client) IsProduction() bool {
	return c.Config.Production
}

// RefillRqsts will sleep long enough to refill all rate limits
func (c *Client) RefillRqsts() {
	c.SG.Lock()
	m, n := c.Acct.Lmts.MatchingEngine, c.Acct.Lmts.NonMatchingEngine
	mb, nb := c.Acct.Lmts.MatchingEngineBurst, c.Acct.Lmts.NonMatchingEngineBurst
	if m <= 0 || n <= 0 || mb <= 0 || nb <= 0 {
		c.SG.Unlock()
		return
	}
	const (
		fnanosecs  float64       = float64(time.Second) / float64(time.Nanosecond)
		minSleepTm time.Duration = 250 * time.Millisecond
	)
	tmch := float64(c.rqstCnt.mch) / float64(m) * fnanosecs
	tnon := float64(c.rqstCnt.non) / float64(n) * fnanosecs
	c.SG.Unlock()
	c.updtRqstTmr()
	c.SG.Lock()
	ub := imin(mb/m, nb/n) // seconds
	tacm := trunc(min(float64(ub), c.rqstTmr.dt.Seconds())) * fnanosecs
	tnet := time.Duration(max(tmch, tnon) - tacm) // Nanoseconds
	c.SG.Unlock()
	if tnet > minSleepTm {
		time.Sleep(tnet)
		c.resetRqstTmr()
		c.SG.Lock()
		c.rqstCnt.mch, c.rqstCnt.non = 0, 0
		c.SG.Unlock()
	}
}

// RefillRqstsCndtnl refills requests if request count
// are above given amounts
func (c *Client) RefillRqstsCndtnl(match int, nonmatch int) {
	if c.rqstCnt.mch > match || c.rqstCnt.non > nonmatch {
		c.RefillRqsts()
	}
}

// RqstCnts returns the number of requsts accumulated
func (c *Client) RqstCnts() (cntMch, cntNon int) {
	cntMch, cntNon = c.rqstCnt.mch, c.rqstCnt.non
	return
}

// SubscribeToChannels subscribes to channels
func (c *Client) SubscribeToChannels(channels []string) (e error) {
	c.SG.Lock()
	c.subscriptions = append(c.subscriptions, channels...)
	c.SG.Unlock()
	if e = c.subscribe(channels); e != nil {
		return e
	}
	// Remove any dupes in c.subscriptions
	c.SG.Lock()
	l := len(c.subscriptionsMap)
	if cap(c.subscriptions) < l {
		c.subscriptions = make([]string, l)
	} else {
		c.subscriptions = c.subscriptions[:l]
	}
	i := 0
	for s := range c.subscriptionsMap {
		c.subscriptions[i] = s
		i++
	}
	c.SG.Unlock()
	return nil
}

// UnsubscribeFromChannels unsubscribes from channels
func (c *Client) UnsubscribeFromChannels(channels []string) {
	var (
		pblcChannels []string
		prvtChannels []string
	)
	for _, v := range c.subscriptions {
		if _, ok := c.subscriptionsMap[v]; ok {
			if strings.HasPrefix(v, "user.") {
				prvtChannels = append(prvtChannels, v)
			} else {
				pblcChannels = append(pblcChannels, v)
			}
		}
	}
	if len(pblcChannels) > 0 {
		_, err := c.UnsubPblc(pblcChannels)
		if err != nil {
			go c.Logger.Println(err.Error())
		}
	}
	if len(prvtChannels) > 0 {
		_, err := c.UnsubPrvt(prvtChannels)
		if err != nil {
			go c.Logger.Println(err.Error())
		}
	}
}
