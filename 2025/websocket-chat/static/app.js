let ws = null
let username = ''

const $chat = document.getElementById('chat')
const $username = document.getElementById('username')
const $loginBtn = document.getElementById('loginBtn')
const $logoutBtn = document.getElementById('logoutBtn')
const $content = document.getElementById('content')
const $sendBtn = document.getElementById('sendBtn')

function appendMessage(text, cls) {
    const d = document.createElement('div')
    d.className = 'msg ' + (cls || '')
    d.textContent = text
    $chat.appendChild(d)
    $chat.scrollTop = $chat.scrollHeight
}

$loginBtn.addEventListener('click', () => {
    if (ws && ws.readyState === WebSocket.OPEN) {
        appendMessage('已经连接', 'system')
        return
    }
    username = $username.value.trim()
    if (!username) {
        alert('请输入用户名')
        return
    }

    ws = new WebSocket("ws://" + location.host + "/ws")

    ws.onopen = () => {
        appendMessage('连接已建立', 'system')
        ws.send(JSON.stringify({type: 'event', from: username, content: 'login'}))
        $loginBtn.disabled = true
        $logoutBtn.disabled = false
        $sendBtn.disabled = false
    }

    ws.onmessage = (evt) => {
        const msg = JSON.parse(evt.data)
        const t = new Date(msg.time * 1000).toLocaleTimeString()
        if (msg.type === 'system') {
            appendMessage(`[系统 ${t}] ${msg.content}`, 'system')
        } else {
            const who = msg.from === username ? '我' : msg.from
            appendMessage(`[${t}] ${who}: ${msg.content}`, msg.from === username ? 'me' : '')
        }
    }

    ws.onclose = () => {
        appendMessage('连接已关闭', 'system')
        $loginBtn.disabled = false
        $logoutBtn.disabled = true
        $sendBtn.disabled = true
    }

    ws.onerror = () => appendMessage('连接错误', 'system')
})

$logoutBtn.addEventListener('click', () => {
    if (ws) ws.close()
    ws = null
    appendMessage('已登出', 'system')
    $loginBtn.disabled = false
    $logoutBtn.disabled = true
    $sendBtn.disabled = true
})

$sendBtn.addEventListener('click', () => {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        alert('未连接')
        return
    }
    const text = $content.value.trim()
    if (!text) return
    ws.send(JSON.stringify({type: 'text', from: username, content: text}))
    $content.value = ''
})

$content.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') $sendBtn.click()
})
