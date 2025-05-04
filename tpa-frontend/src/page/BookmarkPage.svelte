<script lang="ts">
  import { onMount } from "svelte";
  import { fetchBookmarkWithUser, type Thread } from "../lib/model/thread";
  import type { User, UserThread } from "../lib/model/user";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import LeftSidebar from "../lib/components/LeftSidebar.svelte";
  import RightSidebar from "../lib/components/RightSidebar.svelte";
  import BookmarkContent from "../lib/components/BookmarkContent.svelte";

    let isDarkMode = true;
    let user: User | null = null;
    let threads: (Thread & { userData: UserThread })[] = [];
    const currentActiveTab: "for-you" | "following" = "for-you";

    async function refreshData(page: number, search: string = ""): Promise<void> {
        let data;
        if(user) {
            console.log(search);
            data = await fetchBookmarkWithUser(page, user.id, search);

            if(data) threads = data;
            else threads = [];
        }
        // console.log(threads);
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

        await refreshData(1);
    });
</script>

<div class="home-container">
  <LeftSidebar {toggleTheme} {isDarkMode} {user} {refreshData} refreshFollowingData={refreshData} {currentActiveTab} />
  <BookmarkContent {threads} {user} {refreshData} />
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