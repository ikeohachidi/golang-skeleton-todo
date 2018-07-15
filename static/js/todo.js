let tasks = document.getElementsByClassName("task")
let addTodo = document.getElementById("add_todo")
let listTodo = document.getElementById("list_todo")
let todoField = document.getElementById("add_todo_field")
let authTodo = document.getElementById("auth_todo")

let tasksList = []

let user = {
    username: localStorage.getItem("username"),
    token: localStorage.getItem("access_token")
}

// Fetch all of the users todos
let req = new XMLHttpRequest()
req.open("POST", `http://localhost:8000/${user.username}/gettodo`)
req.onload = function () {
    if (req.status >= 200 && req.status < 400) {
        let data = JSON.parse(req.responseText)
        console.log(data)

        for (let i = 0; i < data.length; i++) {
            let li = document.createElement("li")

            // deleteBtn is the button that gets clicked to remove a todo
            let deleteBtn = document.createElement("button")
            let deleteBtnContent = document.createTextNode("X")
            deleteBtn.appendChild(deleteBtnContent)
            deleteBtn.setAttribute("class", "button-x remove_todo")
            deleteTodo(deleteBtn)

            // taskElContent is going to be the todo gotten from the database which 
            // in turn is going to be the text content of the li tag
            let liContent = document.createTextNode(data[i])

            li.appendChild(deleteBtn)
            li.appendChild(liContent)
            li.setAttribute("class", "task")

            listTodo.appendChild(li)
        }

    }
    req.onerror = function () {
        console.log("An error occured")
    }
}
req.send(JSON.stringify(user))


// addTodo listens for a click event 
//  then adds the value of the input string to the DOM and saves it to the database
addTodo.addEventListener("click", (e) => {
    addTask = todoField.value
    if (addTask != "") {
        let li = document.createElement("li")

        // deleteBtn is the button that gets clicked to remove a todo
        let deleteBtn = document.createElement("button")
        let deleteBtnContent = document.createTextNode("X")
        deleteBtn.appendChild(deleteBtnContent)
        deleteBtn.setAttribute("class", "button-x remove_todo")
        deleteTodo(deleteBtn)

        let taskElContent = document.createTextNode(addTask)

        li.appendChild(deleteBtn)
        li.appendChild(taskElContent)
        li.setAttribute("class", "task")

        listTodo.appendChild(li)

        for (let x of tasks) {
            if (x.textContent == addTask) {
                tasksList.push(x.textContent)
            } else {
                tasksList.push(x.textContent.substring(1, ))
            }
        }
        let req = new XMLHttpRequest()
        req.open("POST", `http://localhost:8000/updatetodo/${tasksList}`)
        req.onload = function () {
            if (req.status >= 200 && req.status < 400) {
                data = JSON.parse(req.responseText)
                console.log(data)
            }
            req.onerror = function () {
                console.log("An error Occured")
            }
        }
        req.send(JSON.stringify(user))
        tasksList = []
        todoField.value = ""
    }
})

// deleteTodo remove the specific tasks and then updates the database
function deleteTodo(btn) {
    btn.addEventListener("click", (e) => {
        let li = e.target.parentElement
        li.parentElement.removeChild(li)
        for (x of tasks) {
            tasksList.push(x.textContent.substring(1 , ))
        }
        let req = new XMLHttpRequest()
        req.open("POST", `http://localhost:8000/updatetodo/${tasksList}`)
        req.onload = function () {
            if (req.status >=  200 && req.status < 400) {
                data = JSON.parse(req.responseText)
                console.log(data)
            }
            req.onerror = function () {
                console.log("An error has occured")
            }
        }
        req.send(JSON.stringify(user))
        tasksList = []
    })
}

if (user.username == null && user.token == null) {
    authTodo.style.display = "block";
    todoField.style.display = "none"
    addTodo.textContent = "Really Just sign in"
}