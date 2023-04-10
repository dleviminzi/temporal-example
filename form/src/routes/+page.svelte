<script>
    import { input } from './state'

    function handleGreet() {
        $input.greeted = true
        console.log($input.greet)
    }

    function handleFarewell() {
        console.log($input.farewell)
    }

</script>
	
<h1>Form</h1>
<form on:submit={handleGreet} class="content">
  <label>person to greet
    <input disabled={$input.greeted == true} type="text" bind:value={$input.greet} />
  </label>
  <input disabled={$input.greeted == true} type="submit" value="submit"/>
</form>

<form on:submit={handleFarewell} class="content">
  <label>person to say farewell to
    <input disabled={$input.greeted != true} type="text" bind:value={$input.farewell} />
  </label>
  <input disabled={$input.greeted != true} type="submit" value="submit"/>
</form>
