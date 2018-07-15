let logout = document.getElementById("logout")
let menuBtn = document.getElementById("menu_btn")
let menu = document.getElementById("menu")

logout.addEventListener("click", (e) => {
    localStorage.clear()
    window.location.href = "/"
})

menuBtn.addEventListener("click", (e) => {
    menu.classList.toggle("reveal")
    menu.classList.add("u-full-width")
    for (let x of menu.children) {
        x.style.textAlign = "right"
        x.classList.add("u-full-width")
    }
})