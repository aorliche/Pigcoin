
$ = (selector) => document.querySelector(selector);
$$ = (selector) => document.querySelectorAll(selector);

addEventListener('load', () => {
    function createWallet(wallet) {
        const tr = document.createElement('tr');
        tr.innerHTML = `
        <td>${wallet.Id}</td>
        <td>${wallet.Name}</td>
        <td>${wallet.User}</td>
        <td>${wallet.Balance} PIG</td>
        `;
        return tr;
    }
    addEventListener('hashchange', (e) => {
        switch (location.hash) {
            case '#wallets': 
                $('#wallets').display ='inline-block'; 
                $('#transactions').display = 'none';
                break;
            case '#transactions':
                $('#transactions').display ='inline-block';
                $('#wallets').display = 'none';
                break;
        }
    });
    fetch('/api/wallets')
    .then(res => res.json())
    .then(wallets => {
        wallets.forEach(wallet => {
            $('#wallets tbody').appendChild(createWallet(wallet));
        });
    });
});
