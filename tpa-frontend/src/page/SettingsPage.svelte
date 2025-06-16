<script lang="ts">
  import { onMount } from "svelte";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import LeftSidebar from "../lib/components/LeftSidebar.svelte";
  import RightSidebar from "../lib/components/RightSidebar.svelte";
  import SettingsContent from "../lib/components/SettingsContent.svelte";
  import type { User } from "../lib/model/user";

  let isDarkMode = true;
  let user: User | null = null; // Dummy functions for sidebar compatibility
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  async function refreshData(page: number): Promise<void> {
    // Not needed for settings page
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  async function refreshFollowingData(page: number): Promise<void> {
    // Not needed for settings page
  }

  onMount(async () => {
    await checkAuth(navigate);

    const usr = await getAuthUser();
    if (usr) {
      user = usr;
    }

    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
      isDarkMode = savedTheme === "dark";
      document.documentElement.setAttribute("data-theme", savedTheme);
    } else {
      // Check user's system preference if no saved theme
      const prefersDark = window.matchMedia(
        "(prefers-color-scheme: dark)"
      ).matches;
      isDarkMode = prefersDark;
      document.documentElement.setAttribute(
        "data-theme",
        prefersDark ? "dark" : "light"
      );
    }
  });

  function toggleTheme() {
    isDarkMode = !isDarkMode;
    const newTheme = isDarkMode ? "dark" : "light";
    localStorage.setItem("theme", newTheme);
    document.documentElement.setAttribute("data-theme", newTheme);
  }
</script>

<div class="settings-container">
  <LeftSidebar
    {toggleTheme}
    {isDarkMode}
    {user}
    {refreshData}
    {refreshFollowingData}
    currentActiveTab="for-you"
  />
  {#if user}
    <SettingsContent {user} />
  {:else}
    <div class="loading">Loading...</div>
  {/if}
  <RightSidebar />
</div>

<style>
  .settings-container {
    display: flex;
    width: 100%;
    height: 100vh;
    background-color: var(--background);
    color: var(--text-primary);
  }

  .loading {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: var(--text-secondary);
    background-color: var(--card-background);
    border-left: 1px solid var(--border-color);
    border-right: 1px solid var(--border-color);
    max-width: 600px;
  }

  .loading::before {
    content: "";
    width: 32px;
    height: 32px;
    border: 3px solid var(--border-color);
    border-top: 3px solid var(--button-primary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
</style>
