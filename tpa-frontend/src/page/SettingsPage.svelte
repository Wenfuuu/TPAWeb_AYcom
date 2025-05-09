<script lang="ts">
  import { onMount } from "svelte";
  import { checkAuth } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";

    // let user: User | null = null;

    onMount(async() => {
        await checkAuth(navigate);

        // const usr = await getAuthUser();
        // if(usr) {
        //     user = usr;
        // }
    });

</script>

<div>
    this is settings page
</div>

<style>

</style>