<script lang="ts">
  import { onMount } from "svelte";
  import { checkGuest } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";


    let isDarkMode = true;

    onMount(async() => {
        await checkGuest(navigate);

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
    this is forgotten page
</div>

<style>

</style>