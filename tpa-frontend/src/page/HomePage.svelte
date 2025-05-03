<script lang="ts">
  import { onMount } from "svelte";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import LeftSidebar from "../lib/components/LeftSidebar.svelte";
  import RightSidebar from "../lib/components/RightSidebar.svelte";
  import HomeContent from "../lib/components/HomeContent.svelte";
  import type { User, UserThread } from "../lib/model/user";
  import { fetchFollowingThread, fetchThreadWithUser, type Thread } from "../lib/model/thread";

  let isDarkMode = true;
  let user: User | null = null;
  let threads: (Thread & { userData: UserThread })[] = [];
  const currentActiveTab: "for-you" | "following" = "for-you";

  async function refreshData(page: number): Promise<void> {
    if(user) {
      const data = await fetchThreadWithUser(page, user.id);
      if(data) threads = data;
    }
    // console.log(threads);
  }

  async function refreshFollowingData(page: number): Promise<void> {
    if(user){
      const data = await fetchFollowingThread(page, user.id);
      if(data) threads = data;
    }
  }

  // Initialize theme from localStorage on component mount
  onMount(async() => {
    await checkAuth(navigate);

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

    await refreshData(1);
  });

  function toggleTheme() {
    isDarkMode = !isDarkMode;
    const newTheme = isDarkMode ? "dark" : "light";
    localStorage.setItem("theme", newTheme);
    document.documentElement.setAttribute("data-theme", newTheme);
  }
</script>

<div class="home-container">
  <LeftSidebar {toggleTheme} {isDarkMode} {user} {refreshData} {refreshFollowingData}
  {currentActiveTab} />
  <HomeContent {threads} {user} {refreshData} {refreshFollowingData} {currentActiveTab} />
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
</style>