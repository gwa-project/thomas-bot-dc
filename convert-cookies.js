// Convert JSON cookies to Netscape format for yt-dlp
const fs = require('fs');

const cookies = [
    {
        "domain": ".youtube.com",
        "expirationDate": 1794131695.590516,
        "hostOnly": false,
        "httpOnly": true,
        "name": "__Secure-3PSID",
        "path": "/",
        "sameSite": "no_restriction",
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "g.a0002AjtiVU8b_KcZLNN_HhxMpmLu7dLDo6oZX7PDEoqO9OEKs6SuD5L0l1K1scXSCUP2zgn4QACgYKAd0SARMSFQHGX2MiDY3Nk0oLF3MfwJ2MKSK8YBoVAUF8yKqBMuUu_HFvje8KRA8QS12p0076"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1791531240.789767,
        "hostOnly": false,
        "httpOnly": true,
        "name": "__Secure-1PSIDTS",
        "path": "/",
        "sameSite": null,
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "sidts-CjQBmkD5SwciEUa8B0x9-x-EWQk--Tn8Rx6oY0_eYKeZGJK3zJCNOLUS_8aTkEUkBkgMvB5mEAA"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1794131695.590714,
        "hostOnly": false,
        "httpOnly": false,
        "name": "SAPISID",
        "path": "/",
        "sameSite": null,
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "pVeCye1cax0n1NLr/AuMt_AEfN4vywDt8R"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1791531283.519263,
        "hostOnly": false,
        "httpOnly": true,
        "name": "__Secure-1PSIDCC",
        "path": "/",
        "sameSite": null,
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "AKEyXzU-1-AoRxXA_BOYKCDfoYwDJ5ToyW6a-yVtnM8_7bBKXrZs3q7JcBeoHDs1tBXQAKnikAWF"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1794131695.590616,
        "hostOnly": false,
        "httpOnly": true,
        "name": "SSID",
        "path": "/",
        "sameSite": null,
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "AtYvhOnbwmWm2as2G"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1794131695.590765,
        "hostOnly": false,
        "httpOnly": false,
        "name": "__Secure-1PAPISID",
        "path": "/",
        "sameSite": null,
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "pVeCye1cax0n1NLr/AuMt_AEfN4vywDt8R"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1794131695.590464,
        "hostOnly": false,
        "httpOnly": true,
        "name": "__Secure-1PSID",
        "path": "/",
        "sameSite": null,
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "g.a0002AjtiVU8b_KcZLNN_HhxMpmLu7dLDo6oZX7PDEoqO9OEKs6STDBiaxzwx8hMz_XEOimTjAACgYKARcSARMSFQHGX2MiyxlhmryQSPjbPMLZfaPBmBoVAUF8yKrOlfzps0nCIAVThrpCOAEM0076"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1794131695.590816,
        "hostOnly": false,
        "httpOnly": false,
        "name": "__Secure-3PAPISID",
        "path": "/",
        "sameSite": "no_restriction",
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "pVeCye1cax0n1NLr/AuMt_AEfN4vywDt8R"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1791531283.519324,
        "hostOnly": false,
        "httpOnly": true,
        "name": "__Secure-3PSIDCC",
        "path": "/",
        "sameSite": "no_restriction",
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "AKEyXzUNoh4RVsLal2puF5IFP9JRim16cvwKZJdomPOgE_NnizVZy9EFEI0UNULESs0BmRGY7ac"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1791531240.789902,
        "hostOnly": false,
        "httpOnly": true,
        "name": "__Secure-3PSIDTS",
        "path": "/",
        "sameSite": "no_restriction",
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "sidts-CjQBmkD5SwciEUa8B0x9-x-EWQk--Tn8Rx6oY0_eYKeZGJK3zJCNOLUS_8aTkEUkBkgMvB5mEAA"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1794555272.002253,
        "hostOnly": false,
        "httpOnly": true,
        "name": "LOGIN_INFO",
        "path": "/",
        "sameSite": "no_restriction",
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "AFmmF2swRQIgPnytnmgNFjsfcfLuuW0J-LzvVgUK1qhWNdk7Hlg707cCIQD6cehmpZ1mhkDRFM3fKT3RtArzF1IHJ4Il2VVZYKV43g:QUQ3MjNmdzZvQmtnaG1qdFhrbV9obnl5dl92ZVBtUFdGWEpCeWw1SG1yQnU0V2xja1doMTF0SjI1TTM5aFg4akZZemJNaGRudzB5QkRaVHBhejV0aXJHckFveUF6dFk4M2RqTFE3SXAyN2tIalZtdkg0RzNhTHRDRV9NTDExT0VQYU1vQ0ZPYmZFZ0ZERlBuNkNxMTdhb3FOS2VZOHIzZkJR"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1775721012.784744,
        "hostOnly": false,
        "httpOnly": true,
        "name": "NID",
        "path": "/",
        "sameSite": null,
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "525=lJBlt0a0cFRYFv98Vb1y0C2ugr_Iqtf8wchLnkwhYhATEL-REHPo1ZK6gAFwQBYyTd24hKx9uxbDM2MzANBxuv9YdeuRzzt5QJArTBE7xkLgG8-oLVYY8S8QXUMjqYWdnVBtdvl8jG1BOjRWMRjZnKmJ46nZHv3dygVnucItFKcxFzTionCg9v_KJWxn-ABDQA--nQE-dZvEf_m3tWAbS8twVJtUeqWrklZd5LXjGC4-8LOaC6BBhW1LhFWaesFmfJ-G14VU25m1"
    },
    {
        "domain": ".youtube.com",
        "expirationDate": 1794555274.575499,
        "hostOnly": false,
        "httpOnly": false,
        "name": "PREF",
        "path": "/",
        "sameSite": null,
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "f6=40000000&volume=12&f7=100&tz=Asia.Bangkok&repeat=NONE&guide_collapsed=false&autoplay=true&f4=10000&f5=20000"
    }
];

// Netscape format
let output = "# Netscape HTTP Cookie File\n";

cookies.forEach(cookie => {
    const domain = cookie.domain;
    const flag = cookie.hostOnly ? "FALSE" : "TRUE";
    const path = cookie.path;
    const secure = cookie.secure ? "TRUE" : "FALSE";
    const expiration = Math.floor(cookie.expirationDate || 0);
    const name = cookie.name;
    const value = cookie.value;

    output += `${domain}\t${flag}\t${path}\t${secure}\t${expiration}\t${name}\t${value}\n`;
});

fs.writeFileSync('youtube_cookies.txt', output);
console.log('✅ Cookies converted to youtube_cookies.txt');
