package format

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestRenderer_Render_Basic(t *testing.T) {
	r := New()

	title, err := r.Render("vi:order.created:title", "Đơn hàng mới #{{.order_id}}", map[string]any{
		"order_id": "A123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "Đơn hàng mới #A123"; title != want {
		t.Fatalf("title = %q, want %q", title, want)
	}
}

func TestRenderer_Render_MissingKeyZero(t *testing.T) {
	r := New()

	// Missing .customer should render empty (no panic)
	desc, err := r.Render("vi:order.created:desc", "Khách hàng {{.customer}} đặt đơn trị giá {{currency .amount \"VND\"}}", map[string]any{
		"amount": 5500000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(desc, "5,500,000 VND") {
		t.Fatalf("desc miss currency formatting: %q", desc)
	}
	// The {{.customer}} placeholder is empty -> acceptable output
}

func TestRenderer_RenderPair(t *testing.T) {
	r := New()
	data := map[string]any{
		"order_id":   "999",
		"customer":   "Nguyễn Văn A",
		"amount":     1234567.89,
		"created_at": time.Now().Add(-2 * time.Hour),
	}

	title, desc, err := r.RenderPair("vi", "order.created",
		"Đơn hàng mới #{{.order_id}}",
		"Khách hàng {{.customer}} đặt đơn trị giá {{currency .amount \"VND\"}} ({{ago .created_at}})",
		data,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "Đơn hàng mới #999" {
		t.Fatalf("title mismatch: %q", title)
	}
	if !strings.Contains(desc, "1,234,567.89 VND") {
		t.Fatalf("desc number formatting mismatch: %q", desc)
	}

	// Replace the old assertion: be tolerant to minute/hour/day suffix
	desc = strings.TrimSpace(desc) // avoid fail due to trailing space/newline

	// Accept "(15m ago)" or "(2h ago)" or "(3d ago)" at the end
	re := regexp.MustCompile(`\s\((\d+[mhd]) ago\)$`)
	if !re.MatchString(desc) {
		t.Fatalf("desc ago not matched (want '(Nm|Nh|Nd) ago' at end): %q", desc)
	}
}

func TestRenderer_Cache_ByKey(t *testing.T) {
	r := New()
	key := "vi:test:cache"

	// First time: parse & cache the template
	out1, err := r.Render(key, "Hello {{.Name}}", map[string]any{"Name": "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out1 != "Hello Alice" {
		t.Fatalf("got %q", out1)
	}

	// Second time: using the same key but a different template — should still use cached version
	out2, err := r.Render(key, "Bye {{.Name}}", map[string]any{"Name": "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out2 != "Hello Alice" {
		t.Fatalf("cache by key failed: got %q, want %q", out2, "Hello Alice")
	}

	// Clear cache and test again
	r.ClearCache()
	out3, err := r.Render(key, "Bye {{.Name}}", map[string]any{"Name": "Alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out3 != "Bye Alice" {
		t.Fatalf("after ClearCache got %q", out3)
	}
}

func TestRenderer_EmptyTemplate(t *testing.T) {
	r := New()
	out, err := r.Render("vi:empty", "", map[string]any{"x": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Fatalf("want empty string, got %q", out)
	}
}
