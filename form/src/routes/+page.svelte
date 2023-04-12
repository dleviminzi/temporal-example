<script>
    import { workflowData, pageData } from './state'

    let resultGreet = ""
    let resultFarewell = ""
    let result = ""

    async function handleGreet() {
        $workflowData.greeted = true
		const res = await fetch('http://localhost:8090/greet?'+ new URLSearchParams({
                id: String($pageData.id),
            }), {
			method: 'POST',
			body: JSON.stringify({
                "name": $workflowData.greet,
                "greet_done": true
            })
		})
		
		const json = await res.json()
		resultGreet = JSON.stringify(json)
        console.log($workflowData.greet)
    }

    async function handleFarewell() {
        $workflowData.goodbye = true
		const res = await fetch('http://localhost:8090/farewell?'+ new URLSearchParams({
                id: String($pageData.id),
            }), {
			method: 'POST',
			body: JSON.stringify({
                "farewell_done": $workflowData.goodbye,
            })
		})
		
		const json = await res.json()
		resultFarewell = JSON.stringify(json)
        console.log($workflowData.farewell)
    }

    async function getResult() {
		const res = await fetch('http://localhost:8090/result?' + new URLSearchParams({
                id: String($pageData.id),
            }), {
			method: 'GET',
		})

		const json = await res.json()
		result = JSON.stringify(json)
        console.log(result)
    }

    async function startWorkflow() {
		const res = await fetch('http://localhost:8090/startwf', {
			method: 'POST',
			body: JSON.stringify({
                "id": $pageData.id,
            })
		})
		const json = await res.json()
        console.log(json)
    }

    async function getWorkflowRun() {
		const res = await fetch('http://localhost:8090/activewf?' + new URLSearchParams({
                id: String($pageData.id),
            }), {
			method: 'GET',
		})

		const json = await res.json()
		let r = JSON.stringify(json)
        console.log(r)
    }

    async function getWorkflows() {
		const res = await fetch('http://localhost:8090/wfs?' + new URLSearchParams({
                id: String($pageData.wfQuery),
            }), {
			method: 'GET',
		})

		const json = await res.json()
		result = JSON.stringify(json)
        console.log(result)
    }

</script>

<form on:submit={startWorkflow} class="content">
  <label>enter user id greater than 0 to start new workflow for that id
    <input type="number" bind:value={$pageData.id} />
  </label>
  <input type="submit" value="start"/>
</form>

<form on:submit={getWorkflowRun} class="content">
  <label>get latest active workflow for user id
    <input type="number" bind:value={$pageData.id} />
  </label>
  <input type="submit" value="retrieve"/>
</form>


<form on:submit={getWorkflows} class="content">
  <label>view list of workflow runs for user id
    <input type="number" bind:value={$pageData.wfQuery} />
  </label>
  <input type="submit" value="view"/>
</form>

	
<h1>Form</h1>

<h3> child 1: greet </h3>
<form on:submit={handleGreet} class="content">
  <label>person to greet
    <input disabled={$workflowData.greeted == true} type="text" bind:value={$workflowData.greet} />
  </label>
  <input disabled={$workflowData.greeted == true} type="submit" value="submit"/>
</form>

<p>
greet result:
</p>
{resultGreet}

<h3> child 2: say farewell </h3>
<form on:submit={handleFarewell} class="content">
  say farewell? <input disabled={$workflowData.greeted != true || $workflowData.goodbye == true} type="submit" value="say farewell!"/>
</form>

<p>
farewell result:
</p>
{resultFarewell}

<h3> parent workflow result </h3>
<form on:submit={getResult} class="content">
<input type="submit" disabled={$workflowData.greeted != true || $workflowData.goodbye != true} value="get workflow result"/>
</form>
<p>
workflow result:
</p>
{result}
