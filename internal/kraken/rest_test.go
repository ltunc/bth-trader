package kraken

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"
)

func TestRestClient_sign(t *testing.T) {
	type fields struct {
		apiKey     string
		privateKey string
		decodedKey []byte
		baseUrl    string
		httpClient *http.Client
	}
	type args struct {
		uriPath string
		data    url.Values
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "base",
			fields: fields{
				apiKey:     "notneededhere",
				privateKey: "kQH5HW/8p1uGOVjbgWA7FunAmGO8lsSUXNsu3eow76sz84Q18fWxnyRzBHCd3pd5nE9qa99HAZtuZuj6F1huXg==",
			},
			args: args{
				uriPath: "/0/private/AddOrder",
				data: url.Values{
					"nonce":     {"1616492376594"},
					"ordertype": {"limit"},
					"pair":      {"XBTUSD"},
					"price":     {"37500"},
					"type":      {"buy"},
					"volume":    {"1.25"},
				},
			},
			want: "4/dpxb3iT4tp/ZCVEwSnEsLxx0bqyhLpdfOpc6fn7OR8+UClSV5n9E6aSS8MPtnRfp32bAb0nmbRn6H8ndwLUQ==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestClient{
				apiKey:     tt.fields.apiKey,
				privateKey: tt.fields.privateKey,
			}
			r.decodedKey, _ = base64.StdEncoding.DecodeString(tt.fields.privateKey)
			if got := r.sign(tt.args.uriPath, tt.args.data); got != tt.want {
				t.Errorf("sign() = %v, want %v", got, tt.want)
			}
		})
	}
}
