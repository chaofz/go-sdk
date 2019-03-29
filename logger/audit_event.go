package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/blend/go-sdk/ansi"
)

// these are compile time assertions
var (
	_ Event          = (*AuditEvent)(nil)
	_ TextWritable   = (*AuditEvent)(nil)
	_ json.Marshaler = (*AuditEvent)(nil)
)

// NewAuditEvent returns a new audit event.
func NewAuditEvent(principal, verb string, options ...EventMetaOption) *AuditEvent {
	return &AuditEvent{
		EventMeta: NewEventMeta(Audit, options...),
		Principal: principal,
		Verb:      verb,
	}
}

// NewAuditEventListener returns a new audit event listener.
func NewAuditEventListener(listener func(context.Context, *AuditEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*AuditEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// AuditEvent is a common type of event detailing a business action by a subject.
type AuditEvent struct {
	*EventMeta

	Context       string
	Principal     string
	Verb          string
	Noun          string
	Subject       string
	Property      string
	RemoteAddress string
	UserAgent     string
	Extra         map[string]string
}

// WithContext sets a field.
func (e *AuditEvent) WithContext(context string) *AuditEvent {
	e.Context = context
	return e
}

// WithPrincipal sets a field.
func (e *AuditEvent) WithPrincipal(principal string) *AuditEvent {
	e.Principal = principal
	return e
}

// WithVerb sets a field.
func (e *AuditEvent) WithVerb(verb string) *AuditEvent {
	e.Verb = verb
	return e
}

// WithNoun sets a field.
func (e *AuditEvent) WithNoun(noun string) *AuditEvent {
	e.Noun = noun
	return e
}

// WithSubject sets a field.
func (e *AuditEvent) WithSubject(subject string) *AuditEvent {
	e.Subject = subject
	return e
}

// WithProperty sets a field.
func (e *AuditEvent) WithProperty(property string) *AuditEvent {
	e.Property = property
	return e
}

// WithRemoteAddress sets the remote address.
func (e *AuditEvent) WithRemoteAddress(remoteAddr string) *AuditEvent {
	e.RemoteAddress = remoteAddr
	return e
}

// WithUserAgent sets the user agent.
func (e *AuditEvent) WithUserAgent(userAgent string) *AuditEvent {
	e.UserAgent = userAgent
	return e
}

// WithExtra sets the extra info.
func (e *AuditEvent) WithExtra(extra map[string]string) *AuditEvent {
	e.Extra = extra
	return e
}

// WriteText implements TextWritable.
func (e AuditEvent) WriteText(formatter TextFormatter, wr io.Writer) {
	if len(e.Context) > 0 {
		io.WriteString(wr, formatter.Colorize("Context:", ansi.ColorGray))
		io.WriteString(wr, e.Context)
		io.WriteString(wr, Space)
	}
	if len(e.Principal) > 0 {
		io.WriteString(wr, formatter.Colorize("Principal:", ansi.ColorGray))
		io.WriteString(wr, e.Principal)
		io.WriteString(wr, Space)
	}
	if len(e.Verb) > 0 {
		io.WriteString(wr, formatter.Colorize("Verb:", ansi.ColorGray))
		io.WriteString(wr, e.Verb)
		io.WriteString(wr, Space)
	}
	if len(e.Noun) > 0 {
		io.WriteString(wr, formatter.Colorize("Noun:", ansi.ColorGray))
		io.WriteString(wr, e.Noun)
		io.WriteString(wr, Space)
	}
	if len(e.Subject) > 0 {
		io.WriteString(wr, formatter.Colorize("Subject:", ansi.ColorGray))
		io.WriteString(wr, e.Subject)
		io.WriteString(wr, Space)
	}
	if len(e.Property) > 0 {
		io.WriteString(wr, formatter.Colorize("Property:", ansi.ColorGray))
		io.WriteString(wr, e.Property)
		io.WriteString(wr, Space)
	}
	if len(e.RemoteAddress) > 0 {
		io.WriteString(wr, formatter.Colorize("Remote Addr:", ansi.ColorGray))
		io.WriteString(wr, e.RemoteAddress)
		io.WriteString(wr, Space)
	}
	if len(e.UserAgent) > 0 {
		io.WriteString(wr, formatter.Colorize("UA:", ansi.ColorGray))
		io.WriteString(wr, e.UserAgent)
		io.WriteString(wr, Space)
	}
	if len(e.Extra) > 0 {
		var values []string
		for key, value := range e.Extra {
			values = append(values, fmt.Sprintf("%s%s", formatter.Colorize(key+":", ansi.ColorGray), value))
		}
		io.WriteString(wr, strings.Join(values, " "))
	}
}

// MarshalJSON implements json.Marshaler.
func (e AuditEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"context":    e.Context,
		"principal":  e.Principal,
		"verb":       e.Verb,
		"noun":       e.Noun,
		"subject":    e.Subject,
		"property":   e.Property,
		"remoteAddr": e.RemoteAddress,
		"ua":         e.UserAgent,
		"extra":      e.Extra,
	}))
}
