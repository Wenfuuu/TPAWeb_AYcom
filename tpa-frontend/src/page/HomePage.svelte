<script>
    import { onMount } from "svelte";
    import { navigate } from "svelte-routing";
    import { checkAuth } from "../lib/utils/auth";


    onMount(async() => {
        await checkAuth(navigate);
        // await checkAdmin(navigate);
    });

</script>

<h1>This is Home Page</h1>

<style>

</style>