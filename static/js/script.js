let logout = document.getElementById("logout")

logout.addEventListener("click", (e) => {
    console.log("Hello")
    localStorage.clear()
    window.location.href = "/"
})