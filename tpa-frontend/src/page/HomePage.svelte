<script>
    import { onMount } from "svelte";
    import { navigate } from "svelte-routing";
    import { checkAdmin, checkAuth } from "../lib/utils/auth";


    onMount(async() => {
        checkAuth(navigate);
        // checkAdmin(navigate);
    });

</script>

<h1>This is Home Page</h1>

<style>

</style>