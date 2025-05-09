<script lang="ts">
  import { onMount } from "svelte";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import type { User } from "../lib/model/user";
  import LeftSidebar from "../lib/components/LeftSidebar.svelte";
  import RightSidebar from "../lib/components/RightSidebar.svelte";
  import { fetchCommunityWithCategory, fetchJoinedWithCategory, fetchPendingWithCategory, type Community } from "../lib/model/community";
  import CommunityContent from "../lib/components/CommunityContent.svelte";

    let isDarkMode = true;
    let user: User | null = null;
    let communities: Community[] = [];
    let limit: 25 | 30 | 35 = 25;

    async function refreshData(page: number, search: string = "", currTab: "joined" | "pending" | "available" = "available",
        categoryId: string = "all"
    ): Promise<void> {
        // console.log("refreshing data");
        let data;
        if(user) {
            // console.log(currTab);
            if(currTab === "available") data = await fetchCommunityWithCategory(user.id, page, limit, search, categoryId);
            else if(currTab === "pending") data = await fetchPendingWithCategory(user.id, page, limit, search, categoryId);
            else data = await fetchJoinedWithCategory(user.id, page, limit, search, categoryId);

            if(data) communities = data;
            else communities = [];
        }
        // console.log(communities);
    }

    function changeLimit(count: 25 | 30 | 35) {
        limit = count;
    }

    function toggleTheme() {
      isDarkMode = !isDarkMode;
      const newTheme = isDarkMode ? "dark" : "light";
      localStorage.setItem("theme", newTheme);
      document.documentElement.setAttribute("data-theme", newTheme);
    }

    onMount(async() => {
        await checkAuth(navigate);

        const usr = await getAuthUser();
        if(usr) user = usr;

        const savedTheme = localStorage.getItem("theme");
        if (savedTheme) {
            isDarkMode = savedTheme === "dark";
            document.documentElement.setAttribute("data-theme", savedTheme);
        } else {
            const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
            isDarkMode = prefersDark;
            document.documentElement.setAttribute("data-theme", prefersDark ? "dark" : "light");
        }

        refreshData(1);
    });
</script>

<div class="home-container">
  <LeftSidebar {toggleTheme} {isDarkMode} {user} {refreshData} />
  <CommunityContent {communities} {user} {refreshData} {changeLimit} />
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