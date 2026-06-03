package validate

import (
	"testing"
	"time"
)

func mustTime(h, m int) time.Time {
	return time.Date(2026, 3, 25, h, m, 0, 0, time.FixedZone("ICT", 7*3600))
}

func TestInOnlineTimeRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		now     time.Time
		opening string
		closing string
		want    bool
		wantErr bool
	}{
		{
			name:    "both empty means always available",
			now:     mustTime(10, 0),
			opening: "",
			closing: "",
			want:    true,
			wantErr: false,
		},
		{
			name:    "opening empty only is invalid",
			now:     mustTime(10, 0),
			opening: "",
			closing: "22:00",
			want:    false,
			wantErr: true,
		},
		{
			name:    "closing empty only is invalid",
			now:     mustTime(10, 0),
			opening: "07:00",
			closing: "",
			want:    false,
			wantErr: true,
		},
		{
			name:    "normal range inside",
			now:     mustTime(10, 0),
			opening: "07:00",
			closing: "22:00",
			want:    true,
			wantErr: false,
		},
		{
			name:    "normal range before opening",
			now:     mustTime(6, 59),
			opening: "07:00",
			closing: "22:00",
			want:    false,
			wantErr: false,
		},
		{
			name:    "normal range at opening",
			now:     mustTime(7, 0),
			opening: "07:00",
			closing: "22:00",
			want:    true,
			wantErr: false,
		},
		{
			name:    "normal range just before closing",
			now:     mustTime(21, 59),
			opening: "07:00",
			closing: "22:00",
			want:    true,
			wantErr: false,
		},
		{
			name:    "normal range at closing is excluded",
			now:     mustTime(22, 0),
			opening: "07:00",
			closing: "22:00",
			want:    false,
			wantErr: false,
		},
		{
			name:    "overnight range before midnight inside",
			now:     mustTime(23, 0),
			opening: "22:00",
			closing: "06:00",
			want:    true,
			wantErr: false,
		},
		{
			name:    "overnight range after midnight inside",
			now:     mustTime(1, 0),
			opening: "22:00",
			closing: "06:00",
			want:    true,
			wantErr: false,
		},
		{
			name:    "overnight range daytime outside",
			now:     mustTime(12, 0),
			opening: "22:00",
			closing: "06:00",
			want:    false,
			wantErr: false,
		},
		{
			name:    "overnight range at closing is excluded",
			now:     mustTime(6, 0),
			opening: "22:00",
			closing: "06:00",
			want:    false,
			wantErr: false,
		},
		{
			name:    "same opening and closing means 24h open",
			now:     mustTime(15, 30),
			opening: "00:00",
			closing: "00:00",
			want:    true,
			wantErr: false,
		},
		{
			name:    "invalid opening format",
			now:     mustTime(10, 0),
			opening: "25:00",
			closing: "22:00",
			want:    false,
			wantErr: true,
		},
		{
			name:    "invalid closing format",
			now:     mustTime(10, 0),
			opening: "07:00",
			closing: "abc",
			want:    false,
			wantErr: true,
		},
		{
			name:    "trim spaces still works",
			now:     mustTime(10, 0),
			opening: " 07:00 ",
			closing: " 22:00 ",
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := InOnlineTimeRange(tt.now, tt.opening, tt.closing)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error state: gotErr=%v wantErr=%v err=%v", err != nil, tt.wantErr, err)
			}
			if got != tt.want {
				t.Fatalf("unexpected result: got=%v want=%v", got, tt.want)
			}
		})
	}
}

func TestInOnlineTimeRangeBool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		now     time.Time
		opening string
		closing string
		want    bool
	}{
		{
			name:    "both empty means always available",
			now:     mustTime(10, 0),
			opening: "",
			closing: "",
			want:    true,
		},
		{
			name:    "one side empty returns false",
			now:     mustTime(10, 0),
			opening: "07:00",
			closing: "",
			want:    false,
		},
		{
			name:    "normal range inside",
			now:     mustTime(10, 0),
			opening: "07:00",
			closing: "22:00",
			want:    true,
		},
		{
			name:    "normal range outside",
			now:     mustTime(23, 0),
			opening: "07:00",
			closing: "22:00",
			want:    false,
		},
		{
			name:    "overnight inside",
			now:     mustTime(1, 0),
			opening: "22:00",
			closing: "06:00",
			want:    true,
		},
		{
			name:    "overnight outside",
			now:     mustTime(12, 0),
			opening: "22:00",
			closing: "06:00",
			want:    false,
		},
		{
			name:    "same start end means always open",
			now:     mustTime(12, 0),
			opening: "08:00",
			closing: "08:00",
			want:    true,
		},
		{
			name:    "invalid format returns false",
			now:     mustTime(12, 0),
			opening: "25:00",
			closing: "22:00",
			want:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := InOnlineTimeRangeBool(tt.now, tt.opening, tt.closing)
			if got != tt.want {
				t.Fatalf("unexpected result: got=%v want=%v", got, tt.want)
			}
		})
	}
}

func TestParseHHMM(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in      string
		want    int
		wantErr bool
	}{
		{"00:00", 0, false},
		{"07:30", 450, false},
		{"23:59", 1439, false},
		{"24:00", 0, true},
		{"7:00", 420, false},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			got, err := parseHHMM(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error state: gotErr=%v wantErr=%v err=%v", err != nil, tt.wantErr, err)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("unexpected result: got=%d want=%d", got, tt.want)
			}
		})
	}
}

func BenchmarkInOnlineTimeRange_NormalInside(b *testing.B) {
	now := mustTime(10, 30)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ok, err := InOnlineTimeRange(now, "07:00", "22:00")
		if err != nil || !ok {
			b.Fatalf("unexpected result: ok=%v err=%v", ok, err)
		}
	}
}

func BenchmarkInOnlineTimeRange_OvernightInside(b *testing.B) {
	now := mustTime(1, 30)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ok, err := InOnlineTimeRange(now, "22:00", "06:00")
		if err != nil || !ok {
			b.Fatalf("unexpected result: ok=%v err=%v", ok, err)
		}
	}
}

func BenchmarkInOnlineTimeRangeBool_NormalInside(b *testing.B) {
	now := mustTime(10, 30)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !InOnlineTimeRangeBool(now, "07:00", "22:00") {
			b.Fatal("unexpected false")
		}
	}
}

func BenchmarkInOnlineTimeRangeBool_OvernightInside(b *testing.B) {
	now := mustTime(1, 30)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !InOnlineTimeRangeBool(now, "22:00", "06:00") {
			b.Fatal("unexpected false")
		}
	}
}
