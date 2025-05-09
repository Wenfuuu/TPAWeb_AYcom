<script lang="ts">
  import { onMount } from "svelte";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import type { User, UserThread } from "../lib/model/user";
  import LeftSidebar from "../lib/components/LeftSidebar.svelte";
  import RightSidebar from "../lib/components/RightSidebar.svelte";
  import ExploreContent from "../lib/components/ExploreContent.svelte";
  import { fetchLatestThread, type Thread } from "../lib/model/thread";
  import { fetchCommunity, type Community } from "../lib/model/community";

  let user: User | null = null;
  let isDarkMode = true;
  let latestThreads: (Thread & { userData: UserThread })[] = [];
  let communities: Community[] = [];

  async function refreshData(page: number,
    currTab: "top" | "latest" | "people" | "media" | "communities" = "latest",
    search: string = ""
  ): Promise<void> {
    let threadData;
    let communityData;
    if(user) {
      if(currTab === "latest") threadData = await fetchLatestThread(page, user.id, search);
      else if(currTab === "communities") communityData = await fetchCommunity();

      if(threadData) latestThreads = threadData;
      else latestThreads = [];

      if(communityData) communities = communityData;
      else communities = [];
    }
  }

  onMount(async() => {
    await checkAuth(navigate);

    const usr = await getAuthUser();
    if(usr) {
      user = usr;
    }

    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
      isDarkMode = savedTheme === "dark";
      document.documentElement.setAttribute("data-theme", savedTheme);
    } else {
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

<main class="home-container">
  <LeftSidebar {toggleTheme} {isDarkMode} {user} {refreshData} />
  <ExploreContent {user} {refreshData} {latestThreads}
  {communities} />
  <RightSidebar />
</main>

<style>
  .home-container {
    display: flex;
    width: 100%;
    height: 100vh;
    background-color: var(--background);
    color: var(--text-primary);
  }
</style>