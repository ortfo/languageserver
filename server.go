package languageserver

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var BuiltAt string

// StartServer starts the language server. It reads from stdin and writes to stdout.
// If logClientIn is not empty, it will log the client's request and responses to respectively client-request-from.log and client-response-to.log, in the directory specified by logClientIn.
func StartServer(logger *zap.Logger, configurationPath string, logClientIn string) {
	conn := jsonrpc2.NewConn(jsonrpc2.NewStream(&readWriteCloser{
		reader: os.Stdin,
		writer: os.Stdout,
		logAt:  logClientIn,
	}))
	// notifier := protocol.ClientDispatcher(conn, logger.Named("notify"))
	// handler := languageserver.Handler{
	// 	Logger: logger,
	// 	Server: protocol.ServerDispatcher(conn, logger),
	// }
	handler, ctx, err := NewHandler(context.WithValue(context.Background(), "configpath", configurationPath), protocol.ServerDispatcher(conn, logger), logger)
	if err != nil {
		logger.Sugar().Fatalf("while initializing handler: %w", err)
	}

	conn.Go(ctx, protocol.ServerHandler(handler, jsonrpc2.MethodNotFoundHandler))
	<-conn.Done()
}

type readWriteCloser struct {
	reader io.ReadCloser
	writer io.WriteCloser
	logAt  string
}

func (r *readWriteCloser) Read(b []byte) (int, error) {
	var f *os.File
	if r.logAt != "" {
		f, _ = os.OpenFile(filepath.Join(r.logAt, "client-request-from.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	n, err := r.reader.Read(b)
	if r.logAt != "" {
		if err != nil {
			f.Write([]byte(err.Error() + "\n"))
		} else {
			f.Write(b)
		}
	}
	return n, err
}

func (r *readWriteCloser) Write(b []byte) (int, error) {
	if r.logAt != "" {
		f, _ := os.OpenFile(filepath.Join(r.logAt, "client-response-to.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.Write(b)
	}
	return r.writer.Write(b)
}

func (r *readWriteCloser) Close() error {
	return multierr.Append(r.reader.Close(), r.writer.Close())
}
