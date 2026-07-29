package libreoffice

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	// ErrOverload: all workers busy and the job queue is full — transient,
	// caller retries (→ 503).
	ErrOverload = errors.New("libreoffice converter overloaded")
	// ErrSofficeNotFound: soffice binary missing — doc-to-pdf unavailable (→ 500).
	ErrSofficeNotFound = errors.New("soffice binary not found")
	// ErrInvalidInput: bad type/ext/size — caller error (→ 400).
	ErrInvalidInput = errors.New("invalid doc-to-pdf input")
)

type job struct {
	ctx context.Context
	src []byte
	ext string
	res chan<- result
}

type result struct {
	pdf []byte
	err error
}

// Converter runs headless soffice conversions over a fixed worker pool. Each
// worker owns a persistent hardened profile (so repeated conversions skip
// profile bootstrap) and recovers from panics by re-spawning. A bounded job
// channel bounds in-flight work; beyond cap+queue, Convert fails fast with
// ErrOverload so a caller's client timeout is never spent queueing.
//
// ponytail: a truly resident soffice process (zero-cold-start, no per-call
// exec) needs a UNO sidecar — infeasible in pure Go without adding python-uno
// to the image. Persistent worker profiles capture most of the cold-start win.
type Converter struct {
	jobs     chan job
	stop     chan struct{}
	wg       sync.WaitGroup
	timeout  time.Duration
	restarts atomic.Int64
}

// NewConverter starts `cap` workers, each with its own profile. Missing soffice
// returns ErrSofficeNotFound so the service stays up and surfaces 500 only on
// doc-to-pdf requests.
func NewConverter(cap int, timeout time.Duration) (*Converter, error) {
	if _, err := exec.LookPath("soffice"); err != nil {
		return nil, ErrSofficeNotFound
	}
	if cap < 1 {
		cap = 4
	}
	if timeout <= 0 {
		timeout = 90 * time.Second
	}
	// Queue depth = pool size: up to 2×cap in-flight (cap running + cap queued)
	// before requests are rejected with ErrOverload.
	c := &Converter{
		jobs:    make(chan job, cap),
		stop:    make(chan struct{}),
		timeout: timeout,
	}
	for i := 0; i < cap; i++ {
		c.spawnWorker()
	}
	return c, nil
}

// Convert submits srcBytes (a `.ext` document) to a worker and awaits its pdf.
// Queue saturated → ErrOverload; conversion/timeout failure → wrapped error.
func (c *Converter) Convert(ctx context.Context, srcBytes []byte, ext string) ([]byte, error) {
	res := make(chan result, 1)
	select {
	case c.jobs <- job{ctx: ctx, src: srcBytes, ext: ext, res: res}:
	default:
		return nil, ErrOverload
	}
	select {
	case r := <-res:
		return r.pdf, r.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Restarts reports how many workers have been re-spawned after a panic
// (crash-recovery health signal).
func (c *Converter) Restarts() int64 { return c.restarts.Load() }

// Close stops workers and removes their profile dirs; call on shutdown.
func (c *Converter) Close() {
	close(c.stop)
	c.wg.Wait()
}

// spawnWorker runs one worker; on panic it logs, counts the restart, and
// re-spawns a replacement (unless shutting down).
func (c *Converter) spawnWorker() {
	c.wg.Add(1)
	threading.GoSafe(func() {
		defer c.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				c.restarts.Add(1)
				logx.Errorf("libreoffice worker panicked: %v — restarting", r)
				select {
				case <-c.stop:
				default:
					c.spawnWorker()
				}
			}
		}()
		c.runWorker()
	})
}

func (c *Converter) runWorker() {
	profileDir, initErr := makeProfile()
	if profileDir != "" {
		defer func() { _ = os.RemoveAll(profileDir) }()
	}
	for {
		select {
		case <-c.stop:
			return
		case j := <-c.jobs:
			if initErr != nil {
				j.res <- result{err: initErr}
				continue
			}
			pdf, err := c.runSoffice(j.ctx, profileDir, j.src, j.ext)
			j.res <- result{pdf: pdf, err: err}
		}
	}
}

// makeProfile creates a hardened UserInstallation dir for one worker.
func makeProfile() (string, error) {
	dir, err := os.MkdirTemp("", "lo-profile-*")
	if err != nil {
		return "", errorx.Wrap(err, "create profile dir")
	}
	if err := ensureSafeProfile(dir); err != nil {
		_ = os.RemoveAll(dir)
		return "", errorx.Wrap(err, "write safe profile")
	}
	return dir, nil
}

func (c *Converter) runSoffice(ctx context.Context, profileDir string, srcBytes []byte, ext string) ([]byte, error) {
	tmp, err := os.MkdirTemp("", "lo-conv-*")
	if err != nil {
		return nil, errorx.Wrap(err, "create temp dir")
	}
	defer func() { _ = os.RemoveAll(tmp) }()

	srcPath := filepath.Join(tmp, "src."+ext)
	if err := os.WriteFile(srcPath, srcBytes, 0o600); err != nil {
		return nil, errorx.Wrap(err, "write source file")
	}

	cctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	//nolint:gosec // G204: soffice is a fixed binary and srcPath is a sanitized
	// temp file — no attacker-controlled shell input reaches the subprocess.
	cmd := exec.CommandContext(cctx, "soffice",
		"--headless", "--norestore", "--nologo", "--nofirststartwizard",
		"--convert-to", "pdf",
		"--outdir", tmp,
		"-env:UserInstallation=file://"+profileDir,
		srcPath,
	)
	// New process group so the ctx-deadline kill reaps soffice + its children
	// (oosplash/soffice.bin), preventing zombie/orphan buildup under load.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return os.ErrProcessDone
		}
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if cctx.Err() == context.DeadlineExceeded {
			return nil, errorx.Wrap(err, "soffice conversion timed out")
		}
		return nil, errorx.Wrap(err, "soffice conversion failed: "+stderr.String())
	}

	//nolint:gosec // G304: path under our own temp dir, fixed filename.
	pdf, err := os.ReadFile(filepath.Join(tmp, "src.pdf"))
	if err != nil {
		return nil, errorx.Wrap(err, "read converted pdf")
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		return nil, errorx.New("converted output is not a pdf")
	}
	return pdf, nil
}
