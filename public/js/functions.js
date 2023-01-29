function goToAuthPage() {
    window.location.replace("/views/login.html")
}

function sync() {
    Alpine.store('isSyncing', true)
    let jwt = Alpine.store('jwt')
    fetchHttp('/api/v1/sync/recent', jwt, "POST")
        .then(resp => {
            Alpine.store('lastSync', df(resp.ts))
            Alpine.store('isSyncing', false)
        })
}

async function fetchHttp(url, jwt, method="GET") {
    const response = await fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${jwt}`
        },
    });
    return response.json();
}

function short(item) {
    let l = item.length
    return `${item.substr(0, 4)}...${item.substr(l-4, l)}`
}

function price(item) {
    return `£${item}`
}

function df(d) {
    let date = new Date(d)
    let day = pad(date.getDate())
    let month = pad(date.getMonth() + 1)
    let year = pad(date.getFullYear())
    let hr = pad(date.getHours())
    let min = pad(date.getMinutes())
    let sec = pad(date.getSeconds())
    return `${day}/${month}/${year} @ ${hr}:${min}:${sec}`
}

function pad(n) {
    return n < 10 ? `0${n}` : `${n}`
}