package fakedns

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xtls/xray-core/common"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/uuid"
)

func TestNewFakeDnsHolder(_ *testing.T) {
	_, err := NewFakeDNSHolder()
	common.Must(err)
}

func TestFakeDnsHolderCreateMapping(t *testing.T) {
	fkdns, err := NewFakeDNSHolder()
	common.Must(err)

	addr := fkdns.GetFakeIPForDomain("fakednstest.example.com")
	assert.Equal(t, "240.", addr[0].IP().String()[0:4])
}

func TestFakeDnsHolderCreateMappingMany(t *testing.T) {
	fkdns, err := NewFakeDNSHolder()
	common.Must(err)

	addr := fkdns.GetFakeIPForDomain("fakednstest.example.com")
	assert.Equal(t, "240.", addr[0].IP().String()[0:4])

	addr2 := fkdns.GetFakeIPForDomain("fakednstest2.example.com")
	assert.Equal(t, "240.", addr2[0].IP().String()[0:4])
	assert.NotEqual(t, addr[0].IP().String(), addr2[0].IP().String())
}

func TestFakeDnsHolderCreateMappingManyAndResolve(t *testing.T) {
	fkdns, err := NewFakeDNSHolder()
	common.Must(err)

	addr := fkdns.GetFakeIPForDomain("fakednstest.example.com")
	addr2 := fkdns.GetFakeIPForDomain("fakednstest2.example.com")

	{
		result := fkdns.GetDomainFromFakeDNS(addr[0])
		assert.Equal(t, "fakednstest.example.com", result)
	}

	{
		result := fkdns.GetDomainFromFakeDNS(addr2[0])
		assert.Equal(t, "fakednstest2.example.com", result)
	}
}

func TestFakeDnsHolderCreateMappingManySingleDomain(t *testing.T) {
	fkdns, err := NewFakeDNSHolder()
	common.Must(err)

	addr := fkdns.GetFakeIPForDomain("fakednstest.example.com")
	addr2 := fkdns.GetFakeIPForDomain("fakednstest.example.com")
	assert.Equal(t, addr[0].IP().String(), addr2[0].IP().String())
}

func TestFakeDnsHolderCreateMappingAndRollOver(t *testing.T) {
	fkdns, err := NewFakeDNSHolderConfigOnly(&FakeDnsPool{
		IpPool:  "240.0.0.0/12",
		LruSize: 256,
	})
	common.Must(err)

	err = fkdns.Start()

	common.Must(err)

	addr := fkdns.GetFakeIPForDomain("fakednstest.example.com")
	addr2 := fkdns.GetFakeIPForDomain("fakednstest2.example.com")

	for i := 0; i <= 8192; i++ {
		{
			result := fkdns.GetDomainFromFakeDNS(addr[0])
			assert.Equal(t, "fakednstest.example.com", result)
		}

		{
			result := fkdns.GetDomainFromFakeDNS(addr2[0])
			assert.Equal(t, "fakednstest2.example.com", result)
		}

		{
			uuid := uuid.New()
			domain := uuid.String() + ".fakednstest.example.com"
			tempAddr := fkdns.GetFakeIPForDomain(domain)
			rsaddr := tempAddr[0].IP().String()

			result := fkdns.GetDomainFromFakeDNS(net.ParseAddress(rsaddr))
			assert.Equal(t, domain, result)
		}
	}
}
