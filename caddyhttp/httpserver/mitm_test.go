package httpserver

import (
	"crypto/tls"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestParseClientHello(t *testing.T) {
	for i, test := range []struct {
		inputHex string
		expected rawHelloInfo
	}{
		{
			// curl 7.51.0 (x86_64-apple-darwin16.0) libcurl/7.51.0 SecureTransport zlib/1.2.8
			inputHex: `010000a6030358a28c73a71bdfc1f09dee13fecdc58805dcce42ac44254df548f14645f7dc2c00004400ffc02cc02bc024c023c00ac009c008c030c02fc028c027c014c013c012009f009e006b0067003900330016009d009c003d003c0035002f000a00af00ae008d008c008b01000039000a00080006001700180019000b00020100000d00120010040102010501060104030203050306030005000501000000000012000000170000`,
			expected: rawHelloInfo{
				cipherSuites:       []uint16{255, 49196, 49195, 49188, 49187, 49162, 49161, 49160, 49200, 49199, 49192, 49191, 49172, 49171, 49170, 159, 158, 107, 103, 57, 51, 22, 157, 156, 61, 60, 53, 47, 10, 175, 174, 141, 140, 139},
				extensions:         []uint16{10, 11, 13, 5, 18, 23},
				compressionMethods: []byte{0},
				curves:             []tls.CurveID{23, 24, 25},
				points:             []uint8{0},
			},
		},
		{
			// Chrome 56
			inputHex: `010000c003031dae75222dae1433a5a283ddcde8ddabaefbf16d84f250eee6fdff48cdfff8a00000201a1ac02bc02fc02cc030cca9cca8cc14cc13c013c014009c009d002f0035000a010000777a7a0000ff010001000000000e000c0000096c6f63616c686f73740017000000230000000d00140012040308040401050308050501080606010201000500050100000000001200000010000e000c02683208687474702f312e3175500000000b00020100000a000a0008aaaa001d001700182a2a000100`,
			expected: rawHelloInfo{
				cipherSuites:       []uint16{6682, 49195, 49199, 49196, 49200, 52393, 52392, 52244, 52243, 49171, 49172, 156, 157, 47, 53, 10},
				extensions:         []uint16{31354, 65281, 0, 23, 35, 13, 5, 18, 16, 30032, 11, 10, 10794},
				compressionMethods: []byte{0},
				curves:             []tls.CurveID{43690, 29, 23, 24},
				points:             []uint8{0},
			},
		},
		{
			// Firefox 51
			inputHex: `010000bd030375f9022fc3a6562467f3540d68013b2d0b961979de6129e944efe0b35531323500001ec02bc02fcca9cca8c02cc030c00ac009c013c01400330039002f0035000a010000760000000e000c0000096c6f63616c686f737400170000ff01000100000a000a0008001d001700180019000b00020100002300000010000e000c02683208687474702f312e31000500050100000000ff030000000d0020001e040305030603020308040805080604010501060102010402050206020202`,
			expected: rawHelloInfo{
				cipherSuites:       []uint16{49195, 49199, 52393, 52392, 49196, 49200, 49162, 49161, 49171, 49172, 51, 57, 47, 53, 10},
				extensions:         []uint16{0, 23, 65281, 10, 11, 35, 16, 5, 65283, 13},
				compressionMethods: []byte{0},
				curves:             []tls.CurveID{29, 23, 24, 25},
				points:             []uint8{0},
			},
		},
		{
			// openssl s_client (OpenSSL 0.9.8zh 14 Jan 2016)
			inputHex: `0100012b03035d385236b8ca7b7946fa0336f164e76bf821ed90e8de26d97cc677671b6f36380000acc030c02cc028c024c014c00a00a500a300a1009f006b006a0069006800390038003700360088008700860085c032c02ec02ac026c00fc005009d003d00350084c02fc02bc027c023c013c00900a400a200a0009e00670040003f003e0033003200310030009a0099009800970045004400430042c031c02dc029c025c00ec004009c003c002f009600410007c011c007c00cc00200050004c012c008001600130010000dc00dc003000a00ff0201000055000b000403000102000a001c001a00170019001c001b0018001a0016000e000d000b000c0009000a00230000000d0020001e060106020603050105020503040104020403030103020303020102020203000f000101`,
			expected: rawHelloInfo{
				cipherSuites:       []uint16{49200, 49196, 49192, 49188, 49172, 49162, 165, 163, 161, 159, 107, 106, 105, 104, 57, 56, 55, 54, 136, 135, 134, 133, 49202, 49198, 49194, 49190, 49167, 49157, 157, 61, 53, 132, 49199, 49195, 49191, 49187, 49171, 49161, 164, 162, 160, 158, 103, 64, 63, 62, 51, 50, 49, 48, 154, 153, 152, 151, 69, 68, 67, 66, 49201, 49197, 49193, 49189, 49166, 49156, 156, 60, 47, 150, 65, 7, 49169, 49159, 49164, 49154, 5, 4, 49170, 49160, 22, 19, 16, 13, 49165, 49155, 10, 255},
				extensions:         []uint16{11, 10, 35, 13, 15},
				compressionMethods: []byte{1, 0},
				curves:             []tls.CurveID{23, 25, 28, 27, 24, 26, 22, 14, 13, 11, 12, 9, 10},
				points:             []uint8{0, 1, 2},
			},
		},
	} {
		data, err := hex.DecodeString(test.inputHex)
		if err != nil {
			t.Fatalf("Test %d: Could not decode hex data: %v", i, err)
		}
		actual := parseRawClientHello(data)
		if !reflect.DeepEqual(test.expected, actual) {
			t.Errorf("Test %d: Expected %+v; got %+v", i, test.expected, actual)
		}
	}
}

func TestHeuristicFunctions(t *testing.T) {
	// To test the heuristics, we assemble a collection of real
	// ClientHello messages from various TLS clients. Please be
	// sure to hex-encode them and document the User-Agent
	// associated with the connection.
	//
	// If the TLS client used is not an HTTP client (e.g. s_client),
	// you can leave the userAgent blank, but please use a comment
	// to document crucial missing information such as client name,
	// version, and platform, maybe even the date you collected
	// the sample! Please group similar clients together, ordered
	// by version for convenience.

	// clientHello pairs a User-Agent string to its ClientHello message.
	type clientHello struct {
		userAgent string
		helloHex  string
	}

	// clientHellos groups samples of true (real) ClientHellos by the
	// name of the browser that produced them. We limit the set of
	// browsers to those we are programmed to protect, as well as a
	// category for "Other" which contains real ClientHello messages
	// from clients that we do not recognize, which may be used to
	// test or imitate interception scenarios.
	//
	// Please group similar clients and order by version for convenience
	// when adding to the test cases.
	clientHellos := map[string][]clientHello{
		"Chrome": []clientHello{
			{
				userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
				helloHex:  `010000c003031dae75222dae1433a5a283ddcde8ddabaefbf16d84f250eee6fdff48cdfff8a00000201a1ac02bc02fc02cc030cca9cca8cc14cc13c013c014009c009d002f0035000a010000777a7a0000ff010001000000000e000c0000096c6f63616c686f73740017000000230000000d00140012040308040401050308050501080606010201000500050100000000001200000010000e000c02683208687474702f312e3175500000000b00020100000a000a0008aaaa001d001700182a2a000100`,
			},
		},
		"Firefox": []clientHello{
			{
				userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:51.0) Gecko/20100101 Firefox/51.0",
				helloHex:  `010000bd030375f9022fc3a6562467f3540d68013b2d0b961979de6129e944efe0b35531323500001ec02bc02fcca9cca8c02cc030c00ac009c013c01400330039002f0035000a010000760000000e000c0000096c6f63616c686f737400170000ff01000100000a000a0008001d001700180019000b00020100002300000010000e000c02683208687474702f312e31000500050100000000ff030000000d0020001e040305030603020308040805080604010501060102010402050206020202`,
			},
		},
		// TODO... in the process of downloading a VM...
		// "Edge": []clientHello{
		// 	{
		// 		userAgent: "",
		// 		helloHex:  ``,
		// 	},
		// },
		"Safari": []clientHello{
			{
				userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8",
				helloHex:  `010000d2030358a295b513c8140c6ff880f4a8a73cc830ed2dab2c4f2068eb365228d828732e00002600ffc02cc02bc024c023c00ac009c030c02fc028c027c014c013009d009c003d003c0035002f010000830000000e000c0000096c6f63616c686f7374000a00080006001700180019000b00020100000d00120010040102010501060104030203050306033374000000100030002e0268320568322d31360568322d31350568322d313408737064792f332e3106737064792f3308687474702f312e310005000501000000000012000000170000`,
			},
		},
		"Other": []clientHello{
			{
				// openssl s_client (OpenSSL 0.9.8zh 14 Jan 2016)
				helloHex: `0100012b03035d385236b8ca7b7946fa0336f164e76bf821ed90e8de26d97cc677671b6f36380000acc030c02cc028c024c014c00a00a500a300a1009f006b006a0069006800390038003700360088008700860085c032c02ec02ac026c00fc005009d003d00350084c02fc02bc027c023c013c00900a400a200a0009e00670040003f003e0033003200310030009a0099009800970045004400430042c031c02dc029c025c00ec004009c003c002f009600410007c011c007c00cc00200050004c012c008001600130010000dc00dc003000a00ff0201000055000b000403000102000a001c001a00170019001c001b0018001a0016000e000d000b000c0009000a00230000000d0020001e060106020603050105020503040104020403030103020303020102020203000f000101`,
			},
			{
				// curl 7.51.0 (x86_64-apple-darwin16.0) libcurl/7.51.0 SecureTransport zlib/1.2.8
				userAgent: "curl/7.51.0",
				helloHex:  `010000a6030358a28c73a71bdfc1f09dee13fecdc58805dcce42ac44254df548f14645f7dc2c00004400ffc02cc02bc024c023c00ac009c008c030c02fc028c027c014c013c012009f009e006b0067003900330016009d009c003d003c0035002f000a00af00ae008d008c008b01000039000a00080006001700180019000b00020100000d00120010040102010501060104030203050306030005000501000000000012000000170000`,
			},
		},
	}

	for client, chs := range clientHellos {
		for i, ch := range chs {
			hello, err := hex.DecodeString(ch.helloHex)
			if err != nil {
				t.Errorf("[%s] Test %d: Error decoding ClientHello: %v", client, i, err)
				continue
			}
			parsed := parseRawClientHello(hello)

			isChrome := parsed.looksLikeChrome()
			isFirefox := parsed.looksLikeFirefox()
			isSafari := parsed.looksLikeSafari()
			isEdge := parsed.looksLikeEdge()

			// we want each of the heuristic functions to be as
			// exclusive but as low-maintenance as possible;
			// in other words, if one returns true, the others
			// should return false, with as little logic as possible,
			// but with enough logic to force TLS proxies to do a
			// good job preserving characterstics of the handshake.
			var wrong bool
			switch client {
			case "Chrome":
				wrong = !isChrome || isFirefox || isSafari || isEdge
			case "Firefox":
				wrong = isChrome || !isFirefox || isSafari || isEdge
			case "Safari":
				wrong = isChrome || isFirefox || !isSafari || isEdge
			case "Edge":
				wrong = isChrome || isFirefox || isSafari || !isEdge
			case "Others":
				wrong = isChrome || isFirefox || isSafari || isEdge
			}

			if wrong {
				t.Errorf("[%s] Test %d: Chrome=%v, Firefox=%v, Safari=%v, Edge=%v",
					client, i, isChrome, isFirefox, isSafari, isEdge)
			}
		}
	}
}
