<script lang="ts">
  import { onMount } from "svelte";
  import { checkAdmin, checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import type { User } from "../lib/model/user";

    let isDarkMode = true;
    let user: User | null = null;

    onMount(async() => {
        await checkAuth(navigate);
        await checkAdmin(navigate);

        const usr = await getAuthUser();
        if(usr) user = usr;

        const savedTheme = localStorage.getItem("theme");
        if (savedTheme) {
            isDarkMode = savedTheme === "dark";
            document.documentElement.setAttribute("data-theme", savedTheme);
        } else {
            // Check user's system preference if no saved theme
            const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
            isDarkMode = prefersDark;
            document.documentElement.setAttribute("data-theme", prefersDark ? "dark" : "light");
        }

        // await refreshData(1);
    });

</script>

<div>
    this is admin page
</div>

<style>

</style>