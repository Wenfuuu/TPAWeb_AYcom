<script lang="ts">
  import { onMount } from "svelte";
  import { checkAuth } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import LeftSidebar from "../lib/components/LeftSidebar.svelte";
  import RightSidebar from "../lib/components/RightSidebar.svelte";
  import HomeContent from "../lib/components/HomeContent.svelte";

  let isDarkMode = true; // Default to dark mode

  // Initialize theme from localStorage on component mount
  onMount(async() => {
    await checkAuth(navigate);

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
  });

  // Toggle theme function
  function toggleTheme() {
    isDarkMode = !isDarkMode;
    const newTheme = isDarkMode ? "dark" : "light";
    localStorage.setItem("theme", newTheme);
    document.documentElement.setAttribute("data-theme", newTheme);
  }
</script>

<div class="home-container">
    <!-- <h1>nigga</h1> -->
  <LeftSidebar {toggleTheme} {isDarkMode} />
  <HomeContent />
  <RightSidebar />
</div>

<style>
  .home-container {
    display: flex;
    width: 100%;
    height: 100vh;
    background-color: var(--background);
    color: var(--text-primary);
  }

  /* @media (max-width: 1024px) {
    .home-container {
      justify-content: center;
    }
  } */

  /* @media (max-width: 768px) {
    .home-container {
      flex-direction: column;
    }
  } */
</style>