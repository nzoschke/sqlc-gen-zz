package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/olekukonko/errors"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"google.golang.org/protobuf/proto"
)

func main() {

	if err := run(); err != nil {
		slog.Error("run failed", "error", err)
		fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		os.Exit(2)
	}
}

func run() error {
	ctx := context.Background()

	f, err := os.OpenFile("sqlc-gen-zz.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	l := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(l)

	slog.Info("sqlc-gen-zz", "fn", "run")

	bs, err := io.ReadAll(os.Stdin)
	if err != nil {
		return errors.WithStack(err)
	}

	req := &plugin.GenerateRequest{}
	if err := proto.Unmarshal(bs, req); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	slog.Info("sqlc-gen-zz", "req", req)

	res, err := Gen(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}

	jres, err := Generate(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}

	res.Files = append(res.Files, jres.Files...)

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
