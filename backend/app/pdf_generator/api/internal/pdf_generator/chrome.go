package pdf_generator

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	DefaultMaxConcurrent  = 4
	DefaultTimeout        = 60 * time.Second
	DefaultAcquireTimeout = 10 * time.Second
	DefaultWaitTime       = 1000 * time.Millisecond
)

// ChromePool manages a single Chrome instance with concurrent tab generation
type ChromePool struct {
	allocCtx       context.Context
	allocCancel    context.CancelFunc
	sem            chan struct{}
	maxConcurrent  int
	mu             sync.Mutex
	initialized    bool
	closed         bool
	acquireTimeout time.Duration
	restarting     bool
}

// NewChromePool creates a new Chrome pool
func NewChromePool(maxConcurrent int, acquireTimeout time.Duration) (*ChromePool, error) {
	if maxConcurrent <= 0 {
		maxConcurrent = DefaultMaxConcurrent
	}
	if acquireTimeout <= 0 {
		acquireTimeout = DefaultAcquireTimeout
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // classic headless
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	pool := &ChromePool{
		allocCtx:       allocCtx,
		allocCancel:    allocCancel,
		sem:            make(chan struct{}, maxConcurrent),
		maxConcurrent:  maxConcurrent,
		initialized:    true,
		acquireTimeout: acquireTimeout,
	}

	// verify Chrome instance
	if err := pool.verify(); err != nil {
		_ = pool.Close()
		return nil, fmt.Errorf("chrome pool initialization failed: %w", err)
	}

	logx.Infof("Chrome pool initialized (maxConcurrent=%d)", maxConcurrent)
	return pool, nil
}

// verify navigates to blank page to ensure Chrome is alive
func (p *ChromePool) verify() error {
	tabCtx, tabCancel := chromedp.NewContext(p.allocCtx)
	defer tabCancel()
	ctx, cancel := context.WithTimeout(tabCtx, DefaultTimeout)
	defer cancel()
	return chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
	)
}

// acquireWithTimeout handles semaphore acquire with timeout
func (p *ChromePool) acquireWithTimeout(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, p.acquireTimeout)
	defer cancel()

	select {
	case p.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("chrome busy: acquire timeout")
		}
		return ctx.Err()
	}
}

// Release releases a semaphore slot
func (p *ChromePool) Release() {
	select {
	case <-p.sem:
	default:
	}
}

// URLToPDF generates PDF from URL
func (p *ChromePool) URLToPDF(ctx context.Context, url string, opts ...PDFOption) ([]byte, error) {
	if err := p.acquireWithTimeout(ctx); err != nil {
		return nil, err
	}
	defer p.Release()

	cfg := &pdfConfig{waitTime: DefaultWaitTime}
	for _, opt := range opts {
		opt(cfg)
	}

	tabCtx, tabCancel := chromedp.NewContext(p.allocCtx)
	defer tabCancel()

	var pdfBuf []byte
	err := chromedp.Run(tabCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body"),
		chromedp.Sleep(cfg.waitTime),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				Do(ctx)
			return err
		}),
	)

	if isChromeFatal(err) {
		p.tryRestart(err)
	}

	return pdfBuf, err
}

// HTMLToPDF generates PDF from HTML string
func (p *ChromePool) HTMLToPDF(ctx context.Context, html string, opts ...PDFOption) ([]byte, error) {
	if err := p.acquireWithTimeout(ctx); err != nil {
		return nil, err
	}
	defer p.Release()

	cfg := &pdfConfig{waitTime: DefaultWaitTime}
	for _, opt := range opts {
		opt(cfg)
	}

	// Compress HTML: remove newlines and excess whitespace
	html = compressHTML(html)

	tabCtx, tabCancel := chromedp.NewContext(p.allocCtx)
	defer tabCancel()

	var pdfBuf []byte
	err := chromedp.Run(tabCtx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		// chromedp.WaitVisible("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var readyState string
			for {
				if err := chromedp.Evaluate(`document.readyState`, &readyState).Do(ctx); err != nil {
					return err
				}
				if readyState == "complete" {
					return nil
				}
				time.Sleep(50 * time.Millisecond)
			}
		}),
		chromedp.Sleep(cfg.waitTime),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			return err
		}),
	)

	if isChromeFatal(err) {
		p.tryRestart(err)
	}

	return pdfBuf, err
}

