package languageserver

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	ortfodb "github.com/ortfo/db"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var YAMLSeparator = regexp.MustCompile(ortfodb.PatternYAMLSeparator)
var logger *zap.Logger

type state struct {
	config       ortfodb.Configuration
	tags         []yaml.Node
	technologies []yaml.Node
}

type Handler struct {
	protocol.Server
	Logger *zap.Logger
}

func (h Handler) state(ctx context.Context) state {
	return ctx.Value("state").(state)
}

func (h Handler) config(ctx context.Context) ortfodb.Configuration {
	return h.state(ctx).config
}

func NewHandler(ctx context.Context, server protocol.Server, logger *zap.Logger) (Handler, context.Context, error) {
	config, err := ortfodb.NewConfiguration(ctx.Value("configpath").(string))
	if err != nil {
		return Handler{}, ctx, fmt.Errorf("while loading ortfodb configuration from ./ortfodb.yaml: %w", err)
	}

	tags, err := LoadRepository(config.Tags.Repository)
	if err != nil {
		return Handler{}, ctx, fmt.Errorf("while loading tags from repository: %w", err)
	}

	technologies, err := LoadRepository(config.Technologies.Repository)
	if err != nil {
		return Handler{}, ctx, fmt.Errorf("while loading technologies from repository: %w", err)
	}

	return Handler{
			Server: server,
			Logger: logger,
		}, context.WithValue(ctx, "state", state{
			config:       config,
			tags:         tags,
			technologies: technologies,
		}), nil
}

func (h Handler) Initialize(ctx context.Context, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	logger = h.Logger
	h.Logger.Debug("Initializing ortfols server")
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			DefinitionProvider: true,
			HoverProvider:      true,
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "ortfols",
			Version: ortfodb.Version,
		},
	}, nil
}

func (h Handler) Definition(ctx context.Context, params *protocol.DefinitionParams) ([]protocol.Location, error) {
	h.Logger.Debug("LSP:Definition", zap.Any("state", h.state), zap.Any("params", params))
	file, err := CurrentFile(params.TextDocumentPositionParams)
	if err != nil {
		return []protocol.Location{}, fmt.Errorf("while getting current file: %w", err)
	}

	if key, node, inside := file.InFrontmatter(); inside {
		h.Logger.Debug("Found frontmatter key", zap.String("key", key), zap.Any("node", node))
		switch key {
		case "tags":
			node, _, err := FindInRepository[ortfodb.Tag](node.Value, "tag", h.state(ctx).tags)
			pos := positionOf(node)
			return []protocol.Location{
				{
					URI: uri.File(h.config(ctx).Tags.Repository),
					Range: protocol.Range{
						Start: pos,
						End:   pos,
					},
				},
			}, err
		case "made with":
			node, _, err := FindInRepository[ortfodb.Technology](node.Value, "technology", h.state(ctx).technologies)
			pos := positionOf(node)
			return []protocol.Location{
				{
					URI: uri.File(h.config(ctx).Technologies.Repository),
					Range: protocol.Range{
						Start: pos,
						End:   pos,
					},
				},
			}, err
		}
	}
	return []protocol.Location{}, nil
}

func (h Handler) Initialized(ctx context.Context, params *protocol.InitializedParams) error {
	return nil
}

func (h Handler) Shutdown(ctx context.Context) error {
	return nil
}

func (h Handler) Exit(ctx context.Context) error {
	return nil
}

func (h Handler) WorkDoneProgressCancel(ctx context.Context, params *protocol.WorkDoneProgressCancelParams) error {
	return errors.New("unimplemented")
}

func (h Handler) LogTrace(ctx context.Context, params *protocol.LogTraceParams) error {
	return errors.New("unimplemented")
}

func (h Handler) SetTrace(ctx context.Context, params *protocol.SetTraceParams) error {
	return errors.New("unimplemented")
}

func (h Handler) CodeAction(ctx context.Context, params *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) CodeLens(ctx context.Context, params *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) CodeLensResolve(ctx context.Context, params *protocol.CodeLens) (*protocol.CodeLens, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) ColorPresentation(ctx context.Context, params *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) CompletionResolve(ctx context.Context, params *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Declaration(ctx context.Context, params *protocol.DeclarationParams) ([]protocol.Location, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	return errors.New("unimplemented")
}

func (h Handler) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	return errors.New("unimplemented")
}

func (h Handler) DidChangeWatchedFiles(ctx context.Context, params *protocol.DidChangeWatchedFilesParams) error {
	return errors.New("unimplemented")
}

