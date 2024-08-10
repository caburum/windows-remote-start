first tutorial: https://github.com/asheroto/WinRM-HTTPS-NonDomain-Computers/blob/master/Target-Machine.ps1

tutorial that worked to set up ca certificate (`ca.crt`): https://github.com/jijiechen/winrm-client-certificate-auth

```ps
winrm set winrm/config/service/auth '@{Basic="false"}'
winrm set winrm/config/service '@{AllowUnencrypted="false"}'
```

tutorial that worked to set up client certificates (`user.pfx`, `user.pem`, `key.pem`): https://www.hurryupandwait.io/blog/certificate-password-less-based-authentication-in-winrm

some code to generate a certificate:

```ps
function New-ClientCertificate {
	param([String]$username, [String]$basePath = ((Resolve-Path .).Path))

	$OPENSSL_CONF=[System.IO.Path]::GetTempFileName()

	Set-Content -Path $OPENSSL_CONF -Value @"
	distinguished_name = req_distinguished_name
	[req_distinguished_name]
	[v3_req_client]
	extendedKeyUsage = clientAuth
	subjectAltName = otherName:1.3.6.1.4.1.311.20.2.3;UTF8:$username@localhost
"@

	$env:OPENSSL_CONF=$OPENSSL_CONF;

	$user_path = Join-Path $basePath user.pem
	$key_path = Join-Path $basePath key.pem
	$pfx_path = Join-Path $basePath user.pfx

	& 'C:\Program Files\OpenSSL-Win64\bin\openssl.exe' req -x509 -nodes -days 3650 -newkey rsa:2048 -out $user_path -outform PEM -keyout $key_path -subj "/CN=$username" -extensions v3_req_client 2>&1

	& 'C:\Program Files\OpenSSL-Win64\bin\openssl.exe' pkcs12 -export -in $user_path -inkey $key_path -out $pfx_path -passout pass: 2>&1

	del $OPENSSL_CONF
}
```
