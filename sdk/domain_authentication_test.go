package sendgrid

import (
	"reflect"
	"testing"
)

func Test_parseDomainAuthentication(t *testing.T) {
	type args struct {
		respBody string
	}
	tests := []struct {
		name  string
		args  args
		want  *DomainAuthentication
		want1 RequestError
	}{
		{
			name: "not automatic security",
			args: args{respBody: `{"id":123,"user_id":234,"subdomain":"aaa","domain":"example.com","username":"user","ips":[],"custom_spf":false,"default":false,"legacy":false,"automatic_security":false,"valid":false,"dns":{"mail_server":{"valid":false,"type":"mx","host":"aaa.example.com","data":"mx.sendgrid.net."},"subdomain_spf":{"valid":false,"type":"txt","host":"bbb.example.com","data":"v=spf1 include:sendgrid.net ~all"},"dkim":{"valid":false,"type":"txt","host":"m1._domainkey.aaa.example.com","data":"k=rsa; t=s; p=ccc"}}}`},
			want: &DomainAuthentication{
				ID:                 123,
				UserID:             234,
				Domain:             "example.com",
				Subdomain:          "aaa",
				Username:           "user",
				IPs:                []string{},
				CustomSPF:          false,
				IsDefault:          false,
				AutomaticSecurity:  false,
				CustomDKIMSelector: "",
				Legacy:             false,
				Valid:              false,
				DNS: DomainAuthenticationDNS{
					MailServer:   DomainAuthenticationDNSValue{Valid: false, Type: "mx", Host: "aaa.example.com", Data: "mx.sendgrid.net."},
					DKIM:         DomainAuthenticationDNSValue{Valid: false, Type: "txt", Host: "m1._domainkey.aaa.example.com", Data: "k=rsa; t=s; p=ccc"},
					SubDomainSPF: DomainAuthenticationDNSValue{Valid: false, Type: "txt", Host: "bbb.example.com", Data: "v=spf1 include:sendgrid.net ~all"},
				},
			},
			want1: RequestError{StatusCode: 200, Err: nil},
		},
		{
			name: "automatic security",
			args: args{respBody: `{"id":123,"user_id":234,"subdomain":"aaa","domain":"example.com","username":"user","ips":[],"custom_spf":false,"default":false,"legacy":false,"automatic_security":true,"valid":false,"dns":{"mail_cname":{"valid":false,"type":"cname","host":"aaa.example.com","data":"u234.abc.sendgrid.net"},"dkim1":{"valid":false,"type":"cname","host":"s1._domainkey.example.com","data":"s1.domainkey.u234.abc.sendgrid.net"},"dkim2":{"valid":false,"type":"cname","host":"s2._domainkey.example.com","data":"s2.domainkey.u234.abc.sendgrid.net"}}}`},
			want: &DomainAuthentication{
				ID:                 123,
				UserID:             234,
				Domain:             "example.com",
				Subdomain:          "aaa",
				Username:           "user",
				IPs:                []string{},
				CustomSPF:          false,
				IsDefault:          false,
				AutomaticSecurity:  true,
				CustomDKIMSelector: "",
				Legacy:             false,
				Valid:              false,
				DNS: DomainAuthenticationDNS{
					MailCNAME: DomainAuthenticationDNSValue{Valid: false, Type: "cname", Host: "aaa.example.com", Data: "u234.abc.sendgrid.net"},
					DKIM1:     DomainAuthenticationDNSValue{Valid: false, Type: "cname", Host: "s1._domainkey.example.com", Data: "s1.domainkey.u234.abc.sendgrid.net"},
					DKIM2:     DomainAuthenticationDNSValue{Valid: false, Type: "cname", Host: "s2._domainkey.example.com", Data: "s2.domainkey.u234.abc.sendgrid.net"},
				},
			},
			want1: RequestError{StatusCode: 200, Err: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseDomainAuthentication(tt.args.respBody)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDomainAuthentication() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseDomainAuthentication() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
