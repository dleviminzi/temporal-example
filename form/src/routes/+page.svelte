<script>
	import { workflowData, pageData } from './state';

	let greetResult = '';
	let farewellResult = '';
	let fullWfResult = '';
	let queryResult = '';
	let runQueryResult = '';

	async function startWorkflow() {
		const res = await fetch('http://localhost:8090/startwf', {
			method: 'POST',
			body: JSON.stringify({
				id: $pageData.id
			})
		});
		const json = await res.json();
		console.log(json);
	}

	async function sendGreetInput() {
		$workflowData.greeted = true;
		const res = await fetch(
			'http://localhost:8090/greet?' +
				new URLSearchParams({
					id: String($pageData.id)
				}),
			{
				method: 'POST',
				body: JSON.stringify({
					name: $workflowData.greet,
					greet_done: true
				})
			}
		);

		const json = await res.json();
		greetResult = JSON.stringify(json);
		console.log($workflowData.greet);
	}

	async function sendFarewellInput() {
		$workflowData.goodbye = true;
		const res = await fetch(
			'http://localhost:8090/farewell?' +
				new URLSearchParams({
					id: String($pageData.id)
				}),
			{
				method: 'POST',
				body: JSON.stringify({
					farewell_done: $workflowData.goodbye
				})
			}
		);

		const json = await res.json();
		farewellResult = JSON.stringify(json);
		console.log($workflowData.farewell);
	}

	async function getFullWfResult() {
		const res = await fetch(
			'http://localhost:8090/result?' +
				new URLSearchParams({
					id: String($pageData.id)
				}),
			{
				method: 'GET'
			}
		);

		const json = await res.json();
		fullWfResult = JSON.stringify(json);
		console.log(fullWfResult);
	}

	async function getWorkflowRun() {
		const res = await fetch(
			'http://localhost:8090/activewf?' +
				new URLSearchParams({
					id: String($pageData.idRunQuery)
				}),
			{
				method: 'GET'
			}
		);

		const json = await res.json();
		runQueryResult = JSON.stringify(json);
		console.log(runQueryResult);
	}

	async function getWorkflows() {
		const res = await fetch(
			'http://localhost:8090/wfs?' +
				new URLSearchParams({
					id: String($pageData.idQuery)
				}),
			{
				method: 'GET'
			}
		);

		const json = await res.json();
		queryResult = JSON.stringify(json);
		console.log(queryResult);
	}

</script>

<form on:submit={startWorkflow} class="content">
	<label
		>enter user id greater than 0 to start new workflow for that id
		<input type="number" bind:value={$pageData.id} />
	</label>
	<input type="submit" value="start" />
</form>

<form on:submit={getWorkflowRun} class="content">
	<label
		>get latest active workflow's state for user id
		<input type="number" bind:value={$pageData.idRunQuery} />
	</label>
	<input type="submit" value="retrieve" />
</form>

<p>run query result</p>
{runQueryResult}

<p></p>
<form on:submit={getWorkflows} class="content">
	<label
		>view list of workflow runs for user id
		<input type="number" bind:value={$pageData.idQuery} />
	</label>
	<input type="submit" value="view" />
</form>

<p>query result</p>
{queryResult}

<h1>Form</h1>

<h3>child 1: greet</h3>
<form on:submit={sendGreetInput} class="content">
	<label
		>person to greet
		<input disabled={$workflowData.greeted == true} type="text" bind:value={$workflowData.greet} />
	</label>
	<input disabled={$workflowData.greeted == true} type="submit" value="submit" />
</form>

<p>greet result:</p>
{greetResult}

<h3>child 2: say farewell</h3>
<form on:submit={sendFarewellInput} class="content">
	say farewell? <input
		disabled={$workflowData.greeted != true || $workflowData.goodbye == true}
		type="submit"
		value="say farewell!"
	/>
</form>

<p>farewell result:</p>
{farewellResult}

<h3>parent workflow result</h3>
<form on:submit={getFullWfResult} class="content">
	<input
		type="submit"
		disabled={$workflowData.greeted != true || $workflowData.goodbye != true}
		value="get workflow result"
	/>
</form>
<p>workflow result:</p>
{fullWfResult}
