package main

const dummyNote = `
# Binary Search

Binary search works on sorted arrays by repeatedly halving the search space.
Instead of checking every element, it compares the target against the middle
element and eliminates half the array on each step. Time complexity is O(log n).

## Implementation

` + "```" + `go
func binarySearch(arr []int, target int) int {
	low, high := 0, len(arr)-1

	for low <= high {
		mid := (low + high) / 2

		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1 // not found
}
` + "```" + `

## Usage

` + "```" + `go
arr := []int{1, 3, 5, 7, 9, 11, 13}
idx := binarySearch(arr, 7)
fmt.Println(idx) // 3
` + "```" + `

## Notes

- Array must be sorted before calling
- Returns -1 if target is not found
- Works on any ordered data type

## Graceful Shutdown

` + "```" + `go

func (b *backend) serve() error {
	e := echo.New()
	defer func() { _ = e.Close() }()

	renderer, err := NewCacheRenderer()
	if err != nil {
		return err
	}
	b.renderer = renderer
	e.Renderer = renderer
	b.setupRoutes(e)
	srv := http.Server{
		Handler: e, Addr: b.conf.addr,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       90 * time.Second,
	}
	shutdownErr := make(chan error, 1)
	go func() {
		q := make(chan os.Signal, 1)
		signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
		s := <-q
		b.zlog.Info().
			Str("signal", s.String()).
			Msg("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err = srv.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}
		b.zlog.Info().
			Str("address", srv.Addr).
			Msg("completing background tasks")

		b.wg.Wait()
		shutdownErr <- nil
	}()
	b.zlog.Info().
		Str("env", b.conf.env).
		Str("addr", srv.Addr).Msg("starting server")

	if err = srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("couldn't listen on addr=%s: %w", srv.Addr, err)
		}
	}
	if err = <-shutdownErr; err != nil {
		return err
	}
	b.zlog.Info().
		Str("env", b.conf.env).
		Str("addr", srv.Addr).Msg("server stopped")
	return nil
}
` + "```" + ``
