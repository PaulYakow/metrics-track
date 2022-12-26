package pki

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	// Данные для шифрования
	data := []byte(`
Etiam mollis congue massa vitae iaculis. Maecenas arcu elit, eleifend nec urna pellentesque, hendrerit suscipit massa. Mauris lobortis odio et pulvinar interdum. Nunc ac ex eu lacus semper placerat. Suspendisse potenti. Sed et varius risus. Fusce ut dapibus metus. Nulla ac eleifend dui. Fusce ut urna nec risus fringilla placerat. Curabitur varius ut libero id sollicitudin. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Suspendisse feugiat pretium justo, non bibendum mi posuere at. Vivamus semper dui a efficitur imperdiet. Cras gravida ante mattis sodales molestie. Maecenas ornare elementum tortor, at euismod justo viverra in. Etiam congue scelerisque finibus.
Aliquam et ex tincidunt, maximus tortor quis, aliquet felis. Suspendisse posuere vitae urna in volutpat. Etiam vitae sollicitudin dui. Phasellus vulputate erat vel metus malesuada mattis at sed ex. In et dolor vel justo lobortis pellentesque. Interdum et malesuada fames ac ante ipsum primis in faucibus. Ut non dolor massa. Proin arcu enim, pharetra id tortor in, pharetra venenatis elit. Pellentesque dapibus tristique dui, in hendrerit elit molestie a. Integer sodales molestie erat, eget feugiat justo pharetra sit amet. Integer aliquam imperdiet gravida. Vestibulum tempor commodo velit, id dictum lacus porta vel. Etiam a ligula in mauris tincidunt mattis. Proin eget pellentesque sapien. Maecenas pretium mauris at faucibus ultricies.
Phasellus auctor orci porta pharetra lobortis. Nunc condimentum sagittis lacus, auctor pellentesque eros suscipit et. Aliquam vitae libero sollicitudin, rutrum magna in, semper ipsum. Curabitur sit amet elit non augue vulputate eleifend. Vivamus semper lorem et lectus ultricies consequat. Donec facilisis, neque sit amet congue luctus, tortor nisl hendrerit purus, ac tincidunt leo augue sit amet leo. Nulla tristique metus dolor, ac tempor diam cursus nec. Etiam at ligula sed eros venenatis malesuada nec in velit. Duis dignissim, augue vel fermentum euismod, diam purus mattis erat, fringilla blandit magna leo eget massa. Curabitur ultrices purus in massa iaculis gravida. Mauris orci enim, viverra et placerat ultrices, varius id magna. Curabitur nec lectus aliquam, bibendum ex quis, euismod mi. Morbi et elit fermentum, convallis enim a, sagittis urna. Quisque bibendum eu urna a porta. In fringilla risus quis cursus scelerisque. Maecenas vel erat lectus.
Cras gravida, turpis vitae hendrerit viverra, purus orci imperdiet nibh, vitae cursus dui dolor non ipsum. Integer tempor arcu vel diam laoreet, a congue metus sagittis. Proin volutpat sollicitudin nisi ullamcorper congue. Aenean viverra leo a viverra molestie. Integer dui eros, sagittis sed ex id, dictum rutrum ipsum. Sed at massa at neque rhoncus porta. Nunc sed metus varius, fermentum leo vitae, convallis ex. In quis velit in tortor semper blandit. Morbi ultricies non felis in dapibus. Vivamus eu orci id neque gravida ullamcorper nec eu urna. Sed dignissim sapien eros, eget rutrum ligula laoreet vitae. Aenean volutpat semper tincidunt. Etiam quis augue in risus finibus porta. Praesent ac scelerisque sem, placerat suscipit elit.
Vestibulum quis lorem tristique, congue nisl eget, facilisis nulla. Maecenas non tellus quis libero congue condimentum. Suspendisse potenti. Sed commodo ultrices libero. Suspendisse elementum dolor vitae nisi placerat faucibus. Fusce sed feugiat odio. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Fusce gravida, tellus in gravida hendrerit, lectus nunc eleifend felis, sed gravida mauris ante sed mi. Phasellus vel orci nunc. Sed est justo, accumsan et velit in, efficitur semper eros.
Duis laoreet ante velit, vitae finibus nisl semper id. Maecenas volutpat purus eget aliquam tincidunt. Morbi at neque ut arcu dictum lobortis. Ut ut leo convallis, ultrices purus finibus, ultrices diam. Vivamus id euismod eros. In hac habitasse platea dictumst. Nulla molestie dictum placerat. Quisque ac luctus lorem.
Cras quis mauris fringilla, tempus nisi quis, lobortis risus. Nullam diam elit, varius vitae volutpat quis, efficitur vitae lectus. Sed rutrum purus augue, a dapibus justo viverra non. Fusce sodales urna in quam mattis tempus. Suspendisse fringilla orci turpis, quis euismod est eleifend efficitur. In dignissim velit sit amet eros fringilla, a tempus nunc semper. Aliquam erat volutpat. Mauris tempus est urna, a sagittis quam suscipit sit amet. Nullam suscipit ullamcorper eros, auctor elementum lorem aliquam sed.
Aliquam sapien metus, pulvinar ac tellus at, euismod aliquet eros. Donec nec sapien ut leo mollis mollis. Phasellus nec porta turpis. Ut egestas, dolor ut accumsan dictum, massa leo volutpat ipsum, in tincidunt diam nisi a quam. Mauris at nulla ut metus sollicitudin ultricies. Proin lorem lacus, ornare at metus quis, tempus aliquet nisl. Proin in arcu sapien. Donec justo dolor, dignissim eget tincidunt nec, condimentum vel lorem. Phasellus vel blandit ipsum. Ut iaculis mattis justo at pharetra. Maecenas a porta massa, vitae mollis odio.
Morbi maximus quis erat et venenatis. Sed et ante dignissim, maximus dolor sit amet, dictum massa. In eleifend nisl orci, vel venenatis lacus sollicitudin ut. Duis fermentum quam nisl, ac fermentum augue finibus semper. Suspendisse vel maximus nibh, ut lobortis leo. Sed id congue ligula. Suspendisse quis ex quis nisl mollis commodo vel a justo. Sed eu leo vel orci ultrices congue. Proin dolor ante, suscipit in libero quis, rhoncus consectetur nunc. Quisque in neque tincidunt, lobortis arcu eget, semper orci. Duis non egestas orci. Nam venenatis dui vel est varius, nec vehicula sapien faucibus. Quisque velit augue, hendrerit et massa dignissim, posuere dictum diam. Morbi nisl nibh, faucibus vitae tincidunt tempus, lobortis ut odio. Donec efficitur lacus sit amet ligula ultricies mattis. Nulla eget enim lorem.
Pellentesque vitae massa a risus maximus mollis. Aliquam lobortis est augue. In laoreet ipsum a ullamcorper condimentum. Donec ultricies dapibus augue. Nunc lectus nunc, lobortis vel consequat pulvinar, hendrerit eu ligula. Nam porta leo nisi, id eleifend nisl imperdiet et. Sed id lobortis diam, id vehicula est. Nam sed malesuada dolor. Cras ultricies fermentum mauris eu egestas. Donec tincidunt ante quis odio hendrerit ultricies. Praesent suscipit quam eu neque placerat volutpat. Nulla auctor tempor lacus, eu condimentum nisi sollicitudin quis. Maecenas blandit diam vel augue molestie, quis auctor ipsum sodales.
Mauris sit amet orci id nisl pharetra semper sit amet feugiat nisl. Quisque nunc ante, mollis quis lorem at, feugiat scelerisque arcu. Maecenas id iaculis quam. Integer commodo diam porttitor lacus hendrerit, et varius risus feugiat. Suspendisse sodales dolor nec libero eleifend consectetur in non nulla. Maecenas ac bibendum felis. Vivamus quam elit, venenatis a turpis gravida, sagittis convallis leo. Suspendisse posuere tincidunt erat, id sollicitudin nisl tincidunt a. Mauris tristique pharetra nibh nec imperdiet. Vivamus fringilla pellentesque massa eget pulvinar. Etiam et ex quam. Nunc ut nulla in tellus aliquam gravida sodales sit amet metus.
Curabitur ipsum felis, bibendum non ultrices nec, pharetra quis nibh. Phasellus quis porta metus, non tempus lacus. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Curabitur at sem at nisl consectetur luctus ac at libero. Donec nunc ligula, vestibulum id varius ut, bibendum sit amet eros. Sed maximus mauris vitae turpis egestas facilisis. Vestibulum placerat est at tincidunt egestas. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Maecenas posuere sem urna, vel tristique sapien tincidunt a. Sed non lorem scelerisque est aliquet tincidunt quis id nisl. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.
Morbi ornare lectus eget tempor fringilla. Maecenas a aliquet nunc. Mauris varius ullamcorper libero, sit amet fermentum nisi lobortis a. Proin a dapibus nisl, vitae vestibulum est. Etiam rutrum, dui a suscipit imperdiet, dolor lacus tincidunt nibh, in ullamcorper lacus tortor in lorem. Lorem ipsum dolor sit est.
`)

	// Имитация конструктора NewCryptographer ("читаем" публичный ключ)
	c := &Cryptographer{dataReader: func() ([]byte, error) {
		return []byte(`-----BEGIN PUBLIC KEY-----
MIICCgKCAgEAoGcVIhZgpX7IvB7TXdqruS7Hko7z+1V/wP1p2SrYxUDUstsHsyQ/
GdvjejsTO4G6CQ8irn35e931aOm0+d1brT3ELSd1rYD5nSf4jOyB51GlQJEUqwEN
V5TBuEm4xuZsDCJPv4ibdfp8AnQJctJxrleg8z0tAPLwbaE3mUko+LThTnocx97n
S/cvH0X4IjXUgBsT9xqwgbAOr+L9zMphU01R8rSlP5kXXVwedtoW5IOlpuT2HWXk
flTbag+ARDK3zpllqYF3TBthmDVuE3uqOVbGJtIq4gRU4NOdl3oXNkmEmfc3/v5E
hI/iMtzdAQ7akB1F/TTVKu4wclmzpLr0bo5f6ivz4IOIeHi4iEVbUyzqUt8C0tfv
vlgu43Eu581g4ptgKYiHHYavJEhQRnZFsN/3hDRmHPr6K34UZEiN8SwOjVOWnjgA
VOuwRz3/TcWecR7it2mtngBg0FjixgX73LMyGc0EaQjEPsqHGB/5y1ZjPj2QiS7q
8XTjAzEIlpoPKc+Th8x9kTaq69IyDVRdsKhkK8B4QzW74RORGn5XDXaWxfs3LMq4
lcvLTe52ORTtiVOWJMcKYWXtNpGkbj+b0/nt5auzt+AHKCd2d2vsR6HkRkZCzHRa
IP9BSxObYE7mr+WKDDNKu7bDnNnzo7tpCnhYQ6WpGBBNRoIzO+XsVYMCAwEAAQ==
-----END PUBLIC KEY-----`), nil
	}}

	publicKeyBytes, err := c.dataReader()
	require.NoError(t, err)
	require.NotNil(t, publicKeyBytes)

	publicKey, err := convertBytesToPublicKey(publicKeyBytes)
	require.NoError(t, err)
	require.NotNil(t, publicKey)

	c.publicKey = publicKey

	// Проверка шифрования
	dataEncrypted, err := c.Encrypt(data)
	require.NoError(t, err)
	require.NotNil(t, dataEncrypted)

	// Имитация конструктора NewDecryptor ("читаем" приватный ключ)
	d := &Decryptor{dataReader: func() ([]byte, error) {
		return []byte(`-----BEGIN PRIVATE KEY-----
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
-----END PRIVATE KEY-----`), err
	}}

	privateKeyBytes, err := d.dataReader()
	require.NoError(t, err)
	require.NotNil(t, privateKeyBytes)

	privateKey, err := convertBytesToPrivateKey(privateKeyBytes)
	require.NoError(t, err)
	require.NotEmpty(t, privateKey)

	d.privateKey = privateKey

	// Проверка дешифрования
	dataDecrypted, err := d.Decrypt(dataEncrypted)
	require.NoError(t, err)
	require.Equal(t, data, dataDecrypted)

	// Некорректный путь до файла ключа
	c, err = NewCryptographer("/wrong/path/to/public.key")
	require.Error(t, err)
	require.Nil(t, c)

	d, err = NewDecryptor("/wrong/path/to/private.key")
	require.Error(t, err)
	require.Nil(t, d)

}