// GeneratePDF selects type and generates PDF
func (p *ChromePool) GeneratePDF(ctx context.Context, pdfType, content string, opts ...PDFOption) ([]byte, error) {
	switch pdfType {
	case "url":
		return p.URLToPDF(ctx, content, opts...)
	case "text":
		return p.HTMLToPDF(ctx, content, opts...)
	default:
		return nil, fmt.Errorf("invalid pdf type: %s", pdfType)
	}
}

// Close closes the Chrome pool
func (p *ChromePool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	p.initialized = false

	if p.allocCancel != nil {
		p.allocCancel()
	}

	logx.Info("Chrome pool closed")
	return nil
}

// isChromeFatal detects Chrome crashes or disconnections
func isChromeFatal(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return strings.Contains(msg, "target closed") ||
			strings.Contains(msg, "session deleted") ||
			strings.Contains(msg, "inspector disconnected") ||
			strings.Contains(msg, "chrome not reachable") ||
			strings.Contains(msg, "browser has crashed")
	}

	return strings.Contains(msg, "target closed") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "websocket") ||
		strings.Contains(msg, "browser has disconnected") ||
		strings.Contains(msg, "chrome not reachable") ||
		strings.Contains(msg, "session deleted") ||
		strings.Contains(msg, "inspector disconnected") ||
		strings.Contains(msg, "browser process") ||
		strings.Contains(msg, "chrome failed to start")
}

// tryRestart restarts the Chrome process on fatal error
func (p *ChromePool) tryRestart(cause error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.restarting {
		return
	}
	p.restarting = true
	defer func() { p.restarting = false }()

	logx.Errorf("chrome fatal error: %v, restarting...", cause)

	if p.allocCancel != nil {
		p.allocCancel()
		p.allocCancel = nil
	}
	p.allocCtx = nil

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	tabCtx, tabCancel := chromedp.NewContext(allocCtx)
	defer tabCancel()
	ctx, cancel := context.WithTimeout(tabCtx, DefaultTimeout)
	defer cancel()
	if err := chromedp.Run(ctx, chromedp.Navigate("about:blank")); err != nil {
		allocCancel()
		logx.Errorf("chrome warm-up failed: %v", err)
		return
	}

	p.allocCtx = allocCtx
	p.allocCancel = allocCancel
	p.initialized = true
	logx.Info("chrome restarted successfully")
}

// IsInitialized returns pool status
func (p *ChromePool) IsInitialized() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.initialized
}

// Stats returns pool statistics
func (p *ChromePool) Stats() PoolStats {
	p.mu.Lock()
	defer p.mu.Unlock()
	return PoolStats{
		MaxConcurrent: p.maxConcurrent,
		ActiveTasks:   len(p.sem),
		Initialized:   p.initialized,
		Closed:        p.closed,
	}
}

// PoolStats represents pool statistics
type PoolStats struct {
	MaxConcurrent int
	ActiveTasks   int
	Initialized   bool
	Closed        bool
}

// PDFOption is the option type for PDF generation
type PDFOption func(*pdfConfig)

type pdfConfig struct {
	waitTime time.Duration
}

// WithWaitTime sets the wait time for page load
func WithWaitTime(d time.Duration) PDFOption {
	return func(c *pdfConfig) {
		c.waitTime = d
	}
}

// compressHTML removes newlines and excess whitespace from HTML for chromedp
func compressHTML(html string) string {
	// Remove carriage returns and newlines
	html = strings.ReplaceAll(html, "\r\n", "")
	html = strings.ReplaceAll(html, "\n", " ")
	html = strings.ReplaceAll(html, "\r", " ")
	// Fix escaped quotes from JSON transport: \" → "
	// html = strings.ReplaceAll(html, "\\\"", "\"")

	// Replace multiple spaces with single space (except within tags)
	var result strings.Builder
	inTag := false
	prevSpace := false

	for _, r := range html {
		if r == '<' {
			inTag = true
			result.WriteRune(r)
			prevSpace = false
		} else if r == '>' {
			inTag = false
			result.WriteRune(r)
			prevSpace = false
		} else if inTag {
			// Inside tags, keep original spacing
			result.WriteRune(r)
			prevSpace = false
		} else if r == ' ' || r == '\t' {
			// Outside tags, collapse multiple spaces to one
			if !prevSpace {
				result.WriteRune(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}

	return result.String()
}
