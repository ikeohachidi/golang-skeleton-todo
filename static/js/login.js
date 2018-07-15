let login = document.getElementById("login")
let form = document.getElementsByTagName("form")[0]

login.addEventListener("click", function(e) {
    e.preventDefault()

    // let user = `username=${form.username.value}&password=${form.password.value}`
    let user = {
        username: form.username.value,
        password: form.password.value
    }
    
    let req = new XMLHttpRequest()
    req.open("POST", "http://localhost:8000/authenticate")
    req.setRequestHeader("Content-Type", "application/json")
    req.onload = function() {
        if (req.status >= 200 && req.status < 400) {
            let auth = JSON.parse(req.responseText)
            localStorage.setItem("username", auth.Username)
            localStorage.setItem("access_token", auth.Token)
            window.location.href = `/todo`
        }
        req.onerror = function() {
            console.log("An error has occured")
        }
    }
    req.send(JSON.stringify(user))
})