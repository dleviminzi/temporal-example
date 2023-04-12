import { writable } from 'svelte/store'

export const workflowData = writable({
    greet: "daniel",
    greeted: false,
    farewell: "daniel",
    goodbye: false,
})

export const pageData = writable({
    id: 0,
    idQuery: 0,
    idRunQuery: 0
})

