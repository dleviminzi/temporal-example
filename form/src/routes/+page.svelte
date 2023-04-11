<script>
    import { input } from './state'

    let resultGreet = ""
    let resultFarewell = ""
    let result = ""

    async function handleGreet() {
        $input.greeted = true
		const res = await fetch('http://localhost:8090/greet', {
			method: 'POST',
			body: JSON.stringify({
                "name": $input.greet,
                "greet_done": true
            })
		})
		
		const json = await res.json()
		resultGreet = JSON.stringify(json)
        console.log($input.greet)
    }

    async function handleFarewell() {
        $input.goodbye = true
		const res = await fetch('http://localhost:8090/farewell', {
			method: 'POST',
			body: JSON.stringify({
                "farewell_done": $input.goodbye,
            })
		})
		
		const json = await res.json()
		resultFarewell = JSON.stringify(json)
        console.log($input.farewell)
    }

    async function getResult() {
        console.log("called")
		const res = await fetch('http://localhost:8090/result', {
			method: 'GET',
		})

		const json = await res.json()
		result = JSON.stringify(json)
        console.log(result)
    }

</script>
	
<h1>Form</h1>

<h3> child 1: greet </h3>
<form on:submit={handleGreet} class="content">
  <label>person to greet
    <input disabled={$input.greeted == true} type="text" bind:value={$input.greet} />
  </label>
  <input disabled={$input.greeted == true} type="submit" value="submit"/>
</form>

<p>
greet result:
</p>
{resultGreet}

<h3> child 2: say farewell </h3>
<form on:submit={handleFarewell} class="content">
  say farewell? <input disabled={$input.greeted != true || $input.goodbye == true} type="submit" value="say farewell!"/>
</form>

<p>
farewell result:
</p>
{resultFarewell}

<h3> parent workflow result </h3>
<form on:submit={getResult} class="content">
<input type="submit" disabled={$input.greeted != true || $input.goodbye != true} value="get workflow result"/>
</form>
<p>
workflow result:
</p>
{result}
