package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/errors"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"google.golang.org/protobuf/proto"
)

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		os.Exit(2)
	}
}

func run() error {
	ctx := context.Background()

	bs, err := io.ReadAll(os.Stdin)
	if err != nil {
		return errors.WithStack(err)
	}

	req := &plugin.GenerateRequest{}
	if err := proto.Unmarshal(bs, req); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	res, err := Gen(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}

	bs, err = proto.Marshal(res)
	if err != nil {
		return errors.WithStack(err)
	}

	w := bufio.NewWriter(os.Stdout)
	if _, err := w.Write(bs); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}