func (h Handler) DidChangeWorkspaceFolders(ctx context.Context, params *protocol.DidChangeWorkspaceFoldersParams) error {
	return errors.New("unimplemented")
}

func (h Handler) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	return errors.New("unimplemented")
}

func (h Handler) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	return errors.New("unimplemented")
}

func (h Handler) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	return errors.New("unimplemented")
}

func (h Handler) DocumentColor(ctx context.Context, params *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DocumentHighlight(ctx context.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DocumentLink(ctx context.Context, params *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DocumentLinkResolve(ctx context.Context, params *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DocumentSymbol(ctx context.Context, params *protocol.DocumentSymbolParams) ([]interface{}, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) ExecuteCommand(ctx context.Context, params *protocol.ExecuteCommandParams) (interface{}, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) FoldingRanges(ctx context.Context, params *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	h.Logger.Debug("LSP:Hover", zap.Any("state", h.state(ctx)), zap.Any("params", params))
	file, err := CurrentFile(params.TextDocumentPositionParams)
	if err != nil {
		return nil, fmt.Errorf("while getting current file: %w", err)
	}

	if key, node, inside := file.InFrontmatter(); inside {
		h.Logger.Debug("Found frontmatter key", zap.String("key", key), zap.Any("node", node))
		switch key {
		case "tags":
			_, tag, err := FindInRepository[ortfodb.Tag](node.Value, "tag", h.state(ctx).tags)
			if err != nil {
				return nil, err
			}
			return &protocol.Hover{
				Contents: ReferrableDescription(tag, tag.Description),
			}, nil
		case "made with":
			_, technology, err := FindInRepository[ortfodb.Technology](node.Value, "technology", h.state(ctx).technologies)
			if err != nil {
				return nil, err
			}
			return &protocol.Hover{
				Contents: ReferrableDescription(technology, technology.Description),
			}, nil

		}
	}
	return nil, nil
}

func (h Handler) Implementation(ctx context.Context, params *protocol.ImplementationParams) ([]protocol.Location, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) OnTypeFormatting(ctx context.Context, params *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) PrepareRename(ctx context.Context, params *protocol.PrepareRenameParams) (*protocol.Range, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) RangeFormatting(ctx context.Context, params *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) References(ctx context.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Rename(ctx context.Context, params *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) SignatureHelp(ctx context.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Symbols(ctx context.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) TypeDefinition(ctx context.Context, params *protocol.TypeDefinitionParams) ([]protocol.Location, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) WillSave(ctx context.Context, params *protocol.WillSaveTextDocumentParams) error {
	return errors.New("unimplemented")
}

func (h Handler) WillSaveWaitUntil(ctx context.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) ShowDocument(ctx context.Context, params *protocol.ShowDocumentParams) (*protocol.ShowDocumentResult, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) WillCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DidCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) error {
	return errors.New("unimplemented")
}

func (h Handler) WillRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DidRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) error {
	return errors.New("unimplemented")
}

func (h Handler) WillDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DidDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) error {
	return errors.New("unimplemented")
}

func (h Handler) CodeLensRefresh(ctx context.Context) error {
	return errors.New("unimplemented")
}

func (h Handler) PrepareCallHierarchy(ctx context.Context, params *protocol.CallHierarchyPrepareParams) ([]protocol.CallHierarchyItem, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) IncomingCalls(ctx context.Context, params *protocol.CallHierarchyIncomingCallsParams) ([]protocol.CallHierarchyIncomingCall, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) OutgoingCalls(ctx context.Context, params *protocol.CallHierarchyOutgoingCallsParams) ([]protocol.CallHierarchyOutgoingCall, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) SemanticTokensFull(ctx context.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) SemanticTokensFullDelta(ctx context.Context, params *protocol.SemanticTokensDeltaParams) (interface{}, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) SemanticTokensRange(ctx context.Context, params *protocol.SemanticTokensRangeParams) (*protocol.SemanticTokens, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) SemanticTokensRefresh(ctx context.Context) error {
	return errors.New("unimplemented")
}

func (h Handler) LinkedEditingRange(ctx context.Context, params *protocol.LinkedEditingRangeParams) (*protocol.LinkedEditingRanges, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Moniker(ctx context.Context, params *protocol.MonikerParams) ([]protocol.Moniker, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Request(ctx context.Context, method string, params interface{}) (interface{}, error) {
	return nil, errors.New("unimplemented")
}
