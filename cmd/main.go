package main

import (
	"context"
	"io"
	"os"

	"github.com/ortfo/languageserver"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	// ortfodb "github.com/ortfo/db"
)

func main() {
	logconf := zap.NewDevelopmentConfig()
	logconf.OutputPaths = []string{"/home/uwun/projects/ortfo/languageserver/logs/server.log"}
	logger, _ := logconf.Build()

	conn := jsonrpc2.NewConn(jsonrpc2.NewStream(&readWriteCloser{
		reader: os.Stdin,
		writer: os.Stdout,
	}))
	// notifier := protocol.ClientDispatcher(conn, logger.Named("notify"))
	// handler := languageserver.Handler{
	// 	Logger: logger,
	// 	Server: protocol.ServerDispatcher(conn, logger),
	// }
	handler, ctx, err := languageserver.NewHandler(context.Background(), protocol.ServerDispatcher(conn, logger), logger)
	if err != nil {
		logger.Sugar().Fatalf("while initializing handler: %w", err)
	}

	conn.Go(ctx, protocol.ServerHandler(handler, jsonrpc2.MethodNotFoundHandler))
	<-conn.Done()
}

type readWriteCloser struct {
	reader io.ReadCloser
	writer io.WriteCloser
}

func (r *readWriteCloser) Read(b []byte) (int, error) {
	f, _ := os.OpenFile("/home/uwun/projects/ortfo/languageserver/logs/client-request-from.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	n, err := r.reader.Read(b)
	if err != nil {
		f.Write([]byte(err.Error() + "\n"))
	} else {
		f.Write(b)
	}
	return n, err
}

func (r *readWriteCloser) Write(b []byte) (int, error) {
	f, _ := os.OpenFile("/home/uwun/projects/ortfo/languageserver/logs/client-response-to.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.Write(b)
	return r.writer.Write(b)
}

func (r *readWriteCloser) Close() error {
	return multierr.Append(r.reader.Close(), r.writer.Close())
}
