package transmissionrss

import (
	"reflect"
	"testing"
)

func TestFeed_Parse(t *testing.T) {
	type args struct {
		xml string
	}
	tests := []struct {
		name      string
		f         *Feed
		args      args
		wantItems []FeedItem
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feed{}
			gotItems, err := f.Parse(tt.args.xml)
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotItems, tt.wantItems) {
				t.Errorf("Feed.Parse() = %v, want %v", gotItems, tt.wantItems)
			}
		})
	}
}
