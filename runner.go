package rac

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

type runner struct {
	Timeout         time.Duration
	TryTimeoutCount int
}

func (r *runner) RunCtx(ctx context.Context, command string, args []string) (respond []byte, err error) {

	respond, err = r.runTry(ctx, command, args)

	if err == context.DeadlineExceeded {

		for tryCount := 1; tryCount < r.TryTimeoutCount; tryCount++ {

			respond, err = r.runTry(ctx, command, args)

			if err == context.DeadlineExceeded {
				continue
			}

			return

		}

	}

	return

}

func (r *runner) runTry(ctx context.Context, command string, args []string) (respond []byte, err error) {

	ctx, _ = context.WithTimeout(ctx, r.Timeout)

	cmd := exec.CommandContext(ctx, command, args...)

	cmd.Stdout = new(bytes.Buffer)
	cmd.Stderr = new(bytes.Buffer)
	errch := make(chan error, 1)

	err = cmd.Start()
	if err != nil {
		return respond, fmt.Errorf("Произошла ошибка запуска:\n\terr:%v\n\tПараметры: %v\n\t", err.Error(), cmd.Args)
	}

	go func() {
		errch <- cmd.Wait()
	}()

	select {
	case <-ctx.Done(): // timeout

		return respond, ctx.Err()

	case err := <-errch:
		if err != nil {

			stderr := cmd.Stderr.(*bytes.Buffer).Bytes()
			errText := fmt.Sprintf("Произошла ошибка запуска:\n\terr:%v\n\tПараметры: %v\n\t", err.Error(), cmd.Args)

			stdErrText, _ := decodeOutBytes(stderr)

			if len(stderr) > 0 {
				errText += fmt.Sprintf("StdErr:%v\n", stdErrText)
			}

			return respond, errors.New(errText)

		} else {

			in := cmd.Stdout.(*bytes.Buffer).Bytes()

			respond, err = decodeOutBytes(in)

			if err != nil {
				return respond, err
			}

			return respond, nil
		}
	}
}
