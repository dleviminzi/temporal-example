import { writable } from 'svelte/store'

const input = writable({
  greet: "daniel",
  greeted: false,
  farewell: "daniel",
})

export {input}
