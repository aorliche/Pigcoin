
$ = (selector) => document.querySelector(selector);
$$ = (selector) => [...document.querySelectorAll(selector)];

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
    function createTransaction(transaction) {
        const tr = document.createElement('tr');
        tr.innerHTML = `
        <td>${transaction.Id}</td>
        <td>${transaction.From}</td>
        <td>${transaction.To}</td>
        <td>${transaction.Amount}</td>
        <td>${new Date(transaction.When*1000).toUTCString()}</td>
        `;
        return tr;
    }
    function getWallets() {
        $('#wallets tbody').innerHTML = '';
        fetch('/api/wallets')
        .then(res => res.json())
        .then(wallets => {
            wallets.forEach(wallet => {
                $('#wallets tbody').appendChild(createWallet(wallet));
            });
        });
    }
    function getTransactions() {
        $('#transactions tbody').innerHTML = '';
        fetch('/api/transactions')
        .then(res => res.json())
        .then(transactions => {
            transactions.forEach(transaction => {
                $('#transactions tbody').appendChild(createTransaction(transaction));
            });
        });
    }
    addEventListener('hashchange', (e) => {
        e.preventDefault();
        //console.log(location.hash);
        switch (location.hash) {
            case '#wallets': 
                getWallets();
                $('#wallets').style.display ='table'; 
                $('#transactions').style.display = 'none';
                break;
            case '#transactions':
                getTransactions();
                $('#transactions').style.display ='table';
                $('#wallets').style.display = 'none';
                break;
            case '#transfer':
                $('#transfer').style.display = 'block';
                $('#create').style.display = 'none';
                break;
            case '#create':
                $('#create').style.display = 'block';
                $('#transfer').style.display = 'none';
                break;
        }
    });
    getWallets();
    getTransactions();
    $('#transfer-pig').addEventListener('click', e => {
        e.preventDefault();
        fetch('/api/transfer', {
            body: new FormData($('#transfer')),
            method: 'post'
        })
        .then(res => res.json())
        .then(json => {
            if (json.Error) alert(json.Error);
            else {
                getTransactions();
                getWallets();
            }
        })
        .catch(err => alert(err));
    });
    $('#create-wallet').addEventListener('click', e => {
        e.preventDefault();
        fetch('/api/create-wallet', {
            body: new FormData($('#create')),
            method: 'post'
        })
        .then(res => res.json())
        .then(json => {
            if (json.Error) alert(json.Error);
            else getWallets();
        })
        .catch(err => alert(err));
    });
    function sortWallets(fieldNum) {
        const trs = $$('#wallets tbody tr');
        trs.forEach(tr => {
            tr.parentNode.removeChild(tr);
        });
        trs.sort((a, b) => {
            return parseInt(b.children[fieldNum].innerText) - parseInt(a.children[fieldNum].innerText); 
        });
        trs.forEach(tr => {
            $('#wallets tbody').appendChild(tr);
        });
    }
    $('#sort-balance').addEventListener('click', e => {
        e.preventDefault();
        sortWallets(3);
    });
    $('#sort-id').addEventListener('click', e => {
        e.preventDefault();
        sortWallets(0);
    });
});
