const POST_URL = "http://localhost:8080/api/v1/projects/"

const start = () => {
    id('submit').addEventListener('click', handleClick)
}

const id = (id) => document.getElementById(id)

const handleClick = (event) => {
    event.preventDefault()
    const form = new FormData(id('form'))
    submit(form)
}

const submit = async (form) => {
    const response = await post(form)
    const responseData = await response.json()
    id('result').innerHTML = JSON.stringify(responseData).split(',').join('\n')
    console.log(form.getAll('description'))
}

const post = async (project) => {
    return fetch(POST_URL, {
        method: "POST",
        body: project,
    })
}

start()